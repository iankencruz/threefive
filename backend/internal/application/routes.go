package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/core/middleware"
)

func Routes(app *Application) http.Handler {
	r := chi.NewRouter()

	// ✅ Global middlewares
	r.Use(middleware.RequestLogger)

	// ✅ Versioned API
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

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
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth(app.AuthHandler))

				// === /media ===
				r.Route("/media", func(r chi.Router) {
					r.Get("/", app.MediaHandler.ListMediaHandler)
					r.Post("/", app.MediaHandler.UploadMediaHandler)
					r.Route("/{id}", func(r chi.Router) {
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
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", app.ProjectHandler.Get)
						r.Put("/", app.ProjectHandler.Update)
						r.Delete("/", app.ProjectHandler.Delete)

						r.Post("/media", app.ProjectHandler.AddMedia)
						r.Delete("/media", app.ProjectHandler.RemoveMedia)
						r.Post("/media/sort", app.ProjectHandler.UpdateSortOrder)
					})
				})

				// === TODO: Add galleries, blog, settings etc. here ===
			})
		})
	})

	return r
}
