// internal/contacts/routes.go
package contacts

import (
	authmiddleware "github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/session"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all contact-related routes
func RegisterRoutes(r chi.Router, h *Handler, sessionManager *session.Manager) {
	authMiddleware := authmiddleware.NewMiddleware(sessionManager)

	// Public routes
	r.Post("/contact", h.CreateContact)

	// Admin routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		r.Get("/admin/contacts", h.ListContacts)
		r.Get("/admin/contacts/{id}", h.GetContact)
		r.Patch("/admin/contacts/{id}/status", h.UpdateContactStatus)
		r.Delete("/admin/contacts/{id}", h.DeleteContact)
	})
}
