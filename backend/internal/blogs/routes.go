// backend/internal/blogs/routes.go
package blogs

import (
	"github.com/go-chi/chi/v5"
	authmiddleware "github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/session"
)

// RegisterRoutes registers all blog routes
func RegisterRoutes(r chi.Router, handler *Handler, sessionManager *session.Manager) {
	authMiddleware := authmiddleware.NewMiddleware(sessionManager)

	r.Route("/blogs", func(r chi.Router) {
		// Public routes (no authentication required)
		r.Group(func(r chi.Router) {
			r.Get("/", handler.ListBlogs)         // GET /api/v1/blogs
			r.Get("/{idOrSlug}", handler.GetBlog) // GET /api/v1/blogs/{idOrSlug}
		})

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			r.Post("/", handler.CreateBlog)                   // POST /api/v1/blogs
			r.Put("/{id}", handler.UpdateBlog)                // PUT /api/v1/blogs/{id}
			r.Patch("/{id}/status", handler.UpdateBlogStatus) // PATCH /api/v1/blogs/{id}/status
			r.Delete("/{id}", handler.DeleteBlog)             // DELETE /api/v1/blogs/{id}
			r.Post("/purge", handler.PurgeOldDeletedBlogs)    // POST /api/v1/blogs/purge
		})
	})
}
