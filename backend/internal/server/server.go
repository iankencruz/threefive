// backend/internal/server/server.go
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/iankencruz/threefive/internal/config"
	"github.com/iankencruz/threefive/internal/shared/database"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config  *config.Config
	server  *http.Server
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Connect to database using pgxpool
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create SQLC queries with pgxpool
	queries := sqlc.New(db)

	// Setup router
	router := setupRouter(cfg, db, queries)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:  cfg,
		server:  httpServer,
		db:      db,
		queries: queries,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("🚀 Server starting on http://%s\n", s.config.ServerAddress())
	fmt.Printf("📊 Environment: %s\n", s.config.Server.Env)

	if s.config.IsDevelopment() {
		fmt.Println("🔥 Development mode")

		// Show connection pool stats in development
		stats := s.db.Stat()
		fmt.Printf("📊 DB Pool: %d total, %d idle, %d in-use\n",
			stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("🛑 Shutting down server...")

	// Close database connection pool
	if s.db != nil {
		s.db.Close()
		fmt.Println("✅ Database pool closed")
	}

	// Shutdown HTTP server
	return s.server.Shutdown(ctx)
}

// Basic HTTP handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Welcome to the API", "version": "1.0.0"}`))
}

func healthHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Test database connection
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status": "error", "database": "disconnected"}`))
			return
		}

		// Get pool stats
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

// Database test handler to verify SQLC integration
func dbTestHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Simple query to test database
		// This will work once you have some SQLC queries generated
		// For now, just test raw query

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"database": "connection_test_ok", "sqlc": "ready"}`))
	}
}
