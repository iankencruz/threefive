package middleware

import (
	"context"
	"net/http"

	"github.com/iankencruz/threefive/internal/generated"
)

// contextKey is used to avoid key collisions in request context.
type contextKey string

const userContextKey = contextKey("user")

// AuthService defines what the middleware expects from the auth handler.

type AuthService interface {
	GetUserID(r *http.Request) (int32, error)
	LoadUser(ctx context.Context, userID int32) (any, error)
}

// SessionManager defines just the part of your session manager needed for redirection.
type SessionManager interface {
	Exists(r *http.Request, key string) bool
}

// RequireAuth loads the user from session and injects them into context.
// Redirects to /login if unauthenticated.
func RequireAuth(auth AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := auth.GetUserID(r)
			if err != nil || userID == 0 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			user, err := auth.LoadUser(r.Context(), int32(userID))
			if err != nil || user == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RedirectIfAuthenticated sends logged-in users to the admin dashboard.
func RedirectIfAuthenticated(sm SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sm.Exists(r, "userID") {
				http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts the *auth.User from context.
// Returns nil if not found or invalid type.
func GetUserFromContext(r *http.Request) *generated.User {
	user, ok := r.Context().Value(userContextKey).(*generated.User)
	if !ok {
		return nil
	}
	return user
}
