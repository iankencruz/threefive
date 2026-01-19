package admin

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	// You can add dependencies here, e.g., DB *sql.DB
}

// NewHandler creates a new handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// Handler functions
func (h *Handler) HealthCheckHandler(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "up",
	})
}

func (h *Handler) HelloWorldHandler(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
