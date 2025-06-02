package middleware

import (
	"context"
	"net/http"

	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/contextkeys"
)

// User type should match your actual user model
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	// Add other fields as needed
}

type AuthService interface {
	GetUserID(r *http.Request) (int, error)
	LoadUser(ctx context.Context, userID int) (User, error) // Changed from `any` to `User`
}

func RequireAuth(auth AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := auth.GetUserID(r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			user, err := auth.LoadUser(r.Context(), userID)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.User, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RedirectIfAuthenticated(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app.SessionManager.Exists(r.Context(), "userID") {
				http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
