// internal/services/projects.go
package services

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/pkg/validation"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectService struct {
	queries      *generated.Queries
	mediaService *MediaService
}

func NewProjectService(queries *generated.Queries, mediaService *MediaService) *ProjectService {
	return &ProjectService{
		queries:      queries,
		mediaService: mediaService,
	}
}

// ProjectResponse is the view model for a project with related data
type ProjectResponse struct {
	Project       generated.Project
	FeaturedImage *MediaResponse
	GalleryMedia  []MediaResponse
	Tags          []TagResponse
}

// TagResponse is the view model for a tag
type TagResponse struct {
	Tag        generated.Tag
	UsageCount int64 // Optional: how many projects use this tag
}

// CreateProjectRequest represents the data needed to create a project
type CreateProjectRequest struct {
	Title           string
	Slug            string
	Description     string
	ProjectDate     *time.Time
	ClientName      string
	ProjectYear     int32
	ProjectURL      string
	ProjectStatus   string // "completed", "in-progress", "planned"
	Status          string // "draft", "published", "archived"
	FeaturedImageID *uuid.UUID
	AuthorID        uuid.UUID
	GalleryMediaIDs []uuid.UUID // Media IDs for gallery
	TagNames        []string    // Tag names (will be created if they don't exist)
}

// Add Validate method to CreateProjectRequest
func (r *CreateProjectRequest) Validate() (validation.FieldErrors, error) {
	fields := []validation.Field{
		{
			Name:  "title",
			Value: r.Title,
			Rules: []validation.ValidationRule{
				validation.Required("Title is required"),
				validation.MinLength(3, "Title must be at least 3 characters"),
				validation.MaxLength(200, "Title must be at most 200 characters"),
			},
		},
		{
			Name:  "slug",
			Value: r.Slug,
			Rules: []validation.ValidationRule{
				validation.Required("Slug is required"),
				validation.IsSlug(""),
				validation.MaxLength(200, "Slug must be at most 200 characters"),
			},
		},
		{
			Name:  "description",
			Value: r.Description,
			Rules: []validation.ValidationRule{
				validation.MaxLength(1000, "Description must be at most 1000 characters"),
			},
		},
		{
			Name:  "project_url",
			Value: r.ProjectURL,
			Rules: []validation.ValidationRule{
				validation.IsURL(""),
			},
		},
		{
			Name:  "status",
			Value: r.Status,
			Rules: []validation.ValidationRule{
				validation.OneOf([]string{"draft", "published", "archived"}, ""),
			},
		},
		{
			Name:  "project_status",
			Value: r.ProjectStatus,
			Rules: []validation.ValidationRule{
				validation.OneOf([]string{"completed", "in-progress", "planned"}, ""),
			},
		},
	}

	errors := validation.ValidateFields(fields)
	if errors.HasErrors() {
		return errors, fmt.Errorf("validation failed")
	}

	return nil, nil
}

