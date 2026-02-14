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
	"github.com/labstack/echo/v5"
)

type PageHandler struct {
	logger      *slog.Logger
	pageService *services.PageService
}

func NewPageHandler(logger *slog.Logger, pageService *services.PageService) *PageHandler {
	return &PageHandler{
		logger:      logger,
		pageService: pageService,
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
	component := pages.PageList(pagesResp, currentPath)
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
		component = pages.AdminHome(pageResp, currentPath)
	case "about":
		component = pages.AdminAbout(pageResp, currentPath)
	case "contact":
		component = pages.AdminContact(pageResp, currentPath)
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

	// Build update request
	req := &services.UpdatePageRequest{
		Title:       c.FormValue("title"),
		Header:      c.FormValue("header"),
		SubHeader:   c.FormValue("sub_header"),
		HeroMediaID: c.FormValue("hero_media_id"),
	}

	// Add page-specific fields
	if existing.Page.PageType == "about" {
		req.Content = c.FormValue("content")
		req.ContentImageID = c.FormValue("content_image_id")
		req.CtaText = c.FormValue("cta_text")
		req.CtaLink = c.FormValue("cta_link")

		// Get featured project IDs
		c.Request().ParseForm()
		req.FeaturedProjectIDs = c.Request().Form["featured_project_ids[]"]
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

		// Render appropriate form based on page type with errors
		var component templ.Component
		switch existing.Page.PageType {
		case "home":
			component = pages.AdminHomeForm(existing, fieldErrors)
		case "about":
			component = pages.AdminAboutForm(existing, fieldErrors)
		case "contact":
			component = pages.AdminContactForm(existing, fieldErrors)
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

		var component templ.Component
		switch existing.Page.PageType {
		case "home":
			component = pages.AdminHomeForm(updated, dbErrors)
		case "about":
			component = pages.AdminAboutForm(updated, dbErrors)
		case "contact":
			component = pages.AdminContactForm(updated, dbErrors)
		}

		return responses.RenderError(ctx, c, component, err.Error())
	}

	h.logger.Info("Page updated successfully", "slug", slug)

	// Success - render form with success toast
	ctx := lib.WithUser(c.Request().Context(), getUser(c))

	var component templ.Component
	switch updated.Page.PageType {
	case "home":
		component = pages.AdminHomeForm(updated, nil)
	case "about":
		component = pages.AdminAboutForm(updated, nil)
	case "contact":
		component = pages.AdminContactForm(updated, nil)
	}

	return responses.RenderSuccess(ctx, c, component, "Page updated successfully")
} // Public page handlers (no auth required)

// ShowHomePage displays the home page to visitors
// func (h *PageHandler) ShowHomePage(c *echo.Context) error {
// 	h.logger.Debug("Loading home page")
//
// 	pageResp, err := h.pageService.GetPageBySlug(c.Request().Context(), "home")
// 	if err != nil {
// 		h.logger.Error("failed to get home page", "error", err)
// 		return c.String(500, "Failed to load page")
// 	}
//
// 	component := pages.HomePage(pageResp)
// 	return responses.Render(c.Request().Context(), c, component)
// }
//
// // ShowAboutPage displays the about page with featured projects
// func (h *PageHandler) ShowAboutPage(c *echo.Context) error {
// 	h.logger.Debug("Loading about page")
//
// 	pageResp, err := h.pageService.GetPageBySlug(c.Request().Context(), "about")
// 	if err != nil {
// 		h.logger.Error("failed to get about page", "error", err)
// 		return c.String(500, "Failed to load page")
// 	}
//
// 	component := pages.AboutPage(pageResp)
// 	return responses.Render(c.Request().Context(), c, component)
// }
//
// // ShowContactPage displays the contact page
// func (h *PageHandler) ShowContactPage(c *echo.Context) error {
// 	h.logger.Debug("Loading contact page")
//
// 	pageResp, err := h.pageService.GetPageBySlug(c.Request().Context(), "contact")
// 	if err != nil {
// 		h.logger.Error("failed to get contact page", "error", err)
// 		return c.String(500, "Failed to load page")
// 	}
//
// 	component := pages.ContactPage(pageResp)
// 	return responses.Render(c.Request().Context(), c, component)
// }

// Helper function to get user from echo context
func getUser(c *echo.Context) *generated.User {
	user, ok := c.Get("user").(*generated.User)
	if !ok {
		return nil
	}
	return user
}
