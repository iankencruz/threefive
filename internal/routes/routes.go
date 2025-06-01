package routes

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	// mw "github.com/go-chi/chi/v5/middleware"
	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/handlers"
	"github.com/iankencruz/threefive/internal/middleware"
)

func Routes(app *application.Application) http.Handler {
	r := chi.NewRouter()

	// Log requests/
	r.Use(middleware.RequestLogger(app.Logger))
	// r.Use(mw.Logger)

	// ✅ Serve static files (e.g., /static/css/output.css)
	staticDir := http.Dir(filepath.Join("ui", "static"))
	fileServer := http.FileServer(staticDir)
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Get("/", handlers.HomeHandler(app))

	// Authentication Routes
	r.Get("/login", handlers.LoginUserHandler(app))
	r.Get("/register", handlers.RegisterUserHandler(app))

	r.Post("/login", handlers.LoginUserHandler(app))
	r.Post("/register", handlers.RegisterUserHandler(app))

	return r
}
