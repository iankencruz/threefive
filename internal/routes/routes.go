package routes

import (
	"net/http"

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

	r.Get("/", handlers.HomeHandler(app))

	return r
}
