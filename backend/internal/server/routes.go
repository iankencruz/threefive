package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/config"
	"github.com/iankencruz/threefive/internal/shared/middleware"
	"github.com/iankencruz/threefive/internal/shared/session"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// setupRouter configures basic routes
func setupRouter(cfg *config.Config, db *pgxpool.Pool, queries *sqlc.Queries) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recovery)

	// Application-level middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.Frontend.URL}, // From config
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		MaxAge:           300,
		AllowCredentials: true,
	}))

	// Create session manager
	sessionConfig := session.DefaultConfig()
	sessionManager := session.NewManager(db, queries, sessionConfig)

	// Auth middleware (auth package)
	authMW := auth.NewMiddleware(sessionManager)

	// Create Auth Service
	authService := auth.NewService(db, queries, sessionManager)
	authHandler := auth.NewHandler(authService, sessionManager)
	authMiddleware := auth.NewMiddleware(sessionManager)

	// Mount auth routes
	authRoutes := auth.Routes(authHandler, authMiddleware)
	r.Mount("/auth", authRoutes)

	// Basic routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler(db))

	// ===========================
	// NOTE: API versioning
	// ===========================

	r.Route("/api/v1", func(r chi.Router) {
		// Auth Middleware
		r.Use(authMW.RequireAuth)

		// Register feature routes
		r.Get("/status", statusHandler)
		r.Get("/db-test", dbTestHandler(queries))
	})

	return r
}
