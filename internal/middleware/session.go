// internal/middleware/session.go
package middleware

import (
	"log/slog"
	"net/http"

	"github.com/iankencruz/threefive/internal/session"
	"github.com/labstack/echo/v5"
)

const (
	sessionContextKey = "session_data"
)

type SessionMiddleware struct {
	manager *session.SessionManager
	logger  *slog.Logger
}

func NewSessionMiddleware(manager *session.SessionManager, logger *slog.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		manager: manager,
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
		c.Set(sessionContextKey, sessionData)

		return next(c)
	}
}

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

		// Store user_id in context for easy access
		c.Set("user_id", userID)

		return next(c)
	}
}

// Helper functions to work with session data

func GetSessionData(c *echo.Context) map[string]any {
	data := c.Get(sessionContextKey)
	if data == nil {
		return make(map[string]any)
	}
	return data.(map[string]any)
}

func PutSessionData(c *echo.Context, key string, value any) {
	sessionData := GetSessionData(c)
	sessionData[key] = value
	c.Set(sessionContextKey, sessionData)
}

func GetSessionValue(c *echo.Context, key string) (any, bool) {
	sessionData := GetSessionData(c)
	value, exists := sessionData[key]
	return value, exists
}

func RemoveSessionValue(c *echo.Context, key string) {
	sessionData := GetSessionData(c)
	delete(sessionData, key)
	c.Set(sessionContextKey, sessionData)
}
