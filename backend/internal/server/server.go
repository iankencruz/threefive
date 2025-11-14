// backend/internal/server/server.go
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/config"
	"github.com/iankencruz/threefive/internal/jobs"
	"github.com/iankencruz/threefive/internal/media"
	"github.com/iankencruz/threefive/internal/pages"
	"github.com/iankencruz/threefive/internal/projects"
	"github.com/iankencruz/threefive/internal/shared/database"
	"github.com/iankencruz/threefive/internal/shared/session"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	Config         *config.Config
	Server         *http.Server
	DB             *pgxpool.Pool
	Queries        *sqlc.Queries
	Storage        storage.Storage
	SessionManager *session.Manager
	AuthHandler    *auth.Handler
	MediaHandler   *media.Handler
	PageHandler    *pages.Handler
	ProjectHandler *projects.Handler
	CleanupWorker  *jobs.PageCleanupWorker // Add cleanup worker
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
	pageHandler := pages.NewHandler(db, queries, cfg)
	projectHandler := projects.NewHandler(db, queries)

	// 5. Initialize cleanup worker if enabled
	var cleanupWorker *jobs.PageCleanupWorker
	if cfg.AutoPurgeEnabled {
		cleanupWorker = jobs.NewPageCleanupWorker(queries, cfg.AutoPurgeRetentionDays)
	}

	// Create server instance
	srv := &Server{
		Config:         cfg,
		DB:             db,
		Queries:        queries,
		Storage:        storageInstance,
		SessionManager: sessionManager,
		CleanupWorker:  cleanupWorker,
		AuthHandler:    authHandler,
		MediaHandler:   mediaHandler,
		PageHandler:    pageHandler,
		ProjectHandler: projectHandler,
	}

	// 6. Create default admin user if it doesn't exist
	if err := srv.createDefaultAdminUser(context.Background()); err != nil {
		log.Printf("Warning: Failed to create default admin user: %v", err)
		// Don't fail server startup if this fails
	}

	// 7. Setup router with all initialized components
	router := srv.setupRouter()

	// 8. Create HTTP server
	srv.Server = &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv, nil
}

// Start starts the HTTP server and background workers
func (s *Server) Start() error {
	fmt.Printf("\n🚀 Server starting on http://%s\n", s.Config.ServerAddress())
	fmt.Printf("📊 Environment: %s\n\n", s.Config.Server.Env)
	fmt.Printf("📁 Storage: %s\n\n", s.Storage.Type())

	if s.Config.IsDevelopment() {
		fmt.Println("================================================")
		fmt.Println("🔥 Development mode")

		stats := s.DB.Stat()
		fmt.Printf("📊 DB Pool: %d total, %d idle, %d in-use\n",
			stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
		fmt.Printf("💾 Storage Type: %s\n", s.Storage.Type())
		fmt.Println("================================================")
	}

	// Start background cleanup routine for sessions
	ctx := context.Background()
	s.SessionManager.StartCleanupRoutine(ctx, 1*time.Hour)

	// Start page cleanup worker if enabled
	if s.CleanupWorker != nil {
		s.CleanupWorker.Start(ctx)
		log.Printf("✅ Page cleanup worker started (retention: %d days)", s.Config.AutoPurgeRetentionDays)
	} else {
		log.Println("ℹ️  Page auto-purge disabled")
	}

	return s.Server.ListenAndServe()
}

// Shutdown gracefully shuts down the server and background workers
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("🛑 Shutting down server...")

	// Stop cleanup worker first
	if s.CleanupWorker != nil {
		s.CleanupWorker.Stop()
		fmt.Println("✅ Page cleanup worker stopped")
	}

	// Close database connection
	if s.DB != nil {
		s.DB.Close()
		fmt.Println("✅ Database pool closed")
	}

	// Shutdown HTTP server
	return s.Server.Shutdown(ctx)
}

// createDefaultAdminUser creates a default admin user if one doesn't exist
func (s *Server) createDefaultAdminUser(ctx context.Context) error {
	// Check if any users exist
	// If this is your first time, there won't be a CountUsers method
	// You can implement it or just try to create the user and handle the error

	defaultEmail := "admin@example.com"
	defaultPassword := "Password!123"

	// Check if user exists
	_, err := s.Queries.GetUserByEmail(ctx, defaultEmail)
	if err == nil {
		// User already exists
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        defaultEmail,
		PasswordHash: string(hashedPassword),
		FirstName:    "Admin",
		LastName:     "User",
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Println("================================================")
	log.Println("🔐 Default admin user created!")
	log.Printf("   Email: %s", user.Email)
	log.Printf("   Password: %s", defaultPassword)
	log.Printf("   ⚠️  IMPORTANT: Change this password after first login!")
	log.Println("================================================")

	return nil
}

// Rest of your server methods (setupRouter, handlers, etc.)
// ... keep all your existing methods below this point
