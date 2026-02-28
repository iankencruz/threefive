// internal/handler/pages.go
package handler

import (
	"fmt"
	"log/slog"

	"github.com/a-h/templ"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/pkg/validation"
	"github.com/iankencruz/threefive/templates/lib"
	"github.com/iankencruz/threefive/templates/pages"
	"github.com/iankencruz/threefive/templates/pages/admin"
	"github.com/labstack/echo/v5"
)

type PageHandler struct {
	logger         *slog.Logger
	pageService    *services.PageService
	projectService *services.ProjectService
}

func NewPageHandler(logger *slog.Logger, pageService *services.PageService, projectService *services.ProjectService) *PageHandler {
	return &PageHandler{
		logger:         logger,
		pageService:    pageService,
		projectService: projectService,
	}
}

func (h *PageHandler) ListPages(c *echo.Context) error {
	h.logger.Debug("Loading page list")
	// Fetch all pages
	pagesResp, err := h.pageService.ListPages(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list pages", "error", err)
		return c.String(500, "Failed to load pages")
	}
	// Add user to context for template
	ctx := lib.WithUser(c.Request().Context(), getUser(c))
	currentPath := c.Request().URL.Path
	// Render page list component
	component := admin.PageList(pagesResp, currentPath)
	return responses.Render(ctx, c, component)
}

// ShowEditPage renders the edit form for a page (admin only)
func (h *PageHandler) ShowEditPage(c *echo.Context) error {
	slug := c.Param("slug")

	h.logger.Debug("Loading page edit form", "slug", slug)

	// Get page with enriched data
	pageResp, err := h.pageService.GetPageBySlug(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to get page", "error", err, "slug", slug)
		return c.String(404, "Page not found")
	}

	// Add user to context for template
	ctx := lib.WithUser(c.Request().Context(), getUser(c))
	currentPath := c.Request().URL.Path

	// Render appropriate edit form based on page type
	var component templ.Component
	switch pageResp.Page.PageType {
	case "home":
		// Get all published projects for the home page selector
		availableProjects, err := h.projectService.ListPublishedProjects(c.Request().Context(), 100, 0)
		if err != nil {
			h.logger.Error("failed to list projects", "error", err)
			availableProjects = []services.ProjectResponse{}
		}
		component = admin.AdminHome(pageResp, availableProjects, currentPath)
	case "about":
		component = admin.AdminAbout(pageResp, currentPath)
	case "contact":
		component = admin.AdminContact(pageResp, currentPath)
	default:
		return c.String(400, "Invalid page type")
	}

	return responses.Render(ctx, c, component)
}

