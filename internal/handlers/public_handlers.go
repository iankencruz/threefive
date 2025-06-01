package handlers

import (
	"net/http"

	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/templates/data"
	"github.com/iankencruz/threefive/ui/templates/pages"
)

func HomeHandler(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageData := data.HomeData{
			Headline: "Welcome to ThreeFive",
			Subtext:  "Curated cinematic storytelling.",
		}

		app.Render(w, r, "Home", pages.HomePage(pageData))
	}
}
