package handler

import (
	"net/http"

	"github.com/iankencruz/threefive/templates/pages"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	// We'll add dependencies later (UserService, SessionStore, etc.)
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// ShowLoginPage renders the login page
func (h *AuthHandler) ShowLoginPage(c *echo.Context) error {
	props := pages.LoginPageProps{
		Email: "",
		Error: "",
	}

	// Just render the login page for now
	component := pages.LoginPage(props)
	return component.Render(c.Request().Context(), c.Response())
}

// HandleLogin will handle form submission (placeholder for now)
func (h *AuthHandler) HandleLogin(c *echo.Context) error {
	// TODO: Implement authentication logic
	// For now, just return a placeholder
	return c.String(http.StatusOK, "Login handler - to be implemented")
}
