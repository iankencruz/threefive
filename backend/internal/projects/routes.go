// backend/internal/projects/routes.go
package projects

import (
	"github.com/go-chi/chi/v5"
	authmiddleware "github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/session"
)

// RegisterRoutes creates and returns projects routes
func RegisterRoutes(r chi.Router, handler *Handler, sessionManager *session.Manager) {
	authMiddleware := authmiddleware.NewMiddleware(sessionManager)

	// ==========================================
	// 1. PUBLIC ROUTES
	// Base Path: /api/v1/projects
	// ==========================================
	r.Route("/projects", func(r chi.Router) {
		// Public routes - anyone can view projects
		r.Get("/", handler.ListPublishedProjects)
		r.Get("/{slug}", handler.GetProjectBySlug)
	})

	// ==========================================
	// 2. ADMIN / PROTECTED ROUTES
	// Base Path: /api/v1/admin/projects
	// ==========================================

	// Protected routes - require authentication
	r.Route("/admin/projects", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		r.Get("/", handler.ListProjects)
		r.Post("/", handler.CreateProject)
		r.Get("/{id}", handler.GetProjectByID)
		r.Put("/{id}", handler.UpdateProject)
		r.Patch("/{id}/status", handler.UpdateProjectStatus)
		r.Delete("/{id}", handler.DeleteProject)
	})
}
