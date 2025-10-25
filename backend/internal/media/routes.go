// backend/internal/media/routes.go
package media

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all media routes
func RegisterRoutes(r chi.Router, handler *Handler) {
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
