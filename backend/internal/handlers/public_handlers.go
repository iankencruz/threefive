package handlers

import (
	"net/http"

	"github.com/iankencruz/threefive/internal/templates"
	"github.com/iankencruz/threefive/internal/templates/data"
	"github.com/iankencruz/threefive/ui/templates/layouts"
	"github.com/iankencruz/threefive/ui/templates/pages"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := data.HomeData{
		Headline: "Welcome to ThreeFive",
		Subtext:  "Your go-to platform for all things related to threefive.",
	}
	page := layouts.BaseLayout(data.Headline, pages.HomePage(data))
	templates.Render(w, r, page)
}
