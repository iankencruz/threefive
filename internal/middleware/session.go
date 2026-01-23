// internal/middleware/session.go
package middleware

import (
	"log/slog"
	"net/http"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/internal/session"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
)

const (
	SessionContextKey = "session_data"
	UserContextKey    = "user"
)

type AuthUser struct {
	ID    string
	Email string
}

type SessionMiddleware struct {
	manager *session.SessionManager

	queries *generated.Queries
	logger  *slog.Logger
}

func NewSessionMiddleware(manager *session.SessionManager, queries *generated.Queries, logger *slog.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		manager: manager,

		queries: queries,
		logger:  logger.With("component", "session_middleware"),
	}
}

// Session middleware loads the session data and stores it in the context
func (m *SessionMiddleware) Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Request().Context()

		// Load session data
		sessionData, err := m.manager.Load(ctx, c)
		if err != nil {
			m.logger.Error("failed to load session",
				"error", err,
			)
			sessionData = make(map[string]any)
		}

		// Store session data in context
		c.Set(SessionContextKey, sessionData)

		return next(c)
	}
}

// RequireAuth ensures the user is authenticated

// RequireAuth ensures the user is authenticated
func (m *SessionMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		sessionData := GetSessionData(c)
		userID, exists := sessionData["user_id"]

		if !exists || userID == nil {
			m.logger.Debug("unauthenticated access attempt",
				"path", c.Request().URL.Path,
			)
			return c.Redirect(http.StatusFound, "/login")
		}

		// Convert string to pgtype.UUID
		userIDStr, ok := userID.(string)
		if !ok {
			m.logger.Error("invalid user_id type in session",
				"user_id", userID,
			)
			return c.Redirect(http.StatusFound, "/login")
		}

		var pgUserID pgtype.UUID
		err := pgUserID.Scan(userIDStr)
		if err != nil {
			m.logger.Error("failed to parse user_id",
				"error", err,
				"user_id", userIDStr,
			)
			return c.Redirect(http.StatusFound, "/login")
		}

		// Fetch the full user from database
		user, err := m.queries.GetUserByID(c.Request().Context(), pgUserID)
		if err != nil {
			m.logger.Error("failed to fetch user",
				"error", err,
				"user_id", userIDStr,
			)
			return c.Redirect(http.StatusFound, "/login")
		}

		// Store the full SQLC User struct in context
		c.Set(UserContextKey, &user)

		return next(c)
	}
}

// GetUser retrieves the authenticated user from context
func GetUser(c *echo.Context) *generated.User {
	user := c.Get(UserContextKey)
	if user == nil {
		return nil
	}
	return user.(*generated.User)
}

// Helper functions to work with session data

func GetSessionData(c *echo.Context) map[string]any {
	data := c.Get(SessionContextKey)
	if data == nil {
		return make(map[string]any)
	}
	return data.(map[string]any)
}

func PutSessionData(c *echo.Context, key string, value any) {
	sessionData := GetSessionData(c)
	sessionData[key] = value
	c.Set(SessionContextKey, sessionData)
}

func GetSessionValue(c *echo.Context, key string) (any, bool) {
	sessionData := GetSessionData(c)
	value, exists := sessionData[key]
	return value, exists
}

func RemoveSessionValue(c *echo.Context, key string) {
	sessionData := GetSessionData(c)
	delete(sessionData, key)
	c.Set(SessionContextKey, sessionData)
}
