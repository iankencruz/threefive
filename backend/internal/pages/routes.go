// backend/internal/pages/routes.go
package pages

import (
	"github.com/go-chi/chi/v5"
	authmiddleware "github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/session"
)

// RegisterRoutes creates and returns pages routes
func RegisterRoutes(r chi.Router, handler *Handler, sessionManager *session.Manager) {
	authMiddleware := authmiddleware.NewMiddleware(sessionManager)

	r.Route("/pages", func(r chi.Router) {
		// Public routes - anyone can view pages
		r.Group(func(r chi.Router) {
			r.Get("/", handler.ListPages)           // GET /api/v1/pages
			r.Get("/{slug}", handler.GetPageBySlug) // GET /api/v1/pages/{slug}
		})

		// Protected routes - require authentication
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			r.Post("/", handler.CreatePage)                   // POST /api/v1/pages
			r.Put("/{id}", handler.UpdatePage)                // PUT /api/v1/pages/{id}
			r.Patch("/{id}/status", handler.UpdatePageStatus) // PATCH /api/v1/pages/{id}/status
			r.Delete("/{id}", handler.DeletePage)             // DELETE /api/v1/pages/{id}

			r.Post("/purge", handler.PurgeDeletedPages)
		})
	})
}
