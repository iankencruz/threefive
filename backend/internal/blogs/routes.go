// backend/internal/blogs/routes.go
package blogs

import (
	"log"

	"github.com/go-chi/chi/v5"
	authmiddleware "github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/session"
)

func RegisterRoutes(r chi.Router, handler *Handler, sessionManager *session.Manager) {
	authMiddleware := authmiddleware.NewMiddleware(sessionManager)

	// ==========================================
	// 1. PUBLIC ROUTES
	// Base Path: /api/v1/blogs
	// ==========================================
	r.Route("/blogs", func(r chi.Router) {
		// Anyone can list or read by slug
		r.Get("/", handler.ListBlogs)           // GET /api/v1/blogs
		r.Get("/{slug}", handler.GetBlogBySlug) // GET /api/v1/blogs/my-first-post
	})

	// ==========================================
	// 2. ADMIN / PROTECTED ROUTES
	// Base Path: /api/v1/admin/blogs
	// ==========================================
	r.Route("/admin/blogs", func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		log.Println("📍 Registering GET /admin/blogs/{id}")
		r.Get("/{id}", handler.GetBlogByID)

		r.Post("/", handler.CreateBlog)
		r.Put("/{id}", handler.UpdateBlog)
		r.Patch("/{id}/status", handler.UpdateBlogStatus)
		r.Delete("/{id}", handler.DeleteBlog)
		r.Post("/purge", handler.PurgeOldDeletedBlogs)
	})
}
