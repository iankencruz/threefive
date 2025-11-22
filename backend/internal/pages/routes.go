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

	// ==========================================
	// 1. PUBLIC ROUTES
	// Base Path: /api/v1/pages
	// ==========================================
	r.Route("/pages", func(r chi.Router) {
		// Public routes - anyone can view pages
		r.Get("/", handler.ListPublishedPages)
		r.Get("/{slug}", handler.GetPageBySlug)
	})

	// ==========================================
	// 2. ADMIN / PROTECTED ROUTES
	// Base Path: /api/v1/admin/pages
	// ==========================================
	r.Route("/admin/pages", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		r.Get("/", handler.ListPages)
		r.Post("/", handler.CreatePage)
		r.Get("/{id}", handler.GetPageByID)
		r.Put("/{id}", handler.UpdatePage)
		r.Patch("/{id}/status", handler.UpdatePageStatus)
		r.Delete("/{id}", handler.DeletePage)
		r.Post("/purge", handler.PurgeDeletedPages)
	})
}
