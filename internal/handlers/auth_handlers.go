package handlers

import (
	"net/http"
	"strings"

	"github.com/iankencruz/threefive/internal/application"
	models "github.com/iankencruz/threefive/internal/models/users"
	"github.com/iankencruz/threefive/internal/templates"
	"github.com/iankencruz/threefive/internal/templates/data"
	"github.com/iankencruz/threefive/internal/validators"
	"github.com/iankencruz/threefive/ui/templates/layouts"
	"github.com/iankencruz/threefive/ui/templates/pages"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		form := data.LoginForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		data := data.LoginPageData{
			Title:  "Login",
			Form:   struct{ Email string }{},
			Errors: map[string]string{},
		}

		if r.Method == http.MethodGet {
			page := layouts.BlankLayout(data.Title, pages.LoginPage(data))
			templates.Render(w, r, page)
			return
		}

		// POST logic
		if err := r.ParseForm(); err != nil {
			app.Logger.Error("Failed to parse login form", "error", err)
			return
		}

		v := validators.New()
		v.Require("email", form.Email)
		v.Require("password", form.Password)

		data.Form.Email = form.Email
		data.Errors = v.Errors

		if !v.Valid() {
			page := layouts.BlankLayout(data.Title, pages.LoginPage(data))
			templates.Render(w, r, page)
			return
		}

		// 🔐 Lookup user
		user, err := app.UserRepo.GetUserByEmail(r.Context(), form.Email)
		if err != nil {
			// Avoid leaking which field failed
			app.Logger.Warn("Login failed: user not found", "email", form.Email)
			data.Errors["email"] = "Invalid email or password"
			page := layouts.BlankLayout(data.Title, pages.LoginPage(data))
			templates.Render(w, r, page)
			return
		}

		// 🔐 Compare password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(form.Password))
		if err != nil {
			app.Logger.Warn("Login failed: invalid password", "email", form.Email)
			data.Errors["email"] = "Invalid email or password"
			page := layouts.BlankLayout(data.Title, pages.LoginPage(data))
			templates.Render(w, r, page)
			return
		}

		// ✅ Success – log user in
		err = app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.Logger.Error("Failed to renew session token", "error", err)
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		app.SessionManager.Put(r.Context(), "userID", int(user.ID))

		app.Logger.Info("User logged in", "email", user.Email, "id", user.ID)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}
}

func RegisterUserHandler(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		form := data.RegisterForm{
			FirstName: r.FormValue("first_name"),
			LastName:  r.FormValue("last_name"),
			Email:     strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
			Password:  r.FormValue("password"),
		}

		// ✅ Validate using your Validator
		v := validators.New()
		v.Require("first_name", form.FirstName)
		v.Require("last_name", form.LastName)
		v.Require("email", form.Email)
		v.MatchPattern("email", form.Email, validators.EmailRX, "Invalid email address")
		v.Require("password", form.Password)
		v.MatchPattern("password", form.Password, validators.UppercaseRX, "Must include at least one uppercase letter")
		v.MatchPattern("password", form.Password, validators.NumberRX, "Must include at least one number")

		if !v.Valid() {
			// You can re-render the form with error messages here
			page := layouts.BlankLayout("Register", pages.RegisterPage(data.RegisterPageData{
				Form:   form,
				Errors: v.Errors,
			}))
			templates.Render(w, r, page)
			return
		}

		// ✅ Proceed to hash password and create user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		params := models.CreateUserParams{
			FirstName:    form.FirstName,
			LastName:     form.LastName,
			Email:        form.Email,
			PasswordHash: string(hashedPassword),
		}

		user, err := app.UserRepo.CreateUser(r.Context(), params)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// ✅ Log user in via session
		_ = app.SessionManager.RenewToken(r.Context())
		app.SessionManager.Put(r.Context(), "userID", user.ID)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogoutHandler(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Remove session + regenerate token
		err := app.SessionManager.Destroy(r.Context())
		if err != nil {
			app.Logger.Error("Logout failed", "error", err)
			http.Error(w, "Logout failed", http.StatusInternalServerError)
			return
		}

		app.Logger.Info("User logged out")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
