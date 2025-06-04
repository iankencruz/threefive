package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/iankencruz/threefive/backend/internal/core/sessions"
	"github.com/iankencruz/threefive/backend/internal/core/templates"
	"github.com/iankencruz/threefive/backend/internal/core/validators"
	"github.com/iankencruz/threefive/backend/internal/core/viewdata"
	"github.com/iankencruz/threefive/backend/ui/pages"
)

type Handler struct {
	Service        Service
	SessionManager *sessions.Manager
	Logger         *slog.Logger
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := viewdata.LoginPageData{
			MetaData: viewdata.NewMeta(r, "Login", "Login to your account"),
		}
		templates.Render(w, r, pages.LoginPage(data))
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	v := validators.New()
	v.Require("email", email)
	v.Require("password", password)

	if !v.Valid() {
		data := viewdata.LoginPageData{
			MetaData: viewdata.NewMeta(r, "Login", "Login to your account"),
			Email:    email,
			Errors:   v.Errors,
		}
		templates.Render(w, r, pages.LoginPage(data))
		return
	}

	user, err := h.Service.Login(r.Context(), email, password)
	if err != nil {
		v.Errors["email"] = "Invalid email or password"
		data := viewdata.LoginPageData{
			MetaData: viewdata.NewMeta(r, "Login", "Login to your account"),
			Email:    email,
			Errors:   v.Errors,
		}
		templates.Render(w, r, pages.LoginPage(data))
		return
	}

	// ✅ Create session and set user ID cookie
	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := viewdata.RegisterPageData{
			MetaData: viewdata.NewMeta(r, "Register", "Create your account"),
		}
		templates.Render(w, r, pages.RegisterPage(data))
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")

	v := validators.New()
	v.Require("first_name", firstName)
	v.Require("last_name", lastName)
	v.Require("email", email)
	v.MatchPattern("email", email, validators.EmailRX, "Invalid email address")
	v.Require("password", password)
	v.MatchPattern("password", password, validators.UppercaseRX, "Must include at least one uppercase letter")
	v.MatchPattern("password", password, validators.NumberRX, "Must include at least one number")

	if !v.Valid() {
		data := viewdata.RegisterPageData{
			MetaData:  viewdata.NewMeta(r, "Register", "Create your account"),
			Errors:    v.Errors,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		}
		templates.Render(w, r, pages.RegisterPage(data))
		return
	}

	user, err := h.Service.Register(r.Context(), firstName, lastName, email, password)
	if err != nil {
		v.Errors["email"] = "Registration failed"
		data := viewdata.RegisterPageData{
			MetaData:  viewdata.NewMeta(r, "Register", "Create your account"),
			Errors:    v.Errors,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		}
		templates.Render(w, r, pages.RegisterPage(data))
		return
	}

	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement handler logic
}

func (h *Handler) GetUserID(r *http.Request) (int32, error) {
	return h.SessionManager.GetUserID(r)
}

func (h *Handler) LoadUser(ctx context.Context, userID int32) (any, error) {
	return h.Service.GetUserByID(ctx, userID)
}

func (h *Handler) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.SessionManager.GetUserID(r)
	w.Header().Set("Content-Type", "application/json")

	if err != nil || userID == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"user": nil})
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"user": nil})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
}
