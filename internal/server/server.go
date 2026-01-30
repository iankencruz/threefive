package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/iankencruz/threefive/database"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/internal/session"
	"github.com/iankencruz/threefive/pkg/logger"
	"github.com/labstack/echo/v5"
)

type Server struct {
	Echo              *echo.Echo
	DB                database.Service
	Queries           *generated.Queries
	AuthService       *services.AuthService
	MediaService      *services.MediaService
	SessionManager    *session.SessionManager
	SessionMiddleware *middleware.SessionMiddleware
	Log               *slog.Logger
}

func NewServer() *Server {
	var handler slog.Handler

	env := os.Getenv("ENV")

	if env == "production" {
		// JSON handler for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Pretty handler for development
		handler = logger.NewPrettyHandler(os.Stdout, logger.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		})
	}

	slogger := slog.New(handler)

	slogger.Info("Initializing server...")

	// Initialize database
	db := database.New(slogger)
	slogger.Info("database connection established")

	if err := db.RunMigrations("migrations"); err != nil {
		slogger.Error("failed to run migrations", "error", err)
		panic(err)
	}

	var storage services.StorageProvider
	storageProvider := os.Getenv("STORAGE_PROVIDER")

	if storageProvider == "s3" {
		// Initialize S3 storage (Vultr Object Storage or AWS S3)
		slogger.Info("Initializing S3 storage")
		s3Storage, err := services.NewS3Storage(context.Background(), services.S3Config{
			Bucket:          os.Getenv("S3_BUCKET"),
			Region:          os.Getenv("S3_REGION"),
			Endpoint:        os.Getenv("S3_ENDPOINT"), // For Vultr or other S3-compatible services
			AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
			BaseURL:         os.Getenv("S3_BASE_URL"), // Optional CDN URL
		})
		if err != nil {
			slogger.Error("failed to initialize S3 storage", "error", err)
			panic(err)
		}
		storage = s3Storage
		slogger.Info("S3 storage initialized", "bucket", os.Getenv("S3_BUCKET"), "region", os.Getenv("S3_REGION"))
	} else {
		// Initialize local storage (default for development)
		uploadDir := os.Getenv("LOCAL_UPLOAD_DIR")
		if uploadDir == "" {
			uploadDir = "./uploads"
		}
		baseURL := os.Getenv("LOCAL_BASE_URL")
		if baseURL == "" {
			baseURL = "/uploads"
		}
		storage = services.NewLocalStorage(uploadDir, baseURL)
		slogger.Info("Local storage initialized", "upload_dir", uploadDir, "base_url", baseURL)
	}

	// Initialize SQLC queries
	queries := generated.New(db.Pool())

	// Bootstrap database (ensure admin user exists)
	ctx := context.Background()
	if err := database.Bootstrap(ctx, queries, slogger); err != nil {
		slogger.Error("database bootstrap failed", "error", err)
		panic(err)
	}

	//
	// Initialize services
	authService := services.NewAuthService(db.Pool(), queries, slogger)
	slogger.Info("auth service initialized")

	// Initialize media service
	mediaService := services.NewMediaService(
		db.Pool(),
		queries,
		storage,
		services.MediaConfig{
			MaxFileSize: 250 * 1024 * 1024, // 250MB
			AllowedTypes: []string{
				"image/jpeg",
				"image/png",
				"image/gif",
				"image/webp",
				"video/mp4",
				"video/quicktime",
				"application/pdf",
			},
		},
	)
	sessionStore := session.NewPostgresStore(db.Pool(), queries, slogger)
	slogger.Info("session store initialized")

	// Initialize session manager (7 day lifetime)
	sessionManager := session.NewManager(sessionStore, 7*24*time.Hour, slogger)
	slogger.Info("session manager initialized")

	// Initialize middleware
	sessionMiddleware := middleware.NewSessionMiddleware(sessionManager, queries, slogger)
	slogger.Info("session middleware initialized")

	// Initialize Echo with the config that uses the silent logger
	e := echo.NewWithConfig(echo.Config{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	// Middlewares here
	e.Use(middleware.CustomRequestLogger())

	s := &Server{
		Echo:              e,
		DB:                db,
		Queries:           queries,
		AuthService:       authService,
		MediaService:      mediaService,
		SessionManager:    sessionManager,
		SessionMiddleware: sessionMiddleware,
		Log:               slogger,
	}

	s.RegisterRoutes()

	return s
}

// Start runs the Echo server on a specific port address
func (s *Server) Start(ctx context.Context, port string) error {
	if port == "" {
		port = ":8080" // fallback
	}

	// Start session cleanup
	go s.sessionCleanupWorker(ctx)

	srv := &http.Server{
		Addr:    port,
		Handler: s.Echo,
	}

	// Start server in background
	go func() {
		s.Log.Info("Starting server on %s", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log.Error("Server ListenAndServe failed")
		}
	}()

	// Wait for the signal context to be done
	<-ctx.Done()
	s.Log.Warn("signal received, beginning graceful shutdown")

	// Create a deadline for the shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		s.Log.Error("graceful shutdown failed", "error", err)
		return err
	}

	// Only close database pool after server has shutdown
	s.DB.Close()

	s.Log.Info("Server exited cleanly")
	return nil
}

// sessionCleanupWorker periodically removes expired sessions
func (s *Server) sessionCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	s.Log.Info("session cleanup worker started")

	for {
		select {
		case <-ctx.Done():
			s.Log.Info("session cleanup worker stopped")
			return
		case <-ticker.C:
			if err := s.SessionManager.Cleanup(ctx); err != nil {
				s.Log.Error("failed to cleanup expired sessions", "error", err)
			} else {
				s.Log.Debug("cleaned up expired sessions")
			}
		}
	}
}
