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

	// Define custom time format
	// logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	// 	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
	// 		if a.Key == slog.TimeKey {
	// 			t := a.Value.Time()
	// 			// Format the time as "YYYY-MM-DD HH:MM:SS"
	// 			a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
	// 		}
	// 		return a
	// 	},
	// }),
	// )

	slogger.Info("Initializing server...")

	// Initialize database
	db := database.New(slogger)
	slogger.Info("database connection established")

	// Initialize SQLC queries
	queries := generated.New(db.Pool())

	// Bootstrap database (ensure admin user exists)
	ctx := context.Background()
	if err := database.Bootstrap(ctx, db.Pool(), queries, slogger); err != nil {
		slogger.Error("database bootstrap failed", "error", err)
		panic(err)
	}

	// Initialize services
	authService := services.NewAuthService(db.Pool(), queries, slogger)
	slogger.Info("auth service initialized")

	// Initialize session store
	sessionStore := session.NewPostgresStore(db.Pool(), queries, slogger)
	slogger.Info("session store initialized")

	// Initialize session manager (7 day lifetime)
	sessionManager := session.NewManager(sessionStore, 7*24*time.Hour, slogger)
	slogger.Info("session manager initialized")

	// Initialize middleware
	sessionMiddleware := middleware.NewSessionMiddleware(sessionManager, slogger)
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
