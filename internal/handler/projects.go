// internal/handler/projects.go
package handler

import (
	"log/slog"
	"strconv"

	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/templates/lib"
	"github.com/iankencruz/threefive/templates/pages"
	"github.com/labstack/echo/v5"
)

type ProjectHandler struct {
	logger         *slog.Logger
	projectService *services.ProjectService
	tagService     *services.TagService
}

func NewProjectHandler(logger *slog.Logger, projectService *services.ProjectService, tagService *services.TagService) *ProjectHandler {
	return &ProjectHandler{
		logger:         logger,
		projectService: projectService,
		tagService:     tagService,
	}
}

// ShowProjectsList renders the admin projects list page
func (h *ProjectHandler) ShowProjectsList(c *echo.Context) error {
	h.logger.Debug("Loading projects list")

	// Get pagination parameters
	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := int32(20)
	offset := int32((page - 1) * int(limit))

	// Fetch projects
	projects, err := h.projectService.ListProjects(c.Request().Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list projects", "error", err)
		return c.String(500, "Failed to load projects")
	}

	// Get total count for pagination
	totalCount, err := h.projectService.CountProjects(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to count projects", "error", err)
		totalCount = 0
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	// Add user to context for template
	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
	currentPath := c.Request().URL.Path

	// Render project list page
	component := pages.ProjectsList(pages.ProjectsListProps{
		Projects:    projects,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, currentPath)

	return responses.Render(ctx, c, component)
}

// ShowCreateModal loads the create project modal
// func (h *ProjectHandler) ShowCreateModal(c *echo.Context) error {
// 	h.logger.Debug("Loading create project modal")
//
// 	// Get all tags for the tag selector
// 	tags, err := h.tagService.ListAllTags(c.Request().Context())
// 	if err != nil {
// 		h.logger.Error("failed to list tags", "error", err)
// 		tags = []generated.Tag{} // Empty list on error
// 	}
//
// 	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
//
// 	component := pages.ProjectCreateModal(tags, nil)
// 	return responses.Render(ctx, c, component)
// }
//
// // CreateProject handles project creation from modal
// func (h *ProjectHandler) CreateProject(c *echo.Context) error {
// 	h.logger.Debug("Create project request")
//
// 	// Parse form data
// 	if err := c.Request().ParseForm(); err != nil {
// 		h.logger.Error("failed to parse form", "error", err)
// 		return responses.ErrorToast(c.Request().Context(), c, "Failed to parse form data")
// 	}
//
// 	// Get authenticated user
// 	user := middleware.GetUser(c)
// 	if user == nil {
// 		return responses.ErrorToast(c.Request().Context(), c, "User not authenticated")
// 	}
//
// 	// Extract form values
// 	title := c.FormValue("title")
// 	slug := c.FormValue("slug")
// 	description := c.FormValue("description")
// 	clientName := c.FormValue("client_name")
// 	projectURL := c.FormValue("project_url")
// 	status := c.FormValue("status")
// 	projectStatus := c.FormValue("project_status")
//
// 	// Validate required fields
// 	if title == "" {
// 		tags, _ := h.tagService.ListAllTags(c.Request().Context())
// 		component := pages.ProjectCreateModal(tags, map[string]string{
// 			"title": "Title is required",
// 		})
// 		return responses.RenderError(c.Request().Context(), c, component, "Please fix the errors")
// 	}
//
// 	// Parse project year
// 	var projectYear int32
// 	if yearStr := c.FormValue("project_year"); yearStr != "" {
// 		if year, err := strconv.Atoi(yearStr); err == nil {
// 			projectYear = int32(year)
// 		}
// 	}
//
// 	// Parse project date
// 	var projectDate *time.Time
// 	if dateStr := c.FormValue("project_date"); dateStr != "" {
// 		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
// 			projectDate = &parsed
// 		}
// 	}
//
// 	// Parse featured image ID
// 	var featuredImageID *uuid.UUID
// 	if imgIDStr := c.FormValue("featured_image_id"); imgIDStr != "" {
// 		if imgUUID, err := uuid.Parse(imgIDStr); err == nil {
// 			featuredImageID = &imgUUID
// 		}
// 	}
//
// 	// Parse gallery media IDs
// 	c.Request().ParseForm()
// 	galleryIDStrs := c.Request().Form["gallery_media_ids[]"]
// 	var galleryMediaIDs []uuid.UUID
// 	for _, idStr := range galleryIDStrs {
// 		if idStr != "" {
// 			if mediaUUID, err := uuid.Parse(idStr); err == nil {
// 				galleryMediaIDs = append(galleryMediaIDs, mediaUUID)
// 			}
// 		}
// 	}
//
// 	// Parse tags (comma-separated)
// 	var tagNames []string
// 	if tagsStr := c.FormValue("tags"); tagsStr != "" {
// 		tagNames = strings.Split(tagsStr, ",")
// 		for i, tag := range tagNames {
// 			tagNames[i] = strings.TrimSpace(tag)
// 		}
// 	}
//
// 	// Get user ID
// 	var userID uuid.UUID
// 	if err := userID.Scan(user.ID.Bytes[:]); err != nil {
// 		h.logger.Error("failed to parse user ID", "error", err)
// 		return responses.ErrorToast(c.Request().Context(), c, "Failed to get user ID")
// 	}
//
// 	// Create project
// 	req := &services.CreateProjectRequest{
// 		Title:           title,
// 		Slug:            slug,
// 		Description:     description,
// 		ProjectDate:     projectDate,
// 		ClientName:      clientName,
// 		ProjectYear:     projectYear,
// 		ProjectURL:      projectURL,
// 		ProjectStatus:   projectStatus,
// 		Status:          status,
// 		FeaturedImageID: featuredImageID,
// 		AuthorID:        userID,
// 		GalleryMediaIDs: galleryMediaIDs,
// 		TagNames:        tagNames,
// 	}
//
// 	project, err := h.projectService.CreateProject(c.Request().Context(), req)
// 	if err != nil {
// 		h.logger.Error("failed to create project", "error", err)
//
// 		tags, _ := h.tagService.ListAllTags(c.Request().Context())
// 		component := pages.ProjectCreateModal(tags, map[string]string{
// 			"general": err.Error(),
// 		})
// 		return responses.RenderError(c.Request().Context(), c, component, err.Error())
// 	}
//
// 	// Success - redirect to edit page
// 	return responses.RedirectWithToast(
// 		c.Request().Context(),
// 		c,
// 		"/admin/projects/"+project.Project.Slug,
// 		"Project created successfully",
// 		"success",
// 	)
// }
//
// // ShowEditPage renders the project edit page
// func (h *ProjectHandler) ShowEditPage(c *echo.Context) error {
// 	slug := c.Param("slug")
//
// 	h.logger.Debug("Loading project edit page", "slug", slug)
//
// 	// Get project with all related data
// 	project, err := h.projectService.GetProjectBySlug(c.Request().Context(), slug)
// 	if err != nil {
// 		h.logger.Error("failed to get project", "error", err, "slug", slug)
// 		return c.String(404, "Project not found")
// 	}
//
// 	// Get all tags for the tag selector
// 	tags, err := h.tagService.ListAllTags(c.Request().Context())
// 	if err != nil {
// 		h.logger.Error("failed to list tags", "error", err)
// 		tags = []generated.Tag{}
// 	}
//
// 	// Add user to context for template
// 	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
// 	currentPath := c.Request().URL.Path
//
// 	component := pages.ProjectEditPage(project, tags, currentPath)
// 	return responses.Render(ctx, c, component)
// }
//
// // UpdateProject handles project update
// func (h *ProjectHandler) UpdateProject(c *echo.Context) error {
// 	slug := c.Param("slug")
//
// 	h.logger.Debug("Update project request", "slug", slug)
//
// 	// Get existing project
// 	existing, err := h.projectService.GetProjectBySlug(c.Request().Context(), slug)
// 	if err != nil {
// 		h.logger.Error("failed to get project", "error", err, "slug", slug)
// 		return responses.ErrorToast(c.Request().Context(), c, "Project not found")
// 	}
//
// 	// Parse form data
// 	if err := c.Request().ParseForm(); err != nil {
// 		h.logger.Error("failed to parse form", "error", err)
// 		return responses.ErrorToast(c.Request().Context(), c, "Failed to parse form data")
// 	}
//
// 	// Build update request
// 	req := &services.UpdateProjectRequest{}
//
// 	if title := c.FormValue("title"); title != "" {
// 		req.Title = &title
// 	}
// 	if newSlug := c.FormValue("slug"); newSlug != "" && newSlug != slug {
// 		req.Slug = &newSlug
// 	}
// 	if description := c.FormValue("description"); description != "" {
// 		req.Description = &description
// 	}
// 	if clientName := c.FormValue("client_name"); clientName != "" {
// 		req.ClientName = &clientName
// 	}
// 	if projectURL := c.FormValue("project_url"); projectURL != "" {
// 		req.ProjectURL = &projectURL
// 	}
// 	if status := c.FormValue("status"); status != "" {
// 		req.Status = &status
// 	}
// 	if projectStatus := c.FormValue("project_status"); projectStatus != "" {
// 		req.ProjectStatus = &projectStatus
// 	}
//
// 	// Parse project year
// 	if yearStr := c.FormValue("project_year"); yearStr != "" {
// 		if year, err := strconv.Atoi(yearStr); err == nil {
// 			year32 := int32(year)
// 			req.ProjectYear = &year32
// 		}
// 	}
//
// 	// Parse project date
// 	if dateStr := c.FormValue("project_date"); dateStr != "" {
// 		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
// 			req.ProjectDate = &parsed
// 		}
// 	}
//
// 	// Parse featured image ID
// 	if imgIDStr := c.FormValue("featured_image_id"); imgIDStr != "" {
// 		if imgUUID, err := uuid.Parse(imgIDStr); err == nil {
// 			req.FeaturedImageID = &imgUUID
// 		}
// 	}
//
// 	// Parse gallery media IDs (if provided, replaces entire gallery)
// 	c.Request().ParseForm()
// 	galleryIDStrs := c.Request().Form["gallery_media_ids[]"]
// 	if len(galleryIDStrs) > 0 {
// 		var galleryMediaIDs []uuid.UUID
// 		for _, idStr := range galleryIDStrs {
// 			if idStr != "" {
// 				if mediaUUID, err := uuid.Parse(idStr); err == nil {
// 					galleryMediaIDs = append(galleryMediaIDs, mediaUUID)
// 				}
// 			}
// 		}
// 		req.GalleryMediaIDs = galleryMediaIDs
// 	}
//
// 	// Parse tags (comma-separated, if provided replaces all tags)
// 	if tagsStr := c.FormValue("tags"); tagsStr != "" {
// 		tagNames := strings.Split(tagsStr, ",")
// 		for i, tag := range tagNames {
// 			tagNames[i] = strings.TrimSpace(tag)
// 		}
// 		req.TagNames = tagNames
// 	}
//
// 	// Update project
// 	updated, err := h.projectService.UpdateProjectBySlug(c.Request().Context(), slug, req)
// 	if err != nil {
// 		h.logger.Error("failed to update project", "error", err)
//
// 		tags, _ := h.tagService.ListAllTags(c.Request().Context())
// 		component := pages.ProjectEditPage(existing, tags, c.Request().URL.Path)
// 		return responses.RenderError(c.Request().Context(), c, component, err.Error())
// 	}
//
// 	// If slug changed, redirect to new URL
// 	if updated.Project.Slug != slug {
// 		return responses.RedirectWithToast(
// 			c.Request().Context(),
// 			c,
// 			"/admin/projects/"+updated.Project.Slug,
// 			"Project updated successfully",
// 			"success",
// 		)
// 	}
//
// 	// Success - re-render edit page with success message
// 	tags, _ := h.tagService.ListAllTags(c.Request().Context())
// 	ctx := lib.WithUser(c.Request().Context(), middleware.GetUser(c))
// 	component := pages.ProjectEditPage(updated, tags, c.Request().URL.Path)
// 	return responses.RenderSuccess(ctx, c, component, "Project updated successfully")
// }
//
// // DeleteProject soft-deletes a project
// func (h *ProjectHandler) DeleteProject(c *echo.Context) error {
// 	slug := c.Param("slug")
//
// 	h.logger.Debug("Delete project request", "slug", slug)
//
// 	if err := h.projectService.DeleteProjectBySlug(c.Request().Context(), slug); err != nil {
// 		h.logger.Error("failed to delete project", "error", err)
// 		return responses.ErrorToast(c.Request().Context(), c, "Failed to delete project")
// 	}
//
// 	// Success - return toast and let HTMX handle row removal
// 	return responses.SuccessToast(c.Request().Context(), c, "Project deleted successfully")
// }
//
// // PublishProject publishes a project
// func (h *ProjectHandler) PublishProject(c *echo.Context) error {
// 	slug := c.Param("slug")
//
// 	h.logger.Debug("Publish project request", "slug", slug)
//
// 	// Get project ID
// 	project, err := h.projectService.GetProjectBySlug(c.Request().Context(), slug)
// 	if err != nil {
// 		return responses.ErrorToast(c.Request().Context(), c, "Project not found")
// 	}
//
// 	var projectID uuid.UUID
// 	if err := projectID.Scan(project.Project.ID.Bytes[:]); err != nil {
// 		return responses.ErrorToast(c.Request().Context(), c, "Invalid project ID")
// 	}
//
// 	// Publish
// 	_, err = h.projectService.PublishProject(c.Request().Context(), projectID)
// 	if err != nil {
// 		h.logger.Error("failed to publish project", "error", err)
// 		return responses.ErrorToast(c.Request().Context(), c, "Failed to publish project")
// 	}
//
// 	return responses.SuccessToast(c.Request().Context(), c, "Project published successfully")
// }
//
// // UnpublishProject unpublishes a project
// func (h *ProjectHandler) UnpublishProject(c *echo.Context) error {
// 	slug := c.Param("slug")
//
// 	h.logger.Debug("Unpublish project request", "slug", slug)
//
// 	// Get project ID
// 	project, err := h.projectService.GetProjectBySlug(c.Request().Context(), slug)
// 	if err != nil {
// 		return responses.ErrorToast(c.Request().Context(), c, "Project not found")
// 	}
//
// 	var projectID uuid.UUID
// 	if err := projectID.Scan(project.Project.ID.Bytes[:]); err != nil {
// 		return responses.ErrorToast(c.Request().Context(), c, "Invalid project ID")
// 	}
//
// 	// Unpublish
// 	_, err = h.projectService.UnpublishProject(c.Request().Context(), projectID)
// 	if err != nil {
// 		h.logger.Error("failed to unpublish project", "error", err)
// 		return responses.ErrorToast(c.Request().Context(), c, "Failed to unpublish project")
// 	}
//
// 	return responses.SuccessToast(c.Request().Context(), c, "Project unpublished successfully")
// }
