package application

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/backend/internal/core/middleware"
	"github.com/iankencruz/threefive/backend/internal/public"
)

func Routes(app *Application) http.Handler {
	r := chi.NewRouter()

	// ✅ Request logging middleware
	r.Use(middleware.RequestLogger)

	// ✅ Serve public static assets
	publicStatic := http.Dir(filepath.Join("ui", "static"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(publicStatic)))

	// ✅ Public site pages
	r.Get("/", public.HomeHandler)
	r.Get("/about", public.AboutHandler)
	r.Get("/contact", public.ContactHandler)

	// ✅ Auth endpoints
	r.Group(func(r chi.Router) {
		r.Get("/login", app.AuthHandler.LoginHandler)
		r.Post("/login", app.AuthHandler.LoginHandler)
		r.Get("/register", app.AuthHandler.RegisterHandler)
		r.Post("/register", app.AuthHandler.RegisterHandler)
	})

	// ✅ Serve SvelteKit Admin SPA — protected by middleware
	adminStatic := http.Dir(filepath.Join("static", "admin"))
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(app.AuthHandler))
		r.Handle("/admin/*", http.StripPrefix("/admin", http.FileServer(adminStatic)))
	})

	// ✅ Versioned API group for /admin — scoped under /api/v1/admin
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/admin", func(r chi.Router) {
				// Public admin auth endpoints
				r.Post("/login", app.AuthHandler.LoginHandler)
				r.Post("/register", app.AuthHandler.RegisterHandler)

				// Protected admin endpoints
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireAuth(app.AuthHandler))
					r.Post("/logout", app.AuthHandler.LogoutHandler)
					r.Get("/me", app.AuthHandler.GetAuthenticatedUser)

					// Media management
					r.Get("/media", app.MediaHandler.ListMediaHandler)
					r.Post("/media", app.MediaHandler.UploadMediaHandler)
					r.Put("/media/{id}", func(w http.ResponseWriter, r *http.Request) {
						fmt.Print("Update media metadata Endpoint:")
					})
					r.Delete("/media/{id}", func(w http.ResponseWriter, r *http.Request) {

						fmt.Print("Delete media instance Endpoint:")
					})
					r.Post("/sort", func(w http.ResponseWriter, r *http.Request) {
						fmt.Print("Sort media Endpoint:")
					})

					// Future protected admin APIs (e.g. projects, galleries, etc.)
				})
			})
		})
	})

	return r
}
