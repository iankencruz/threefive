package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/handlers"
)

func Routes(app *application.Application) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handlers.HomeHandler(app))

	return r
}
