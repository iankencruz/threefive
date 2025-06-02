package handlers

import (
	"net/http"

	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/templates"
	"github.com/iankencruz/threefive/ui/templates/admin"
	"github.com/iankencruz/threefive/ui/templates/layouts"
)

func DashboardPage(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := layouts.AdminLayout(r, "admin", admin.AdminDashboardPage())
		templates.Render(w, r, page)
	}
}

func UsersPage(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := layouts.AdminLayout(r, "admin", admin.AdminDashboardPage())
		templates.Render(w, r, page)
	}
}
