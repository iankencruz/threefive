package public

import (
	"net/http"

	"github.com/iankencruz/threefive/backend/internal/core/templates"
	"github.com/iankencruz/threefive/backend/internal/core/viewdata"
	"github.com/iankencruz/threefive/backend/ui/pages"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	data := viewdata.HomePageData{
		MetaData: viewdata.NewMeta(r, "Home", "Welcome to the ThreeFive Project"),
		Feature:  "Cinematic Storytelling",
	}

	templates.Render(w, r, pages.HomePage(data))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := viewdata.AboutPageData{
		MetaData: viewdata.NewMeta(r, "About Page", "Learn more about the ThreeFive Project"),
	}
	templates.Render(w, r, pages.AboutPage(data))
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	data := viewdata.ContactPageData{
		MetaData:     viewdata.NewMeta(r, "Contact", "Get in touch with the ThreeFive Project team"),
		ContactEmail: "contactemail@contact.com",
	}
	templates.Render(w, r, pages.ContactPage(data))
}
