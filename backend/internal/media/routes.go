// backend/internal/media/routes.go
package media

import (
	"github.com/a-h/templ/cfg"
	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes registers all media routes
func RegisterRoutes(r chi.Router, db *pgxpool.Pool, queries *sqlc.Queries) {
	// Initialize storage (using local for now)
	// storage := NewLocalStorage("./uploads")
	storage := storage.NewS3Storage()

	// Initialize service
	service := NewService(db, queries)

	// Initialize handler
	handler := NewHandler(service)

	// Media routes
	r.Route("/media", func(r chi.Router) {
		// Upload endpoint - requires authentication
		r.Post("/upload", handler.UploadHandler)

		// Get single media
		r.Get("/{id}", handler.GetMediaHandler)

		// List all media (with pagination)
		r.Get("/", handler.ListMediaHandler)

		// Delete media (soft delete by default, ?hard=true for permanent)
		r.Delete("/{id}", handler.DeleteMediaHandler)

		// Link/unlink media to entities
		r.Post("/{id}/link", handler.LinkMediaHandler)
		r.Delete("/{id}/link", handler.UnlinkMediaHandler)

		// Get media for an entity
		r.Get("/entity/{type}/{id}", handler.GetEntityMediaHandler)

		// Statistics
		r.Get("/stats", handler.GetStatsHandler)
	})
}
