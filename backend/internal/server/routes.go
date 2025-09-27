package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/config"
	"github.com/iankencruz/threefive/internal/shared/middleware"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// setupRouter configures basic routes
func setupRouter(cfg *config.Config, db *pgxpool.Pool, queries *sqlc.Queries) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.CORS)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recovery)

	// Basic routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler(db))

	// ===========================
	// NOTE: API versioning
	// ===========================

	r.Route("/api/v1", func(r chi.Router) {
		// Register feature routes

		r.Get("/status", statusHandler)
		r.Get("/db-test", dbTestHandler(queries))
	})

	return r
}
