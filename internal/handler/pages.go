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
	"github.com/iankencruz/threefive/templates/pages/admin"
	"github.com/jackc/pgx/v5/pgtype"
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

	pagesResp, err := h.pageService.ListPages(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list pages", "error", err)
		return c.String(500, "Failed to load pages")
	}

	ctx := lib.WithUser(c.Request().Context(), getUser(c))
	currentPath := c.Request().URL.Path
	component := admin.PageList(pagesResp, currentPath)
	return responses.Render(ctx, c, component)
}

// ShowEditPage renders the edit form for a page (admin only)
func (h *PageHandler) ShowEditPage(c *echo.Context) error {
	slug := c.Param("slug")
	h.logger.Debug("Loading page edit form", "slug", slug)

	pageResp, err := h.pageService.GetPageBySlug(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to get page", "error", err, "slug", slug)
		return c.String(404, "Page not found")
	}

	ctx := lib.WithUser(c.Request().Context(), getUser(c))
	currentPath := c.Request().URL.Path

	var component templ.Component
	switch pageResp.Page.PageType {
	case "home":
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

	if err := c.Request().ParseForm(); err != nil {
		h.logger.Error("failed to parse form", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to parse form data")
	}

	h.logger.Debug("Form data received", "form", c.Request().Form)

	// Build common update params
	params := generated.UpdatePageParams{
		ID: existing.Page.ID,
	}

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

	// Home page: update featured projects
	if existing.Page.PageType == "home" {
		featuredProjectIDs := c.Request().Form["featured_project_ids[]"]
		h.logger.Debug("Home page featured project IDs", "ids", featuredProjectIDs, "count", len(featuredProjectIDs))

		if err := h.pageService.UpdateFeaturedProjects(c.Request().Context(), existing.Page.ID, featuredProjectIDs); err != nil {
			h.logger.Error("failed to update featured projects", "error", err)
			return responses.ErrorToast(c.Request().Context(), c, "Failed to update featured projects")
		}
	}

	// About page specific fields + featured projects
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

		featuredProjectIDs := c.Request().Form["featured_project_ids[]"]
		if err := h.pageService.UpdateFeaturedProjects(c.Request().Context(), existing.Page.ID, featuredProjectIDs); err != nil {
			h.logger.Error("failed to update featured projects", "error", err)
			return responses.ErrorToast(c.Request().Context(), c, "Failed to update featured projects")
		}
	}

	// Contact page specific fields
	if existing.Page.PageType == "contact" {
		if email := c.FormValue("email"); email != "" {
			params.Email = pgtype.Text{String: email, Valid: true}
		}

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

	// UpdatePageBySlug now returns a full *PageResponse (enriched)
	updated, err := h.pageService.UpdatePageBySlug(c.Request().Context(), slug, params)
	if err != nil {
		h.logger.Error("failed to update page", "error", err, "slug", slug)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to update page")
	}

	h.logger.Info("Page updated successfully", "slug", slug)

	ctx := lib.WithUser(c.Request().Context(), getUser(c))

	// Build the success response component with fresh enriched data
	var component templ.Component
	switch updated.Page.PageType {
	case "home":
		availableProjects, _ := h.projectService.ListPublishedProjects(c.Request().Context(), 100, 0)
		component = admin.AdminHomeForm(updated, availableProjects, nil)
	case "about":
		component = admin.AdminAboutForm(updated, nil)
	case "contact":
		component = admin.AdminContact(updated, c.Request().URL.Path)
	default:
		return responses.ErrorToast(c.Request().Context(), c, "Unknown page type")
	}

	return responses.RenderSuccess(ctx, c, component, "Page updated successfully")
}

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

// toastVariant is a helper alias so we don't need to import toast everywhere
var _ = toast.VariantSuccess
