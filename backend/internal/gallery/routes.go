package gallery

import (
	"github.com/go-chi/chi/v5"
	authmiddleware "github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/session"
)

func RegisterRoutes(r chi.Router, handler *Handler, sessionManager *session.Manager) {
	authMiddleware := authmiddleware.NewMiddleware(sessionManager)

	r.Route("/galleries", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/", handler.ListGalleries)
			r.Get("/{id}", handler.GetGallery)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			r.Post("/", handler.CreateGallery)
			r.Put("/{id}", handler.UpdateGallery)
			r.Delete("/{id}", handler.DeleteGallery)

			// Media linking/unlinking
			r.Post("/{id}/media", handler.LinkMedia)
			r.Delete("/{id}/media/{mediaId}", handler.UnlinkMedia)
		})
	})
}
