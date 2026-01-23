package handler

import (
	"log/slog"

	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/templates/lib"
	"github.com/iankencruz/threefive/templates/pages"
	"github.com/labstack/echo/v5"
)

type AdminHandler struct {
	logger *slog.Logger
}

func NewAdminHandler(logger *slog.Logger) *AdminHandler {
	return &AdminHandler{
		logger: logger,
	}
}

// ShowDashboard renders the admin dashboard
func (h *AdminHandler) ShowDashboard(c *echo.Context) error {
	// Get authenticated user
	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))

	currentPath := c.Request().URL.Path

	// Create a context variable that inherits from a parent, and sets the value "test".

	// TODO: Get real stats from services
	// For now, use mock data
	stats := pages.DashboardStats{
		TotalProjects:     12,
		TotalBlogs:        24,
		TotalPages:        3,
		TotalMedia:        156,
		PublishedProjects: 8,
		PublishedBlogs:    18,
		DraftProjects:     4,
		DraftBlogs:        6,
	}

	component := pages.Dashboard(stats, currentPath)
	return responses.Render(ctx, c, component)
}

// ShowDashboard renders the admin dashboard
func (h *AdminHandler) ShowProjects(c *echo.Context) error {
	currentPath := c.Request().URL.Path

	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))

	// TODO: Get real stats from services
	// For now, use mock data
	stats := pages.ProjectProps{
		TotalProjects:     12,
		TotalBlogs:        24,
		TotalPages:        3,
		TotalMedia:        156,
		PublishedProjects: 8,
		PublishedBlogs:    18,
		DraftProjects:     4,
		DraftBlogs:        6,
	}

	component := pages.Projects(stats, currentPath)
	return responses.Render(ctx, c, component)
}
