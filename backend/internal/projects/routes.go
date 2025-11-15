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

	r.Route("/projects", func(r chi.Router) {
		// Public routes - anyone can view projects
		r.Group(func(r chi.Router) {
			r.Get("/", handler.ListProjects)           // GET /api/v1/projects
			r.Get("/{slug}", handler.GetProjectBySlug) // GET /api/v1/projects/{slug}
		})

		// Protected routes - require authentication
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			r.Post("/", handler.CreateProject)                   // POST /api/v1/projects
			r.Put("/{id}", handler.UpdateProject)                // PUT /api/v1/projects/{id}
			r.Patch("/{id}/status", handler.UpdateProjectStatus) // PATCH /api/v1/projects/{id}/status
			r.Delete("/{id}", handler.DeleteProject)             // DELETE /api/v1/projects/{id}
		})
	})
}
