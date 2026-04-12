// internal/handler/tags.go
package handler

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/templates/components/toast"
	"github.com/iankencruz/threefive/templates/lib"
	"github.com/iankencruz/threefive/templates/pages/admin"
	"github.com/labstack/echo/v5"
)

type TagHandler struct {
	logger     *slog.Logger
	tagService *services.TagService
}

func NewTagHandler(logger *slog.Logger, tagService *services.TagService) *TagHandler {
	return &TagHandler{
		logger:     logger,
		tagService: tagService,
	}
}

// ShowTagsList renders the admin tags list page
func (h *TagHandler) ShowTagsList(c *echo.Context) error {
	h.logger.Debug("Loading tags list")

	// Get search parameter
	search := c.QueryParam("search")

	// Get pagination parameters
	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := int32(20)
	offset := int32((page - 1) * int(limit))

	// Fetch tags with usage counts (with search if provided)
	var tagsWithUsage []services.TagResponse
	var err error

	if search != "" {
		tagsWithUsage, err = h.tagService.SearchTagsWithUsage(c.Request().Context(), search, limit, offset)
	} else {
		tagsWithUsage, err = h.tagService.ListTagsWithUsage(c.Request().Context(), limit, offset)
	}

	if err != nil {
		h.logger.Error("failed to list tags", "error", err)
		return c.String(500, "Failed to load tags")
	}

	// Get total count for pagination
	totalCount, err := h.tagService.CountTags(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to count tags", "error", err)
		totalCount = 0
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	// Add user to context for template
	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
	currentPath := c.Request().URL.Path

	// Render tags list page
	component := admin.TagsList(admin.TagsListProps{
		Tags:        tagsWithUsage,
		TotalCount:  int(totalCount),
		CurrentPage: page,
		TotalPages:  totalPages,
		SearchTerm:  search,
	}, currentPath)

	return responses.Render(ctx, c, component)
}

// ShowCreateModal shows the create tag modal
func (h *TagHandler) ShowCreateModal(c *echo.Context) error {
	h.logger.Debug("Loading create tag modal")

	ctx := c.Request().Context()
	component := lib.TagCreateModal(nil)

	return responses.Render(ctx, c, component)
}

// CreateTag handles tag creation from the modal
func (h *TagHandler) CreateTag(c *echo.Context) error {
	h.logger.Debug("Creating tag")

	// Parse form
	if err := c.Request().ParseForm(); err != nil {
		h.logger.Error("failed to parse form", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Invalid form data")
	}

	req := &services.CreateTagRequest{
		Name: c.Request().FormValue("name"),
		Slug: c.Request().FormValue("slug"),
	}

	// Validate
	if errs := req.Validate(); len(errs) > 0 {
		h.logger.Warn("Tag validation failed", "errors", errs)
		component := lib.TagCreateModal(errs)
		return responses.RenderError(c.Request().Context(), c, component, "Please fix the errors")
	}

	// Create tag
	tag, err := h.tagService.CreateTag(c.Request().Context(), req)
	if err != nil {
		h.logger.Error("failed to create tag", "error", err)
		component := lib.TagCreateModal(map[string]string{
			"general": err.Error(),
		})
		return responses.RenderError(c.Request().Context(), c, component, err.Error())
	}

	h.logger.Info("Tag created successfully", "slug", tag.Slug)

	// Success - close modal and redirect to edit page
	return responses.RedirectWithToast(
		c.Request().Context(),
		c,
		"/admin/tags/"+tag.Slug,
		"Tag created successfully",
		toast.VariantSuccess,
	)
}

// ShowEditPage shows the tag edit page
func (h *TagHandler) ShowEditPage(c *echo.Context) error {
	slug := c.Param("slug")
	h.logger.Debug("Loading tag edit page", "slug", slug)

	tagWithUsage, err := h.tagService.GetTagBySlugWithUsage(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to get tag", "slug", slug, "error", err)
		return c.String(404, "Tag not found")
	}

	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
	currentPath := c.Request().URL.Path

	component := admin.TagEdit(tagWithUsage, nil, currentPath)
	return responses.Render(ctx, c, component)
}

// UpdateTag handles tag updates
func (h *TagHandler) UpdateTag(c *echo.Context) error {
	slug := c.Param("slug")
	h.logger.Debug("Updating tag", "slug", slug)

	// Get existing tag first
	existing, err := h.tagService.GetTagBySlugWithUsage(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to get existing tag", "slug", slug, "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Tag not found")
	}

	// Parse form
	if err := c.Request().ParseForm(); err != nil {
		h.logger.Error("failed to parse form", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Invalid form data")
	}

	req := &services.UpdateTagRequest{
		Name: c.Request().FormValue("name"),
		Slug: c.Request().FormValue("slug"),
	}

	// Validate
	if errs := req.Validate(); len(errs) > 0 {
		h.logger.Warn("Tag validation failed", "errors", errs)
		ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
		currentPath := c.Request().URL.Path
		component := admin.TagEdit(existing, errs, currentPath)
		return responses.RenderError(ctx, c, component, "Please fix the errors")
	}

	// Update tag
	updated, err := h.tagService.UpdateTagBySlug(c.Request().Context(), slug, req)
	if err != nil {
		h.logger.Error("failed to update tag", "error", err)
		ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
		currentPath := c.Request().URL.Path
		component := admin.TagEdit(existing, map[string]string{
			"general": err.Error(),
		}, currentPath)
		return responses.RenderError(ctx, c, component, err.Error())
	}

	h.logger.Info("Tag updated successfully", "slug", updated.Slug)

	// If slug changed, redirect to new URL
	if updated.Slug != slug {
		return responses.RedirectWithToast(
			c.Request().Context(),
			c,
			"/admin/tags/"+updated.Slug,
			"Tag updated successfully",
			toast.VariantSuccess,
		)
	}

	// Normal update - show success message
	// Get updated tag with usage
	updatedWithUsage, _ := h.tagService.GetTagBySlugWithUsage(c.Request().Context(), updated.Slug)
	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
	currentPath := c.Request().URL.Path
	component := admin.TagEdit(updatedWithUsage, nil, currentPath)
	return responses.RenderSuccess(ctx, c, component, "Tag updated successfully")
}

// DeleteTag handles tag deletion
func (h *TagHandler) DeleteTag(c *echo.Context) error {
	slug := c.Param("slug")
	h.logger.Debug("Deleting tag", "slug", slug)

	err := h.tagService.DeleteTagBySlug(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to delete tag", "slug", slug, "error", err)
		return responses.ErrorToast(c.Request().Context(), c, err.Error())
	}

	h.logger.Info("Tag deleted successfully", "slug", slug)

	// Return success toast - row will be removed via HTMX
	return responses.SuccessToast(c.Request().Context(), c, "Tag deleted successfully")
}

// ShowUnusedTags shows a list of unused tags
func (h *TagHandler) ShowUnusedTags(c *echo.Context) error {
	h.logger.Debug("Loading unused tags")

	tags, err := h.tagService.GetUnusedTags(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to get unused tags", "error", err)
		return c.String(500, "Failed to load unused tags")
	}

	ctx := c.Request().Context()
	component := lib.UnusedTagsList(tags)

	return responses.Render(ctx, c, component)
}

// DeleteUnusedTags bulk deletes all unused tags
func (h *TagHandler) DeleteUnusedTags(c *echo.Context) error {
	h.logger.Debug("Deleting unused tags")

	deleted, err := h.tagService.DeleteUnusedTags(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to delete unused tags", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to delete unused tags")
	}

	h.logger.Info("Deleted unused tags", "count", deleted)

	if deleted == 0 {
		return responses.ToastOnly(c.Request().Context(), c, "No unused tags to delete", toast.VariantInfo)
	}

	// Redirect back to tags list with success message
	return responses.RedirectWithToast(
		c.Request().Context(),
		c,
		"/admin/tags",
		fmt.Sprintf("Deleted %d unused tag(s)", deleted),
		toast.VariantSuccess,
	)
}
