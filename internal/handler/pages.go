// internal/handler/pages.go
package handler

import (
	"encoding/json"
	"log/slog"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/components/toast"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/templates/lib"
	"github.com/iankencruz/threefive/templates/pages"
	"github.com/jackc/pgx/v5/pgtype"
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

func (h *PageHandler) ShowPageList(c *echo.Context) error {
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

	// Build update params based on form values
	params := generated.UpdatePageParams{
		ID: existing.Page.ID,
	}

	// Common fields for all page types
	if title := c.FormValue("title"); title != "" {
		params.Title = pgtype.Text{String: title, Valid: true}
	}
	if header := c.FormValue("header"); header != "" {
		params.Header = pgtype.Text{String: header, Valid: true}
	}
	if subHeader := c.FormValue("sub_header"); subHeader != "" {
		params.SubHeader = pgtype.Text{String: subHeader, Valid: true}
	}
	if heroMediaID := c.FormValue("hero_media_id"); heroMediaID != "" {
		if mediaUUID, err := uuid.Parse(heroMediaID); err == nil {
			params.HeroMediaID = pgtype.UUID{Bytes: mediaUUID, Valid: true}
		}
	}

	// About page specific fields
	if existing.Page.PageType == "about" {
		if content := c.FormValue("content"); content != "" {
			params.Content = pgtype.Text{String: content, Valid: true}
		}
		if contentImageID := c.FormValue("content_image_id"); contentImageID != "" {
			if imgUUID, err := uuid.Parse(contentImageID); err == nil {
				params.ContentImageID = pgtype.UUID{Bytes: imgUUID, Valid: true}
			}
		}
		if ctaText := c.FormValue("cta_text"); ctaText != "" {
			params.CtaText = pgtype.Text{String: ctaText, Valid: true}
		}
		if ctaLink := c.FormValue("cta_link"); ctaLink != "" {
			params.CtaLink = pgtype.Text{String: ctaLink, Valid: true}
		}

		// Handle featured projects (up to 3)
		c.Request().ParseForm()
		featuredProjectIDs := c.Request().Form["featured_project_ids[]"]
		if len(featuredProjectIDs) > 0 {
			if err := h.pageService.UpdateFeaturedProjects(c.Request().Context(), existing.Page.ID, featuredProjectIDs); err != nil {
				h.logger.Error("failed to update featured projects", "error", err)
				return responses.ErrorToast(c.Request().Context(), c, "Failed to update featured projects")
			}
		}
	}

	// Contact page specific fields
	if existing.Page.PageType == "contact" {
		if email := c.FormValue("email"); email != "" {
			params.Email = pgtype.Text{String: email, Valid: true}
		}

		// Handle social links JSON
		// Expected form fields: social_twitter, social_linkedin, social_github, social_instagram
		socialLinks := services.SocialLinks{}
		hasAny := false

		if twitter := c.FormValue("social_twitter"); twitter != "" {
			socialLinks.Twitter = twitter
			hasAny = true
		}
		if linkedin := c.FormValue("social_linkedin"); linkedin != "" {
			socialLinks.LinkedIn = linkedin
			hasAny = true
		}
		if github := c.FormValue("social_github"); github != "" {
			socialLinks.GitHub = github
			hasAny = true
		}
		if instagram := c.FormValue("social_instagram"); instagram != "" {
			socialLinks.Instagram = instagram
			hasAny = true
		}

		if hasAny {
			socialLinksJSON, _ := json.Marshal(socialLinks)
			params.SocialLinks = socialLinksJSON
		}
	}

	h.logger.Info("Updating page", "slug", slug, "page_type", existing.Page.PageType)

	// Update page
	updated, err := h.pageService.UpdatePageBySlug(c.Request().Context(), slug, params)
	if err != nil {
		h.logger.Error("failed to update page", "error", err, "slug", slug)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to update page")
	}

	h.logger.Info("Page updated successfully", "slug", slug)

	// If slug changed, redirect to new URL
	if updated.Slug != slug {
		return responses.RedirectWithToast(c.Request().Context(), c,
			"/admin/pages/"+updated.Slug,
			"Page updated successfully",
			toast.VariantInfo,
		)
	}

	// Re-fetch enriched data for display
	pageResp, err := h.pageService.GetPageBySlug(c.Request().Context(), updated.Slug)
	if err != nil {
		return responses.ErrorToast(c.Request().Context(), c, "Page updated but failed to reload")
	}

	// Add user to context
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

	return responses.RenderSuccess(ctx, c, component, "Page updated successfully")
}

// Public page handlers (no auth required)

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
