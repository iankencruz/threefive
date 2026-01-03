// backend/internal/projects/service.go
package projects

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/seo"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

// NewService creates a new projects service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries) *Service {
	return &Service{
		db:      db,
		queries: queries,
	}
}

// CreateProject creates a new project with media and SEO
func (s *Service) CreateProject(ctx context.Context, req *CreateProjectRequest, userID uuid.UUID) (*ProjectResponse, error) {
	// Validate slug uniqueness
	exists, err := s.queries.CheckProjectSlugExists(ctx, sqlc.CheckProjectSlugExistsParams{
		Slug:      req.Slug,
		ExcludeID: uuid.Nil,
	})
	if err != nil {
		return nil, errors.Internal("Failed to check slug", err)
	}
	if exists {
		return nil, errors.Conflict("Slug already exists", "slug_exists")
	}

	// Validate featured_image_id is in media_ids if provided
	if req.FeaturedImageID != nil && len(req.MediaIDs) > 0 {
		featuredFound := false
		for _, mid := range req.MediaIDs {
			if mid == *req.FeaturedImageID {
				featuredFound = true
				break
			}
		}
		if !featuredFound {
			return nil, errors.BadRequest("featured_image_id must be one of the media_ids", "invalid_featured_image")
		}
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errors.Internal("Failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Prepare technologies JSON
	technologiesJSON, err := json.Marshal(req.Technologies)
	if err != nil {
		return nil, errors.Internal("Failed to marshal technologies", err)
	}

	// Parse project date if provided
	var projectDate pgtype.Date
	if req.ProjectDate != nil && *req.ProjectDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.ProjectDate)
		if err != nil {
			return nil, errors.BadRequest("Invalid project date format, use YYYY-MM-DD", "invalid_date")
		}
		projectDate = pgtype.Date{
			Time:  parsedDate,
			Valid: true,
		}
	}

	// Create project
	project, err := qtx.CreateProject(ctx, sqlc.CreateProjectParams{
		Title:           req.Title,
		Slug:            req.Slug,
		Description:     utils.StrToPg(req.Description),
		ProjectDate:     projectDate,
		Status:          sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(req.Status), Valid: true},
		ClientName:      utils.StrToPg(req.ClientName),
		ProjectYear:     utils.IntToPg(req.ProjectYear),
		ProjectUrl:      utils.StrToPg(req.ProjectURL),
		Technologies:    technologiesJSON,
		ProjectStatus:   sqlc.NullProjectStatus{ProjectStatus: sqlc.ProjectStatus(req.ProjectStatus), Valid: true},
		FeaturedImageID: utils.UUIDToPg(req.FeaturedImageID),
		PublishedAt:     pgtype.Timestamptz{Valid: false},
	})
	if err != nil {
		return nil, errors.Internal("Failed to create project", err)
	}

	// Link media to project
	if len(req.MediaIDs) > 0 {
		for i, mediaID := range req.MediaIDs {
			_, err := qtx.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
				MediaID:    mediaID,
				EntityType: "project",
				EntityID:   project.ID,
				SortOrder:  pgtype.Int4{Int32: int32(i), Valid: true},
			})
			if err != nil {
				return nil, errors.Internal("Failed to link media to project", err)
			}
		}
	}

	// Create SEO if provided
	if req.SEO != nil {
		if err := seo.Create(ctx, qtx, "project", project.ID, req.SEO); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Internal("Failed to commit transaction", err)
	}

	// Get full project with all relations
	return s.GetProjectByID(ctx, project.ID)
}

// GetProjectByID retrieves a project by ID with all relations
func (s *Service) GetProjectByID(ctx context.Context, projectID uuid.UUID) (*ProjectResponse, error) {
	// Get project
	project, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Project not found", "project_not_found")
		}
		return nil, errors.Internal("Failed to get project", err)
	}

	return s.buildProjectResponse(ctx, project)
}

// GetProjectBySlug retrieves a project by slug with all relations
func (s *Service) GetProjectBySlug(ctx context.Context, slug string) (*ProjectResponse, error) {
	// Get project
	project, err := s.queries.GetProjectBySlug(ctx, slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Project not found", "project_not_found")
		}
		return nil, errors.Internal("Failed to get project", err)
	}

	return s.buildProjectResponse(ctx, project)
}

