// internal/handler/auth.go
package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/internal/session"
	"github.com/iankencruz/threefive/pkg/errors"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/templates/pages"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService    *services.AuthService
	sessionManager *session.SessionManager
	logger         *slog.Logger
}

func NewAuthHandler(authService *services.AuthService, sessionManager *session.SessionManager, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionManager: sessionManager,
		logger:         logger.With("component", "auth_handler"),
	}
}

// ShowLoginPage renders the login page
func (h *AuthHandler) ShowLoginPage(c *echo.Context) error {
	// Check if already authenticated
	sessionData := middleware.GetSessionData(c)
	if _, exists := sessionData["user_id"]; exists {
		h.logger.Debug("user already authenticated, redirecting")
		return responses.Redirect(c, "/admin")
	}

	props := pages.LoginPageProps{
		Email: "",
		Error: "",
	}

	component := pages.LoginPage(props)
	return responses.Render(c.Request().Context(), c, component)
}

// HandleLogin processes login form submission
func (h *AuthHandler) HandleLogin(c *echo.Context) error {
	ctx := c.Request().Context()
	email := c.Request().FormValue("email")
	password := c.Request().FormValue("password")

	h.logger.Info("login attempt",
		"email", email,
	)

	// Validate input
	if email == "" || password == "" {
		h.logger.Warn("login validation failed - empty fields",
			"email", email,
		)

		props := pages.LoginPageProps{
			Email: email,
			Error: "Email and password are required",
		}
		component := pages.LoginPage(props)
		return responses.RenderWithStatus(ctx, c, http.StatusBadRequest, component)
	}

	// Authenticate user
	user, err := h.authService.Authenticate(ctx, email, password)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Warn("authentication failed",
				"email", email,
				"error", appErr.Message,
			)

			props := pages.LoginPageProps{
				Email: email,
				Error: appErr.Message,
			}
			component := pages.LoginPage(props)
			return responses.RenderWithStatus(ctx, c, appErr.Code, component)
		}

		h.logger.Error("unexpected authentication error",
			"email", email,
			"error", err,
		)

		props := pages.LoginPageProps{
			Email: email,
			Error: "An unexpected error occurred",
		}
		component := pages.LoginPage(props)
		return responses.RenderWithStatus(ctx, c, http.StatusInternalServerError, component)
	}

	userIDStr := uuid.UUID(user.ID.Bytes).String()

	// Create new session with user_id
	sessionData := middleware.GetSessionData(c)

	sessionData["user_id"] = userIDStr
	sessionData["user_email"] = user.Email

	err = h.sessionManager.Save(ctx, c, sessionData)
	if err != nil {
		h.logger.Error("failed to save session",
			"user_id", userIDStr,
			"error", err,
		)

		props := pages.LoginPageProps{
			Email: email,
			Error: "Failed to create session. Please try again.",
		}
		component := pages.LoginPage(props)
		return responses.RenderWithStatus(ctx, c, http.StatusInternalServerError, component)
	}

	h.logger.Info("user logged in successfully",
		"email", email,
		"user_id", userIDStr,
	)

	// Redirect to home/dashboard
	return responses.HTMXRedirect(c, "/admin")
}

// HandleLogout destroys the session and redirects to login
func (h *AuthHandler) HandleLogout(c *echo.Context) error {
	ctx := c.Request().Context()

	sessionData := middleware.GetSessionData(c)
	if userID, exists := sessionData["user_id"]; exists {
		h.logger.Info("user logging out",
			"user_id", userID,
		)
	}

	// Destroy session
	err := h.sessionManager.Destroy(ctx, c)
	if err != nil {
		h.logger.Error("failed to destroy session",
			"error", err,
		)
	}

	h.logger.Debug("session destroyed, redirecting to login")

	// Redirect to login
	return responses.Redirect(c, "/login")
}
