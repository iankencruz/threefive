package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/config"
	"github.com/iankencruz/threefive/internal/media"
	"github.com/iankencruz/threefive/internal/shared/database"
	"github.com/iankencruz/threefive/internal/shared/session"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config         *config.Config
	server         *http.Server
	db             *pgxpool.Pool
	queries        *sqlc.Queries
	storage        storage.Storage
	sessionManager *session.Manager
	authHandler    *auth.Handler
	mediaHandler   *media.Handler
}

// New creates a new server instance with all dependencies initialized
func New(cfg *config.Config) (*Server, error) {
	// 1. Connect to database
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	queries := sqlc.New(db)

	// 2. Initialize storage based on config
	var storageInstance storage.Storage
	if cfg.Storage.Type == "s3" {
		storageInstance, err = storage.NewS3Storage(
			cfg.Storage.S3Bucket,
			cfg.Storage.S3Region,
			cfg.Storage.S3AccessKey,
			cfg.Storage.S3SecretKey,
			cfg.Storage.S3Endpoint,
			cfg.Storage.S3PublicURL,
			cfg.Storage.S3UseSSL,
		)
	} else {
		storageInstance, err = storage.NewLocalStorage(
			cfg.Storage.LocalBasePath,
			cfg.Storage.LocalBaseURL,
		)
	}
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// 3. Initialize session manager
	sessionConfig := session.DefaultConfig()
	sessionManager := session.NewManager(db, queries, sessionConfig)

	// 4. Initialize feature handlers (they create their own services)
	authHandler := auth.NewHandler(db, queries, sessionManager)
	mediaHandler := media.NewHandler(db, queries, storageInstance)
	// userHandler := user.NewHandler(db, queries)
	// projectHandler := project.NewHandler(db, queries)

	// Create server instance
	srv := &Server{
		config:  cfg,
		db:      db,
		queries: queries,
		storage: storageInstance,

		sessionManager: sessionManager,

		authHandler:  authHandler,
		mediaHandler: mediaHandler,
	}

	// 5. Setup router with all initialized components
	router := srv.setupRouter()

	// 6. Create HTTP server
	srv.server = &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("\n🚀 Server starting on http://%s\n", s.config.ServerAddress())
	fmt.Printf("📊 Environment: %s\n\n", s.config.Server.Env)

	if s.config.IsDevelopment() {
		fmt.Println("================================================")
		fmt.Println("🔥 Development mode")

		stats := s.db.Stat()
		fmt.Printf("📊 DB Pool: %d total, %d idle, %d in-use\n",
			stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
		fmt.Printf("💾 Storage Type: %s\n", s.storage.Type())
		fmt.Println("================================================")
	}

	// Start background cleanup routine
	ctx := context.Background()
	s.sessionManager.StartCleanupRoutine(ctx, 1*time.Hour)

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("🛑 Shutting down server...")

	if s.db != nil {
		s.db.Close()
		fmt.Println("✅ Database pool closed")
	}

	return s.server.Shutdown(ctx)
}

// Basic handlers...
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Welcome to the API", "version": "1.0.0"}`))
}

func healthHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status": "error", "database": "disconnected"}`))
			return
		}

		stats := db.Stat()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`{
            "status": "ok", 
            "service": "backend",
            "database": "connected",
            "pool": {
                "total": %d,
                "idle": %d,
                "in_use": %d
            }
        }`, stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
		w.Write([]byte(response))
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"api": "v1", "status": "running"}`))
}

func dbTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"database": "connection_test_ok", "sqlc": "ready"}`))
}
