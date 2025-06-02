package routes

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/handlers"
	"github.com/iankencruz/threefive/internal/middleware"
)

func Routes(app *application.Application) http.Handler {
	r := chi.NewRouter()

	guestOnly := middleware.RedirectIfAuthenticated(app)
	adminOnly := middleware.RequireAuth(app)

	// Log requests/
	r.Use(middleware.RequestLogger(app.Logger))
	// r.Use(mw.Logger)

	// ✅ Serve static files (e.g., /static/css/output.css)
	staticDir := http.Dir(filepath.Join("ui", "static"))
	fileServer := http.FileServer(staticDir)
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Get("/", handlers.HomeHandler(app))

	// Authentication Routes
	r.Group(func(r chi.Router) {
		r.Use(guestOnly)
		r.Get("/login", handlers.LoginHandler(app))
		r.Post("/login", handlers.LoginHandler(app))
		r.Get("/register", handlers.RegisterUserHandler(app))
		r.Post("/register", handlers.RegisterUserHandler(app))
	})

	// Admin Routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(adminOnly)
		r.Post("/logout", handlers.LogoutHandler(app))
		r.Get("/dashboard", handlers.DashboardPage(app))
		r.Get("/users", handlers.UsersPage(app))
		// r.Get("/dashboard", handlers.AdminDashboardHandler(app))
		// other admin routes...
	})

	return r
}
