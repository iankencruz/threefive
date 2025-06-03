package application

import (
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
		r.Use(middleware.RedirectIfAuthenticated(app.SessionManager))
		r.Get("/login", app.AuthHandler.LoginHandler)
		r.Post("/login", app.AuthHandler.LoginHandler)
		r.Get("/register", app.AuthHandler.RegisterHandler)
		r.Post("/register", app.AuthHandler.RegisterHandler)
	})

	// ✅ Serve SvelteKit admin SPA
	adminStatic := http.Dir(filepath.Join("static", "admin"))
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(app.AuthHandler))
		r.Handle("/admin/*", http.StripPrefix("/admin", http.FileServer(adminStatic)))
	})

	// ✅ Example API group for admin SPA
	r.Route("/api/admin", func(r chi.Router) {
		r.Use(middleware.RequireAuth(app.AuthHandler))
		r.Get("/me", app.AuthHandler.GetAuthenticatedUser)
		// Add other secured endpoints here
	})

	return r
}