// ListProjects retrieves projects with optional filtering and pagination
func (s *Service) ListProjects(ctx context.Context, params ListProjectsParams) (*ProjectListResponse, error) {
	// Convert filters to strings (empty string = no filter)
	statusStr := ""
	if params.StatusFilter != nil {
		statusStr = *params.StatusFilter
	}

	// Get total count with filters
	totalCount, err := s.queries.CountProjects(ctx, statusStr)
	if err != nil {
		return nil, errors.Internal("Failed to count projects", err)
	}

	// Get projects with filters
	projects, err := s.queries.ListProjects(ctx, sqlc.ListProjectsParams{
		Status:    statusStr,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
		OffsetVal: params.Offset,
		LimitVal:  params.Limit,
	})
	if err != nil {
		return nil, errors.Internal("Failed to list projects", err)
	}

	// Build response for each project
	projectResponses := make([]ProjectResponse, 0, len(projects))
	for _, proj := range projects {
		response, err := s.buildProjectResponse(ctx, proj)
		if err != nil {
			return nil, err
		}
		projectResponses = append(projectResponses, *response)
	}

	// Calculate pagination
	totalPages := int(totalCount) / int(params.Limit)
	if int(totalCount)%int(params.Limit) > 0 {
		totalPages++
	}
	currentPage := int(params.Offset)/int(params.Limit) + 1

	return &ProjectListResponse{
		Projects: projectResponses,
		Pagination: Pagination{
			Page:       currentPage,
			Limit:      int(params.Limit),
			TotalPages: totalPages,
			TotalCount: int(totalCount),
		},
	}, nil
}

// UpdateProject updates a project and its related data in a transaction
func (s *Service) UpdateProject(ctx context.Context, projectID uuid.UUID, req *UpdateProjectRequest) (*ProjectResponse, error) {
	// Check slug uniqueness if slug is being updated
	if req.Slug != nil {
		exists, err := s.queries.CheckProjectSlugExists(ctx, sqlc.CheckProjectSlugExistsParams{
			Slug:      *req.Slug,
			ExcludeID: projectID,
		})
		if err != nil {
			return nil, errors.Internal("Failed to check slug", err)
		}
		if exists {
			return nil, errors.Conflict("Slug already exists", "slug_exists")
		}
	}

	// Validate featured_image_id is in media_ids if both provided
	if req.FeaturedImageID != nil && req.MediaIDs != nil && len(*req.MediaIDs) > 0 {
		featuredFound := false
		for _, mid := range *req.MediaIDs {
			if mid == *req.FeaturedImageID {
				featuredFound = true
				break
			}
		}
		if !featuredFound {
			return nil, errors.BadRequest("featured_image_id must be one of the media_ids", "invalid_featured_image")
		}
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errors.Internal("Failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Prepare technologies JSON if provided
	var technologiesJSON []byte
	if req.Technologies != nil {
		technologiesJSON, err = json.Marshal(*req.Technologies)
		if err != nil {
			return nil, errors.Internal("Failed to marshal technologies", err)
		}
	}

	// Parse project date if provided
	var projectDate pgtype.Date
	if req.ProjectDate != nil {
		if *req.ProjectDate != "" {
			parsedDate, err := time.Parse("2006-01-02", *req.ProjectDate)
			if err != nil {
				return nil, errors.BadRequest("Invalid project date format, use YYYY-MM-DD", "invalid_date")
			}
			projectDate = pgtype.Date{
				Time:  parsedDate,
				Valid: true,
			}
		}
	}

	// Update project
	_, err = qtx.UpdateProject(ctx, sqlc.UpdateProjectParams{
		ID:              projectID,
		Title:           utils.PtrStr(req.Title),
		Slug:            utils.PtrStr(req.Slug),
		Description:     utils.StrToPg(req.Description),
		ProjectDate:     projectDate,
		Status:          utils.StatusToPg(req.Status),
		ClientName:      utils.StrToPg(req.ClientName),
		ProjectYear:     utils.IntToPg(req.ProjectYear),
		ProjectUrl:      utils.StrToPg(req.ProjectURL),
		Technologies:    technologiesJSON,
		ProjectStatus:   utils.ProjectStatusToPg(req.ProjectStatus),
		FeaturedImageID: utils.UUIDToPg(req.FeaturedImageID),
		PublishedAt:     pgtype.Timestamptz{Valid: false},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Project not found", "project_not_found")
		}
		return nil, errors.Internal("Failed to update project", err)
	}

	// Update media if provided
	if req.MediaIDs != nil {
		// Get existing media relations
		existingMedia, err := qtx.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
			EntityType: "project",
			EntityID:   projectID,
		})
		if err != nil {
			return nil, errors.Internal("Failed to get existing media", err)
		}

		// Create map of existing media IDs
		existingMap := make(map[uuid.UUID]bool)
		for _, m := range existingMedia {
			existingMap[m.ID] = true
		}

		// Create map of new media IDs
		newMap := make(map[uuid.UUID]bool)
		for _, mid := range *req.MediaIDs {
			newMap[mid] = true
		}

		// Remove media that are no longer in the list
		for _, m := range existingMedia {
			if !newMap[m.ID] {
				err := qtx.UnlinkMediaFromEntity(ctx, sqlc.UnlinkMediaFromEntityParams{
					MediaID:    m.ID,
					EntityType: "project",
					EntityID:   projectID,
				})
				if err != nil {
					return nil, errors.Internal("Failed to unlink media", err)
				}
			}
		}

		// Add/update new media with proper sort order
		for i, mediaID := range *req.MediaIDs {
			_, err := qtx.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
				MediaID:    mediaID,
				EntityType: "project",
				EntityID:   projectID,
				SortOrder:  pgtype.Int4{Int32: int32(i), Valid: true},
			})
			if err != nil {
				return nil, errors.Internal("Failed to link media to project", err)
			}
		}
	}

	// Update SEO if provided
	if req.SEO != nil {
		if err := seo.Upsert(ctx, qtx, "project", projectID, req.SEO); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Internal("Failed to commit transaction", err)
	}

	// Get updated project with all relations
	return s.GetProjectByID(ctx, projectID)
}

