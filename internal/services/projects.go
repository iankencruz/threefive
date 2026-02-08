// internal/services/projects.go
package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
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

// UpdateProjectRequest represents the data needed to update a project
type UpdateProjectRequest struct {
	Title           *string
	Slug            *string
	Description     *string
	ProjectDate     *time.Time
	ClientName      *string
	ProjectYear     *int32
	ProjectURL      *string
	ProjectStatus   *string
	Status          *string
	FeaturedImageID *uuid.UUID
	GalleryMediaIDs []uuid.UUID // If provided, replaces entire gallery
	TagNames        []string    // If provided, replaces all tags
}

// CreateProject creates a new project with gallery and tags
func (s *ProjectService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectResponse, error) {
	// Validate slug
	if req.Slug == "" {
		req.Slug = GenerateSlug(req.Title)
	}

	if !IsValidSlug(req.Slug) {
		return nil, fmt.Errorf("invalid slug format: must be lowercase, alphanumeric with hyphens only")
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

// UpdateProjectBySlug updates a project by slug
func (s *ProjectService) UpdateProjectBySlug(ctx context.Context, slug string, req *UpdateProjectRequest) (*ProjectResponse, error) {
	// Get project ID from slug
	projectIDResult, err := s.queries.GetProjectIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	projectID, err := uuid.FromBytes(projectIDResult.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	// Validate new slug if provided
	if req.Slug != nil && *req.Slug != slug {
		if !IsValidSlug(*req.Slug) {
			return nil, fmt.Errorf("invalid slug format")
		}

		exists, err := s.queries.CheckProjectSlugExists(ctx, generated.CheckProjectSlugExistsParams{
			Slug:      *req.Slug,
			ProjectID: projectIDResult,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("slug already exists: %s", *req.Slug)
		}
	}

	// Build update params
	updateParams := generated.UpdateProjectParams{
		ID: projectIDResult,
	}

	if req.Title != nil {
		updateParams.Title = pgtype.Text{String: *req.Title, Valid: true}
	}
	if req.Slug != nil {
		updateParams.Slug = pgtype.Text{String: *req.Slug, Valid: true}
	}
	if req.Description != nil {
		updateParams.Description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.ProjectDate != nil {
		updateParams.ProjectDate = pgtype.Date{Time: *req.ProjectDate, Valid: true}
	}
	if req.ClientName != nil {
		updateParams.ClientName = pgtype.Text{String: *req.ClientName, Valid: true}
	}
	if req.ProjectYear != nil {
		updateParams.ProjectYear = pgtype.Int4{Int32: *req.ProjectYear, Valid: true}
	}
	if req.ProjectURL != nil {
		updateParams.ProjectUrl = pgtype.Text{String: *req.ProjectURL, Valid: true}
	}
	if req.ProjectStatus != nil {
		updateParams.ProjectStatus = pgtype.Text{String: *req.ProjectStatus, Valid: true}
	}
	if req.Status != nil {
		updateParams.Status = pgtype.Text{String: *req.Status, Valid: true}
	}
	if req.FeaturedImageID != nil {
		updateParams.FeaturedImageID = pgtype.UUID{Bytes: *req.FeaturedImageID, Valid: true}
	}

	// Update project
	_, err = s.queries.UpdateProject(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Update gallery if provided
	if req.GalleryMediaIDs != nil {
		// Clear existing gallery
		if err := s.queries.DeleteAllMediaRelationsForEntity(ctx, generated.DeleteAllMediaRelationsForEntityParams{
			EntityType: "project",
			EntityID:   pgtype.UUID{Bytes: projectID, Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("failed to clear gallery: %w", err)
		}
		// Add new gallery
		if len(req.GalleryMediaIDs) > 0 {
			if err := s.addGalleryMedia(ctx, projectID, req.GalleryMediaIDs); err != nil {
				return nil, fmt.Errorf("failed to update gallery: %w", err)
			}
		}
	}

	// Update tags if provided
	if req.TagNames != nil {
		// Clear existing tags
		if err := s.queries.ClearProjectTags(ctx, pgtype.UUID{Bytes: projectID, Valid: true}); err != nil {
			return nil, fmt.Errorf("failed to clear tags: %w", err)
		}
		// Add new tags
		if len(req.TagNames) > 0 {
			if err := s.addTags(ctx, projectID, req.TagNames); err != nil {
				return nil, fmt.Errorf("failed to update tags: %w", err)
			}
		}
	}

	// Return updated project
	return s.GetProjectByID(ctx, projectID)
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
