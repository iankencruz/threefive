package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/core/middleware"
)

func Routes(app *Application) http.Handler {
	r := chi.NewRouter()

	// CORS Middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "http://localhost:5173" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	})

	// ✅ Global middlewares
	r.Use(middleware.RequestLogger)

	// ✅ Versioned API
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// === / public routes ===
			r.Group(func(r chi.Router) {
				r.Get("/home", app.PageHandler.HomePage)

				r.Get("/about", app.PageHandler.AboutPage)

				r.Get("/projects", app.MediaHandler.ListMediaHandler) // replace with PageHandler for projects

				r.Get("/contact", app.PageHandler.ContactPage)
			})

			// === /auth routes (split inside)
			r.Route("/auth", func(r chi.Router) {
				// Public endpoints
				r.Post("/login", app.AuthHandler.LoginHandler)
				r.Post("/register", app.AuthHandler.RegisterHandler)

				// Authenticated endpoints
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireAuth(app.AuthHandler))

					r.Get("/me", app.AuthHandler.GetAuthenticatedUser)
					r.Post("/logout", app.AuthHandler.LogoutHandler)
				})
			})

			// === Authenticated API routes
			r.Route("/admin", func(r chi.Router) {
				r.Use(middleware.RequireAuth(app.AuthHandler))

				// === /media ===
				r.Route("/media", func(r chi.Router) {
					r.Get("/", app.MediaHandler.ListMediaHandler)
					r.Post("/", app.MediaHandler.UploadMediaHandler)
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", app.MediaHandler.Get)
						r.Put("/", app.MediaHandler.UpdateMediaHandler)
						r.Delete("/", app.MediaHandler.DeleteMediaHandler)
					})
					r.Post("/sort", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, "Media sort endpoint (not yet implemented)")
					})
				})

				// === /projects ===
				r.Route("/projects", func(r chi.Router) {
					r.Get("/", app.ProjectHandler.List)
					r.Post("/", app.ProjectHandler.Create)
					r.Route("/{slug}", func(r chi.Router) {
						r.Get("/", app.ProjectHandler.GetAdminProjects)
						r.Put("/", app.ProjectHandler.Update)
						r.Delete("/", app.ProjectHandler.Delete)

						r.Post("/media", app.ProjectHandler.AddMedia)
						r.Delete("/media", app.ProjectHandler.RemoveMedia)
						r.Put("/media/sort", app.ProjectHandler.UpdateSortOrder)
					})
				})

				// === /pages ===
				r.Route("/pages", func(r chi.Router) {
					r.Get("/", app.PageHandler.List)
					r.Post("/", app.PageHandler.Create)
					r.Route("/{slug}", func(r chi.Router) {
						r.Get("/", app.PageHandler.GetAdminPages)
						r.Put("/", app.PageHandler.Update)
						r.Delete("/", app.PageHandler.Delete)
					})
				})

				// === /blogs ===
				r.Route("/blogs", func(r chi.Router) {
					r.Get("/", app.PageHandler.List)
					r.Post("/", app.PageHandler.Create)
					r.Route("/{slug}", func(r chi.Router) {
						r.Get("/", app.PageHandler.GetAdminPages)
						r.Put("/", app.PageHandler.Update)
						// r.Delete("/", app.PageHandler.Delete)
						// r.Put("/blocks/sort", app.PageHandler.SortBlocks)
					})
				})

				// === /blocks ===
				r.Route("/blocks", func(r chi.Router) {
					r.Post("/", app.BlockHandler.Create)
					r.Put("/{id}", app.BlockHandler.Update)
					r.Delete("/{id}", app.BlockHandler.Delete)
				})

				// === TODO: Add galleries, blog, settings etc. here ===
			})
		})
	})

	return r
}