// CreateProject creates a new project with gallery and tags
func (s *ProjectService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectResponse, error) {
	// Validate slug
	if req.Slug == "" {
		req.Slug = GenerateSlug(req.Title)
	}

	// Check slug uniqueness
	exists, err := s.queries.CheckProjectSlugExists(ctx, generated.CheckProjectSlugExistsParams{
		Slug:      req.Slug,
		ProjectID: pgtype.UUID{Valid: false}, // Empty UUID for new project
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("slug already exists: %s", req.Slug)
	}

	projectID := uuid.New()

	// Convert dates
	var projectDate pgtype.Date
	if req.ProjectDate != nil {
		projectDate = pgtype.Date{
			Time:  *req.ProjectDate,
			Valid: true,
		}
	}

	// Convert featured image ID
	var featuredImageID pgtype.UUID
	if req.FeaturedImageID != nil {
		featuredImageID = pgtype.UUID{
			Bytes: *req.FeaturedImageID,
			Valid: true,
		}
	}

	// Create project
	_, err = s.queries.CreateProject(ctx, generated.CreateProjectParams{
		ID: pgtype.UUID{
			Bytes: projectID,
			Valid: true,
		},
		Title:           req.Title,
		Slug:            req.Slug,
		Description:     pgtype.Text{String: req.Description, Valid: req.Description != ""},
		ProjectDate:     projectDate,
		Status:          pgtype.Text{String: req.Status, Valid: req.Status != ""},
		ClientName:      pgtype.Text{String: req.ClientName, Valid: req.ClientName != ""},
		ProjectYear:     pgtype.Int4{Int32: req.ProjectYear, Valid: req.ProjectYear > 0},
		ProjectUrl:      pgtype.Text{String: req.ProjectURL, Valid: req.ProjectURL != ""},
		ProjectStatus:   pgtype.Text{String: req.ProjectStatus, Valid: req.ProjectStatus != ""},
		FeaturedImageID: featuredImageID,
		AuthorID: pgtype.UUID{
			Bytes: req.AuthorID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Add gallery media
	if len(req.GalleryMediaIDs) > 0 {
		if err := s.addGalleryMedia(ctx, projectID, req.GalleryMediaIDs); err != nil {
			return nil, fmt.Errorf("failed to add gallery media: %w", err)
		}
	}

	// Add tags
	if len(req.TagNames) > 0 {
		if err := s.addTags(ctx, projectID, req.TagNames); err != nil {
			return nil, fmt.Errorf("failed to add tags: %w", err)
		}
	}

	// Load and return full project response
	return s.GetProjectByID(ctx, projectID)
}

// UpdateProjectRequest represents the data needed to update a project
type UpdateProjectRequest struct {
	Title           string
	Slug            string
	Description     string
	ClientName      string
	ProjectYear     string
	ProjectDate     string
	ProjectURL      string
	Status          string
	ProjectStatus   string
	Tags            string
	FeaturedImageID string
	GalleryMediaIDs string
}

// Add Validate method to UpdateProjectRequest
func (r *UpdateProjectRequest) Validate() (validation.FieldErrors, error) {
	fields := []validation.Field{
		{
			Name:  "title",
			Value: r.Title,
			Rules: []validation.ValidationRule{
				validation.Required("Title is required"),
				validation.MinLength(3, "Title must be at least 3 characters"),
				validation.MaxLength(200, "Title must be at most 200 characters"),
			},
		},
		{
			Name:  "slug",
			Value: r.Slug,
			Rules: []validation.ValidationRule{
				validation.Required("Slug is required"),
				validation.IsSlug(""),
				validation.MaxLength(200, "Slug must be at most 200 characters"),
			},
		},
		{
			Name:  "description",
			Value: r.Description,
			Rules: []validation.ValidationRule{
				validation.MaxLength(1000, "Description must be at most 1000 characters"),
			},
		},
		{
			Name:  "client_name",
			Value: r.ClientName,
			Rules: []validation.ValidationRule{
				validation.MaxLength(200, "Client name must be at most 200 characters"),
			},
		},
		{
			Name:  "project_url",
			Value: r.ProjectURL,
			Rules: []validation.ValidationRule{
				validation.IsURL(""),
			},
		},
		{
			Name:  "project_year",
			Value: r.ProjectYear,
			Rules: []validation.ValidationRule{
				validation.IsYear(""),
			},
		},
		{
			Name:  "project_date",
			Value: r.ProjectDate,
			Rules: []validation.ValidationRule{
				validation.IsDate(""),
			},
		},
		{
			Name:  "status",
			Value: r.Status,
			Rules: []validation.ValidationRule{
				validation.OneOf([]string{"draft", "published", "archived"}, ""),
			},
		},
		{
			Name:  "project_status",
			Value: r.ProjectStatus,
			Rules: []validation.ValidationRule{
				validation.OneOf([]string{"completed", "in-progress", "planned"}, ""),
			},
		},
	}

	// Validate featured image ID format if provided
	if r.FeaturedImageID != "" {
		if _, err := uuid.Parse(r.FeaturedImageID); err != nil {
			errors := validation.FieldErrors{"featured_image_id": "Invalid image ID format"}
			return errors, fmt.Errorf("validation failed")
		}
	}

	// Validate gallery media IDs format if provided
	if r.GalleryMediaIDs != "" {
		ids := strings.Split(r.GalleryMediaIDs, ",")
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				if _, err := uuid.Parse(idStr); err != nil {
					errors := validation.FieldErrors{"gallery_media_ids": "Invalid gallery media ID format"}
					return errors, fmt.Errorf("validation failed")
				}
			}
		}
	}

	errors := validation.ValidateFields(fields)
	if errors.HasErrors() {
		return errors, fmt.Errorf("validation failed")
	}

	return nil, nil
}

// UpdateProjectBySlug updates a project by slug
func (s *ProjectService) UpdateProjectBySlug(ctx context.Context, slug string, req *UpdateProjectRequest) (*ProjectResponse, error) {
	// Get existing project
	existing, err := s.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if req.Slug != "" && req.Slug != slug {
		exists, err := s.queries.CheckProjectSlugExists(ctx, generated.CheckProjectSlugExistsParams{
			Slug:      req.Slug,
			ProjectID: existing.ID, // Exclude current project
		})
		if err != nil {
			return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("slug already exists: %s", req.Slug)
		}
	}

	// Build update params
	params := generated.UpdateProjectParams{
		ID: existing.ID,
	}

	if req.Title != "" {
		params.Title = pgtype.Text{String: req.Title, Valid: true}
	}
	if req.Slug != "" {
		params.Slug = pgtype.Text{String: req.Slug, Valid: true}
	}
	if req.Description != "" {
		params.Description = pgtype.Text{String: req.Description, Valid: true}
	}
	if req.ClientName != "" {
		params.ClientName = pgtype.Text{String: req.ClientName, Valid: true}
	}
	if req.ProjectURL != "" {
		params.ProjectUrl = pgtype.Text{String: req.ProjectURL, Valid: true}
	}
	if req.Status != "" {
		params.Status = pgtype.Text{String: req.Status, Valid: true}
	}
	if req.ProjectStatus != "" {
		params.ProjectStatus = pgtype.Text{String: req.ProjectStatus, Valid: true}
	}
	if req.ProjectYear != "" {
		if year, err := strconv.Atoi(req.ProjectYear); err == nil {
			params.ProjectYear = pgtype.Int4{Int32: int32(year), Valid: true}
		}
	}
	if req.ProjectDate != "" {
		if t, err := time.Parse("2006-01-02", req.ProjectDate); err == nil {
			params.ProjectDate = pgtype.Date{Time: t, Valid: true}
		}
	}

	// Handle featured image — clear if empty string, set if valid UUID
	if req.FeaturedImageID != "" {
		if imgUUID, err := uuid.Parse(req.FeaturedImageID); err == nil {
			params.FeaturedImageID = pgtype.UUID{Bytes: imgUUID, Valid: true}
		}
	} else {
		params.FeaturedImageID = pgtype.UUID{Valid: false}
	}

	// Update project core fields
	updated, err := s.queries.UpdateProject(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Sync tags
	if err := s.syncProjectTags(ctx, updated.ID, req.Tags); err != nil {
		return nil, fmt.Errorf("failed to sync tags: %w", err)
	}

	// Sync gallery media
	if err := s.syncGalleryMedia(ctx, updated.ID, req.GalleryMediaIDs); err != nil {
		return nil, fmt.Errorf("failed to sync gallery: %w", err)
	}

	return s.GetProjectBySlug(ctx, updated.Slug)
}

// ListProjects retrieves a paginated list of projects
func (s *ProjectService) ListProjects(ctx context.Context, limit, offset int32) ([]ProjectResponse, error) {
	projects, err := s.queries.ListProjects(ctx, generated.ListProjectsParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return s.buildProjectResponses(ctx, projects)
}

// GetProjectByID retrieves a project by ID with all related data
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID uuid.UUID) (*ProjectResponse, error) {
	project, err := s.queries.GetProjectByID(ctx, pgtype.UUID{
		Bytes: projectID,
		Valid: true,
	})
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	return s.buildProjectResponse(ctx, &project)
}

// GetProjectBySlug retrieves a project by slug with all related data
func (s *ProjectService) GetProjectBySlug(ctx context.Context, slug string) (*ProjectResponse, error) {
	project, err := s.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	return s.buildProjectResponse(ctx, &project)
}

// ListPublishedProjects retrieves published projects only
func (s *ProjectService) ListPublishedProjects(ctx context.Context, limit, offset int32) ([]ProjectResponse, error) {
	projects, err := s.queries.ListPublishedProjects(ctx, generated.ListPublishedProjectsParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list published projects: %w", err)
	}

	return s.buildProjectResponses(ctx, projects)
}

// DeleteProjectBySlug soft-deletes a project
func (s *ProjectService) DeleteProjectBySlug(ctx context.Context, slug string) error {
	projectID, err := s.queries.GetProjectIDBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if err := s.queries.SoftDeleteProject(ctx, projectID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// PublishProject publishes a project
func (s *ProjectService) PublishProject(ctx context.Context, projectID uuid.UUID) (*ProjectResponse, error) {
	_, err := s.queries.PublishProject(ctx, pgtype.UUID{
		Bytes: projectID,
		Valid: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to publish project: %w", err)
	}

	return s.GetProjectByID(ctx, projectID)
}

// UnpublishProject unpublishes a project
func (s *ProjectService) UnpublishProject(ctx context.Context, projectID uuid.UUID) (*ProjectResponse, error) {
	_, err := s.queries.UnpublishProject(ctx, pgtype.UUID{
		Bytes: projectID,
		Valid: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to unpublish project: %w", err)
	}

	return s.GetProjectByID(ctx, projectID)
}

// CountProjects returns the total count of projects
func (s *ProjectService) CountProjects(ctx context.Context) (int64, error) {
	return s.queries.CountProjects(ctx)
}

// Helper functions

func (s *ProjectService) buildProjectResponse(ctx context.Context, project *generated.Project) (*ProjectResponse, error) {
	response := &ProjectResponse{
		Project: *project,
	}

	// Load featured image
	if project.FeaturedImageID.Valid {
		featuredMedia, err := s.mediaService.GetMediaByID(ctx, project.FeaturedImageID)
		if err == nil {
			mediaResp := s.mediaService.ToMediaResponse(featuredMedia)
			response.FeaturedImage = &mediaResp
		}
	}

	// Load gallery media
	galleryMedia, err := s.mediaService.GetGalleryMediaForEntity(ctx, "project", project.ID)
	if err == nil && len(galleryMedia) > 0 {
		response.GalleryMedia = s.mediaService.ToMediaResponses(galleryMedia)
	}

	// Load tags
	tags, err := s.queries.GetProjectTags(ctx, project.ID)
	if err == nil && len(tags) > 0 {
		response.Tags = make([]TagResponse, len(tags))
		for i, tag := range tags {
			response.Tags[i] = TagResponse{Tag: tag}
		}
	}

	return response, nil
}

func (s *ProjectService) buildProjectResponses(ctx context.Context, projects []generated.Project) ([]ProjectResponse, error) {
	responses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		resp, err := s.buildProjectResponse(ctx, &project)
		if err != nil {
			return nil, err
		}
		responses[i] = *resp
	}
	return responses, nil
}

func (s *ProjectService) addGalleryMedia(ctx context.Context, projectID uuid.UUID, mediaIDs []uuid.UUID) error {
	for i, mediaID := range mediaIDs {
		relationID := uuid.New()
		_, err := s.queries.CreateMediaRelation(ctx, generated.CreateMediaRelationParams{
			ID: pgtype.UUID{
				Bytes: relationID,
				Valid: true,
			},
			MediaID: pgtype.UUID{
				Bytes: mediaID,
				Valid: true,
			},
			EntityType:   "project",
			EntityID:     pgtype.UUID{Bytes: projectID, Valid: true},
			RelationType: "gallery",
			SortOrder:    pgtype.Int4{Int32: int32(i), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create media relation: %w", err)
		}
	}
	return nil
}

func (s *ProjectService) addTags(ctx context.Context, projectID uuid.UUID, tagNames []string) error {
	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		tagSlug := GenerateSlug(tagName)
		tagID := uuid.New()

		// Find or create tag
		tag, err := s.queries.FindOrCreateTag(ctx, generated.FindOrCreateTagParams{
			ID: pgtype.UUID{
				Bytes: tagID,
				Valid: true,
			},
			Name: tagName,
			Slug: tagSlug,
		})
		if err != nil {
			return fmt.Errorf("failed to find or create tag: %w", err)
		}

		// Add project tag
		_, err = s.queries.AddProjectTag(ctx, generated.AddProjectTagParams{
			ProjectID: pgtype.UUID{Bytes: projectID, Valid: true},
			TagID:     tag.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to add project tag: %w", err)
		}
	}
	return nil
}

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}

// IsValidSlug checks if a slug is valid
func IsValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	// Must be lowercase, alphanumeric with hyphens only
	matched, _ := regexp.MatchString("^[a-z0-9]+(?:-[a-z0-9]+)*$", slug)
	return matched
}

// syncGalleryMedia replaces all gallery media for a project
func (s *ProjectService) syncGalleryMedia(ctx context.Context, projectID pgtype.UUID, galleryMediaIDs string) error {
	// Clear existing gallery relations
	if err := s.queries.DeleteGalleryMediaForEntity(ctx, generated.DeleteGalleryMediaForEntityParams{
		EntityType: "project",
		EntityID:   projectID,
	}); err != nil {
		return fmt.Errorf("failed to clear gallery: %w", err)
	}

	if galleryMediaIDs == "" {
		return nil
	}

	ids := strings.Split(galleryMediaIDs, ",")
	for i, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		mediaUUID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		if err := s.queries.AddGalleryMedia(ctx, generated.AddGalleryMediaParams{
			MediaID:    pgtype.UUID{Bytes: mediaUUID, Valid: true},
			EntityType: "project",
			EntityID:   projectID,
			SortOrder:  pgtype.Int4{Int32: int32(i), Valid: true},
		}); err != nil {
			return fmt.Errorf("failed to add gallery media %s: %w", idStr, err)
		}
	}

	return nil
}

// syncProjectTags replaces all tags for a project
func (s *ProjectService) syncProjectTags(ctx context.Context, projectID pgtype.UUID, tagsCSV string) error {
	// Clear existing tags
	if err := s.queries.ClearProjectTags(ctx, projectID); err != nil {
		return fmt.Errorf("failed to clear tags: %w", err)
	}

	if tagsCSV == "" {
		return nil
	}

	tagNames := strings.Split(tagsCSV, ",")
	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		tagSlug := GenerateSlug(tagName)
		tagID := uuid.New()

		tag, err := s.queries.FindOrCreateTag(ctx, generated.FindOrCreateTagParams{
			ID:   pgtype.UUID{Bytes: tagID, Valid: true},
			Name: tagName,
			Slug: tagSlug,
		})
		if err != nil {
			return fmt.Errorf("failed to find or create tag %q: %w", tagName, err)
		}

		_, err = s.queries.AddProjectTag(ctx, generated.AddProjectTagParams{
			ProjectID: projectID,
			TagID:     tag.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to add tag %q: %w", tagName, err)
		}
	}

	return nil
}