// UpdateProjectStatus updates only the status of a project
func (s *Service) UpdateProjectStatus(ctx context.Context, projectID uuid.UUID, status string) error {
	_, err := s.queries.UpdateProjectStatus(ctx, sqlc.UpdateProjectStatusParams{
		ID:     projectID,
		Status: sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(status), Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("Project not found", "project_not_found")
		}
		return errors.Internal("Failed to update project status", err)
	}
	return nil
}

// DeleteProject soft deletes a project
func (s *Service) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	err := s.queries.SoftDeleteProject(ctx, projectID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("Project not found", "project_not_found")
		}
		return errors.Internal("Failed to delete project", err)
	}
	return nil
}

// ============================================
// Helper Functions
// ============================================

// buildProjectResponse builds a complete ProjectResponse with all related data
func (s *Service) buildProjectResponse(ctx context.Context, project sqlc.Projects) (*ProjectResponse, error) {
	// Parse technologies
	var technologies []string
	if len(project.Technologies) > 0 {
		if err := json.Unmarshal(project.Technologies, &technologies); err != nil {
			return nil, errors.Internal("Failed to unmarshal technologies", err)
		}
	}

	// Get media for project
	media, err := s.queries.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
		EntityType: "project",
		EntityID:   project.ID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get project media", err)
	}

	// Convert media to MediaItem
	mediaItems := make([]MediaItem, 0, len(media))
	for i, m := range media {
		mediaItems = append(mediaItems, MediaItem{
			ID:               m.ID,
			Filename:         m.Filename,
			OriginalFilename: m.OriginalFilename,
			MimeType:         m.MimeType,
			URL:              *utils.PgToStr(m.Url),
			ThumbnailURL:     *utils.PgToStr(m.ThumbnailUrl),
			MediumURL:        *utils.PgToStr(m.MediumUrl),
			LargeURL:         *utils.PgToStr(m.LargeUrl),
			Width:            utils.PgToInt(m.Width),
			Height:           utils.PgToInt(m.Height),
			SizeBytes:        m.SizeBytes,
			SortOrder:        i,
		})
	}

	// Get featured image if set
	var featuredImage *MediaItem
	if project.FeaturedImageID.Valid {
		featuredID := uuid.UUID(project.FeaturedImageID.Bytes)
		featuredMedia, err := s.queries.GetMediaByID(ctx, featuredID)
		if err == nil {
			featuredImage = &MediaItem{
				ID:               featuredMedia.ID,
				Filename:         featuredMedia.Filename,
				OriginalFilename: featuredMedia.OriginalFilename,
				MimeType:         featuredMedia.MimeType,
				URL:              *utils.PgToStr(featuredMedia.Url),
				ThumbnailURL:     *utils.PgToStr(featuredMedia.ThumbnailUrl),
				MediumURL:        *utils.PgToStr(featuredMedia.MediumUrl),
				LargeURL:         *utils.PgToStr(featuredMedia.LargeUrl),
				Width:            utils.PgToInt(featuredMedia.Width),
				Height:           utils.PgToInt(featuredMedia.Height),
				SizeBytes:        featuredMedia.SizeBytes,
			}
		}
	}

	// Get SEO
	var seoResponse *seo.Response
	seoData, err := seo.Get(ctx, s.queries, "project", project.ID)
	if err == nil && seoData != nil {
		seoResponse = seoData
	}

	// Format project date
	var projectDate *string
	if project.ProjectDate.Valid {
		formatted := project.ProjectDate.Time.Format("2006-01-02")
		projectDate = &formatted
	}

	// Format published_at
	var publishedAt *time.Time
	if project.PublishedAt.Valid {
		publishedAt = &project.PublishedAt.Time
	}

	return &ProjectResponse{
		ID:              project.ID,
		Title:           project.Title,
		Slug:            project.Slug,
		Description:     utils.PgToStr(project.Description),
		ProjectDate:     projectDate,
		Status:          string(project.Status.PageStatus),
		ClientName:      utils.PgToStr(project.ClientName),
		ProjectYear:     utils.PgToInt(project.ProjectYear),
		ProjectURL:      utils.PgToStr(project.ProjectUrl),
		Technologies:    technologies,
		ProjectStatus:   string(project.ProjectStatus.ProjectStatus),
		FeaturedImageID: utils.PgToUUID(project.FeaturedImageID),
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		PublishedAt:     publishedAt,
		Media:           mediaItems,
		FeaturedImage:   featuredImage,
		SEO:             seoResponse,
	}, nil
}