// UpdatePage handles page update form submission
func (h *PageHandler) UpdatePage(c *echo.Context) error {
	slug := c.Param("slug")

	h.logger.Debug("Update page request", "slug", slug)

	// Get existing page to determine type
	existing, err := h.pageService.GetPageBySlug(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to get page", "error", err, "slug", slug)
		return responses.ErrorToast(c.Request().Context(), c, "Page not found")
	}

	// Parse form data
	if err := c.Request().ParseForm(); err != nil {
		h.logger.Error("failed to parse form", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to parse form data")
	}

	// DEBUG: Log all form values
	h.logger.Debug("Form data", "form", c.Request().Form)

	// Build update request
	req := &services.UpdatePageRequest{
		Title:       c.FormValue("title"),
		Header:      c.FormValue("header"),
		SubHeader:   c.FormValue("sub_header"),
		HeroMediaID: c.FormValue("hero_media_id"),
	}

	// Add page-specific fields
	if existing.Page.PageType == "home" {
		// Get featured project IDs for home page
		req.FeaturedProjectIDs = c.Request().Form["featured_project_ids[]"]
		h.logger.Debug("Home page featured projects", "ids", req.FeaturedProjectIDs, "count", len(req.FeaturedProjectIDs))

		if err := h.pageService.UpdateFeaturedProjects(c.Request().Context(), existing.Page.ID, req.FeaturedProjectIDs); err != nil {
			h.logger.Error("failed to update featured projects", "error", err)
			return responses.ErrorToast(c.Request().Context(), c, "Failed to update featured projects")
		}
	}

	if existing.Page.PageType == "about" {
		req.Content = c.FormValue("content")
		req.ContentImageID = c.FormValue("content_image_id")
		req.CtaText = c.FormValue("cta_text")
		req.CtaLink = c.FormValue("cta_link")

		// Get featured project IDs
		c.Request().ParseForm()
		featuredProjectIDs := c.Request().Form["featured_project_ids[]"]
		if err := h.pageService.UpdateFeaturedProjects(c.Request().Context(), existing.Page.ID, featuredProjectIDs); err != nil {
			h.logger.Error("failed to update featured projects", "error", err)
			return responses.ErrorToast(c.Request().Context(), c, "Failed to update featured projects")
		}
	}

	if existing.Page.PageType == "contact" {
		req.Email = c.FormValue("email")
		req.SocialTwitter = c.FormValue("social_twitter")
		req.SocialLinkedIn = c.FormValue("social_linkedin")
		req.SocialGitHub = c.FormValue("social_github")
		req.SocialInstagram = c.FormValue("social_instagram")
	}

	// ✅ VALIDATE
	fieldErrors, err := req.Validate(existing.Page.PageType)
	if err != nil {
		h.logger.Warn("validation failed", "errors", fieldErrors)

		ctx := lib.WithUser(c.Request().Context(), getUser(c))

		// Create error message with count
		errorCount := len(fieldErrors)
		toastMessage := fmt.Sprintf("Please fix %d validation error(s)", errorCount)

		// Get available projects for forms that need them
		var availableProjects []services.ProjectResponse
		if existing.Page.PageType == "home" || existing.Page.PageType == "about" {
			availableProjects, _ = h.projectService.ListPublishedProjects(c.Request().Context(), 100, 0)
		}

		// Render appropriate form based on page type with errors
		var component templ.Component
		switch existing.Page.PageType {
		case "home":
			component = admin.AdminHomeForm(existing, availableProjects, fieldErrors)
		case "about":
			component = admin.AdminAboutForm(existing, fieldErrors)
		case "contact":
			component = admin.AdminContactForm(existing, fieldErrors)
		}

		return responses.RenderError(ctx, c, component, toastMessage)
	}

	// Update page
	updated, err := h.pageService.UpdatePageBySlug(c.Request().Context(), slug, req)
	if err != nil {
		h.logger.Error("failed to update page", "error", err, "slug", slug)

		ctx := lib.WithUser(c.Request().Context(), getUser(c))

		// Database error
		dbErrors := validation.FieldErrors{"general": err.Error()}

		// Get available projects for forms that need them
		var availableProjects []services.ProjectResponse
		if existing.Page.PageType == "home" || existing.Page.PageType == "about" {
			availableProjects, _ = h.projectService.ListPublishedProjects(c.Request().Context(), 100, 0)
		}

		var component templ.Component
		switch existing.Page.PageType {
		case "home":
			component = admin.AdminHomeForm(updated, availableProjects, dbErrors)
		case "about":
			component = admin.AdminAboutForm(updated, dbErrors)
		case "contact":
			component = admin.AdminContactForm(updated, dbErrors)
		}

		return responses.RenderError(ctx, c, component, err.Error())
	}

	h.logger.Info("Page updated successfully", "slug", slug)

	// Success - render form with success toast
	ctx := lib.WithUser(c.Request().Context(), getUser(c))

	// Get available projects for forms that need them
	var availableProjects []services.ProjectResponse
	if updated.Page.PageType == "home" || updated.Page.PageType == "about" {
		availableProjects, _ = h.projectService.ListPublishedProjects(c.Request().Context(), 100, 0)
	}

	var component templ.Component
	switch updated.Page.PageType {
	case "home":
		component = admin.AdminHomeForm(updated, availableProjects, nil)
	case "about":
		component = admin.AdminAboutForm(updated, nil)
	case "contact":
		component = admin.AdminContactForm(updated, nil)
	}

	return responses.RenderSuccess(ctx, c, component, "Page updated successfully")
}

// Public page handlers (no auth required)

// ShowPublicHome renders the public home page
func (h *PageHandler) ShowPublicHome(c *echo.Context) error {
	page, err := h.pageService.GetPageBySlug(c.Request().Context(), "home")
	if err != nil {
		h.logger.Error("failed to get home page", "error", err)
		return c.String(500, "Failed to load home page")
	}

	component := pages.Home(page)
	return responses.Render(c.Request().Context(), c, component)
}

// ShowPublicAbout renders the public about page
func (h *PageHandler) ShowPublicAbout(c *echo.Context) error {
	page, err := h.pageService.GetPageBySlug(c.Request().Context(), "about")
	if err != nil {
		h.logger.Error("failed to get about page", "error", err)
		return c.String(500, "Failed to load about page")
	}

	component := pages.About(page)
	return responses.Render(c.Request().Context(), c, component)
}

// Helper function to get user from echo context
func getUser(c *echo.Context) *generated.User {
	user, ok := c.Get("user").(*generated.User)
	if !ok {
		return nil
	}
	return user
}
