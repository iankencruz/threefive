// backend/internal/media/routes.go
package media

import (
	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/session"
)

// RegisterRoutes registers all media routes
func RegisterRoutes(r chi.Router, handler *Handler, sessionManager *session.Manager) {
	authMiddleware := auth.NewMiddleware(sessionManager)

	// Media routes
	r.Route("/media", func(r chi.Router) {
		// Public routes - anyone can view media
		r.Group(func(r chi.Router) {
			r.Get("/{id}", handler.GetMediaHandler)
			r.Get("/", handler.ListMediaHandler)
			r.Get("/entity/{type}/{id}", handler.GetEntityMediaHandler)
		})

		// Protected routes - require authentication for all admin operations
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			r.Post("/upload", handler.UploadHandler)
			r.Delete("/{id}", handler.DeleteMediaHandler)
			r.Post("/{id}/link", handler.LinkMediaHandler)
			r.Delete("/{id}/link", handler.UnlinkMediaHandler)
			r.Get("/stats", handler.GetStatsHandler)
		})
	})
}
