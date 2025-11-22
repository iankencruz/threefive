package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/responses"
	"github.com/iankencruz/threefive/internal/shared/session"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
)

// Context keys for storing user data in request context
type contextKey string

const (
	UserContextKey    contextKey = "user"
	SessionContextKey contextKey = "session"
)

// UserInfo represents the authenticated user info stored in context
type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// FullName returns the user's full name
func (u *UserInfo) FullName() string {
	return u.FirstName + " " + u.LastName
}

// Middleware handles authentication middleware
type Middleware struct {
	sessionManager *session.Manager
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(sessionManager *session.Manager) *Middleware {
	return &Middleware{
		sessionManager: sessionManager,
	}
}

// RequireAuth middleware ensures the user is authenticated
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionWithUser, err := m.validateSession(r)
		if err != nil {
			responses.WriteErr(w, errors.Unauthorized("Authentication required", "auth_required"))
			return
		}

		// Add user and session to context
		userInfo := &UserInfo{
			ID:        sessionWithUser.UserID,
			Email:     sessionWithUser.Email,
			FirstName: sessionWithUser.FirstName,
			LastName:  sessionWithUser.LastName,
		}

		ctx := r.Context()

		// Add values to the existing context instead of creating a new one
		ctx = context.WithValue(ctx, UserContextKey, userInfo)
		ctx = context.WithValue(ctx, SessionContextKey, sessionWithUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware validates auth if present, but doesn't require it
func (m *Middleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionWithUser, err := m.validateSession(r)
		if err == nil {
			// Add user and session to context if valid
			userInfo := &UserInfo{
				ID:        sessionWithUser.UserID,
				Email:     sessionWithUser.Email,
				FirstName: sessionWithUser.FirstName,
				LastName:  sessionWithUser.LastName,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, userInfo)
			ctx = context.WithValue(ctx, SessionContextKey, sessionWithUser)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// RequireNoAuth middleware ensures the user is NOT authenticated (for login/register pages)
func (m *Middleware) RequireNoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.validateSession(r)
		if err == nil {
			responses.WriteErr(w, errors.Conflict("Already authenticated", "already_authenticated"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware ensures the user has admin privileges
// Note: This is a placeholder - you'll need to implement role-based access control
func (m *Middleware) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			responses.WriteErr(w, errors.Unauthorized("Authentication required", "auth_required"))
			return
		}

		// TODO: Add admin role check when you implement roles
		// For now, you might want to check against specific user IDs or email domains
		// Example:
		// if !user.IsAdmin {
		//     responses.WriteErr(w, errors.Forbidden("Admin access required", "admin_required"))
		//     return
		// }

		// For now, allow all authenticated users (remove this when you add roles)
		next.ServeHTTP(w, r)
	})
}

// SessionRefresh middleware automatically refreshes sessions when needed
func (m *Middleware) SessionRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Session refresh is handled automatically in sessionManager.GetSession
		// This middleware is here for future use if we need explicit refresh logic
		next.ServeHTTP(w, r)
	})
}

// Helper methods

func (m *Middleware) validateSession(r *http.Request) (sqlc.GetSessionByTokenRow, error) {
	token := m.sessionManager.GetSessionTokenFromRequest(r)
	if token == "" {
		return sqlc.GetSessionByTokenRow{}, errors.Unauthorized("No session token", "missing_token")
	}

	sessionWithUser, err := m.sessionManager.GetSession(r.Context(), token)
	if err != nil {
		return sqlc.GetSessionByTokenRow{}, errors.Unauthorized("Invalid session", "invalid_session")
	}

	return sessionWithUser, nil
}

// Context helper functions

// GetUserFromContext retrieves the authenticated user from request context
func GetUserFromContext(ctx context.Context) *UserInfo {
	if user, ok := ctx.Value(UserContextKey).(*UserInfo); ok {
		return user
	}
	return nil
}

// GetSessionFromContext retrieves the session data from request context
func GetSessionFromContext(ctx context.Context) *sqlc.GetSessionByTokenRow {
	if session, ok := ctx.Value(SessionContextKey).(*sqlc.GetSessionByTokenRow); ok {
		return session
	}
	return nil
}

// MustGetUserFromContext retrieves the user from context or panics
// Use this only in handlers that are protected by RequireAuth middleware
func MustGetUserFromContext(ctx context.Context) *UserInfo {
	user := GetUserFromContext(ctx)
	if user == nil {
		panic("user not found in context - ensure RequireAuth middleware is applied")
	}
	return user
}

// MustGetSessionFromContext retrieves the session from context or panics
// Use this only in handlers that are protected by RequireAuth middleware
func MustGetSessionFromContext(ctx context.Context) *sqlc.GetSessionByTokenRow {
	session := GetSessionFromContext(ctx)
	if session == nil {
		panic("session not found in context - ensure RequireAuth middleware is applied")
	}
	return session
}
