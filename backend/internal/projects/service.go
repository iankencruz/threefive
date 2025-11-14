// backend/internal/projects/service.go
package projects

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service handles projects business logic
type Service struct {
	db           *pgxpool.Pool
	queries      *sqlc.Queries
	blockService *blocks.Service
}

// NewService creates a new projects service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries, blockService *blocks.Service) *Service {
	return &Service{
		db:           db,
		queries:      queries,
		blockService: blockService,
	}
}

// CreateProject creates a new project with blocks and SEO in a transaction
func (s *Service) CreateProject(ctx context.Context, req CreateProjectRequest, userID uuid.UUID) (*ProjectResponse, error) {

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
	if req.ProjectDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.ProjectDate)
		if err != nil {
			return nil, errors.BadRequest("Invalid project date format, use YYYY-MM-DD", "invalid_date")
		}
		projectDate = pgtype.Date{
			Time:  parsedDate,
			Valid: true,
		}
	}

	// 1. Create project
	project, err := qtx.CreateProject(ctx, sqlc.CreateProjectParams{
		Title:           req.Title,
		Slug:            req.Slug,
		Description:     stringToPgText(req.Description),
		ProjectDate:     projectDate,
		Status:          sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(req.Status), Valid: true},
		ClientName:      stringToPgText(req.ClientName),
		ProjectYear:     intToPgInt4(req.ProjectYear),
		ProjectUrl:      stringToPgText(req.ProjectURL),
		Technologies:    technologiesJSON,
		ProjectStatus:   sqlc.NullProjectStatus{ProjectStatus: sqlc.ProjectStatus(req.ProjectStatus), Valid: true},
		FeaturedImageID: uuidToPgUUID(req.FeaturedImageID),
		PublishedAt:     pgtype.Timestamptz{Valid: false},
	})

	if err != nil {
		return nil, errors.Internal("Failed to create project", err)
	}

	// 2. Create blocks using blocks service
	if len(req.Blocks) > 0 {
		if err := s.blockService.CreateBlocks(ctx, qtx, "project", project.ID, req.Blocks); err != nil {
			return nil, err
		}
	}

	// 3. Create SEO if provided
	if req.SEO != nil {
		if err := s.createSEO(ctx, qtx, project.ID, req.SEO); err != nil {
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

// ListProjects retrieves projects with pagination
func (s *Service) ListProjects(ctx context.Context, status *string, limit, offset int32) (*ProjectListResponse, error) {
	// Build status parameter
	var statusParam sqlc.PageStatus
	if status != nil && *status != "" {
		statusParam = sqlc.PageStatus(*status)
	}

	// Get total count
	totalCount, err := s.queries.CountProjects(ctx, statusParam)
	if err != nil {
		return nil, errors.Internal("Failed to count projects", err)
	}

	// Get projects
	projects, err := s.queries.ListProjects(ctx, sqlc.ListProjectsParams{
		Status:    statusParam,
		SortBy:    "created_at_desc",
		OffsetVal: offset,
		LimitVal:  limit,
	})
	if err != nil {
		return nil, errors.Internal("Failed to list projects", err)
	}

	// Build responses
	projectResponses := make([]ProjectResponse, 0, len(projects))
	for _, project := range projects {
		resp, err := s.buildProjectResponse(ctx, project)
		if err != nil {
			return nil, err
		}
		projectResponses = append(projectResponses, *resp)
	}

	// Calculate pagination
	totalPages := int(totalCount) / int(limit)
	if int(totalCount)%int(limit) > 0 {
		totalPages++
	}
	currentPage := int(offset)/int(limit) + 1

	return &ProjectListResponse{
		Projects: projectResponses,
		Pagination: Pagination{
			Page:       currentPage,
			Limit:      int(limit),
			TotalPages: totalPages,
			TotalCount: int(totalCount),
		},
	}, nil
}

// UpdateProject updates a project and its related data in a transaction
func (s *Service) UpdateProject(ctx context.Context, projectID uuid.UUID, req UpdateProjectRequest) (*ProjectResponse, error) {
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
		Title:           pointerToString(req.Title),
		Slug:            pointerToString(req.Slug),
		Description:     stringToPgText(req.Description),
		ProjectDate:     projectDate,
		Status:          statusToNullPageStatus(req.Status),
		ClientName:      stringToPgText(req.ClientName),
		ProjectYear:     intToPgInt4(req.ProjectYear),
		ProjectUrl:      stringToPgText(req.ProjectURL),
		Technologies:    technologiesJSON,
		ProjectStatus:   projectStatusToNullProjectStatus(req.ProjectStatus),
		FeaturedImageID: uuidToPgUUID(req.FeaturedImageID),
		PublishedAt:     pgtype.Timestamptz{Valid: false},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Project not found", "project_not_found")
		}
		return nil, errors.Internal("Failed to update project", err)
	}

	// Update blocks if provided
	if req.Blocks != nil {
		// Delete existing blocks
		if err := qtx.DeleteBlocksByEntity(ctx, sqlc.DeleteBlocksByEntityParams{
			EntityType: "project",
			EntityID:   projectID,
		}); err != nil {
			return nil, errors.Internal("Failed to delete existing blocks", err)
		}

		// Create new blocks
		if len(*req.Blocks) > 0 {
			if err := s.blockService.CreateBlocks(ctx, qtx, "project", projectID, *req.Blocks); err != nil {
				return nil, err
			}
		}
	}

	// Update SEO if provided
	if req.SEO != nil {
		if err := s.upsertSEO(ctx, qtx, projectID, req.SEO); err != nil {
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
// Helper functions
// ============================================

// buildProjectResponse builds a complete project response with all relations
func (s *Service) buildProjectResponse(ctx context.Context, project sqlc.Projects) (*ProjectResponse, error) {
	resp := &ProjectResponse{
		ID:            project.ID,
		Title:         project.Title,
		Slug:          project.Slug,
		Status:        string(project.Status.PageStatus),
		ProjectStatus: string(project.ProjectStatus.ProjectStatus),
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
	}

	// Description
	if project.Description.Valid {
		desc := project.Description.String
		resp.Description = &desc
	}

	// Project date
	if project.ProjectDate.Valid {
		dateStr := project.ProjectDate.Time.Format("2006-01-02")
		resp.ProjectDate = &dateStr
	}

	// Client name
	if project.ClientName.Valid {
		clientName := project.ClientName.String
		resp.ClientName = &clientName
	}

	// Project year
	if project.ProjectYear.Valid {
		year := int(project.ProjectYear.Int32)
		resp.ProjectYear = &year
	}

	// Project URL
	if project.ProjectUrl.Valid {
		url := project.ProjectUrl.String
		resp.ProjectURL = &url
	}

	// Technologies
	var technologies []string
	if err := json.Unmarshal(project.Technologies, &technologies); err != nil {
		technologies = []string{}
	}
	resp.Technologies = technologies

	// Featured image
	if project.FeaturedImageID.Valid {
		featuredImageID := uuid.UUID(project.FeaturedImageID.Bytes)
		resp.FeaturedImageID = &featuredImageID
	}

	// Published at
	if project.PublishedAt.Valid {
		resp.PublishedAt = &project.PublishedAt.Time
	}

	// Deleted at
	if project.DeletedAt.Valid {
		resp.DeletedAt = &project.DeletedAt.Time
	}

	// Get blocks
	blockResponses, err := s.blockService.GetBlocksByEntity(ctx, "project", project.ID)
	if err != nil {
		return nil, err
	}
	resp.Blocks = blockResponses

	// Get SEO
	seo, err := s.queries.GetSEO(ctx, sqlc.GetSEOParams{
		EntityType: "project",
		EntityID:   project.ID,
	})
	if err != nil && err != pgx.ErrNoRows {
		return nil, errors.Internal("Failed to get project SEO", err)
	}
	if err == nil {
		resp.SEO = buildSEOResponse(seo)
	}

	return resp, nil
}

// createSEO creates SEO data for a project
func (s *Service) createSEO(ctx context.Context, qtx *sqlc.Queries, projectID uuid.UUID, req *SEORequest) error {
	_, err := qtx.CreateSEO(ctx, sqlc.CreateSEOParams{
		EntityType:      "project",
		EntityID:        projectID,
		MetaTitle:       stringToPgText(req.MetaTitle),
		MetaDescription: stringToPgText(req.MetaDescription),
		OgTitle:         stringToPgText(req.OGTitle),
		OgDescription:   stringToPgText(req.OGDescription),
		OgImageID:       uuidToPgUUID(req.OGImageID),
		CanonicalUrl:    stringToPgText(req.CanonicalURL),
		RobotsIndex:     boolToPgBool(req.RobotsIndex),
		RobotsFollow:    boolToPgBool(req.RobotsFollow),
	})
	if err != nil {
		return errors.Internal("Failed to create SEO", err)
	}
	return nil
}

// upsertSEO updates or creates SEO data using the upsert query
func (s *Service) upsertSEO(ctx context.Context, qtx *sqlc.Queries, projectID uuid.UUID, req *SEORequest) error {
	_, err := qtx.UpsertSEO(ctx, sqlc.UpsertSEOParams{
		EntityType:      "project",
		EntityID:        projectID,
		MetaTitle:       stringToPgText(req.MetaTitle),
		MetaDescription: stringToPgText(req.MetaDescription),
		OgTitle:         stringToPgText(req.OGTitle),
		OgDescription:   stringToPgText(req.OGDescription),
		OgImageID:       uuidToPgUUID(req.OGImageID),
		CanonicalUrl:    stringToPgText(req.CanonicalURL),
		RobotsIndex:     boolToPgBool(req.RobotsIndex),
		RobotsFollow:    boolToPgBool(req.RobotsFollow),
	})
	if err != nil {
		return errors.Internal("Failed to upsert SEO", err)
	}
	return nil
}

// buildSEOResponse builds an SEO response from database model
func buildSEOResponse(seo sqlc.Seo) *SEOResponse {
	resp := &SEOResponse{
		ID:           seo.ID,
		RobotsIndex:  seo.RobotsIndex.Bool,
		RobotsFollow: seo.RobotsFollow.Bool,
		CreatedAt:    seo.CreatedAt,
		UpdatedAt:    seo.UpdatedAt,
	}

	if seo.MetaTitle.Valid {
		title := seo.MetaTitle.String
		resp.MetaTitle = &title
	}

	if seo.MetaDescription.Valid {
		desc := seo.MetaDescription.String
		resp.MetaDescription = &desc
	}

	if seo.OgTitle.Valid {
		ogTitle := seo.OgTitle.String
		resp.OGTitle = &ogTitle
	}

	if seo.OgDescription.Valid {
		ogDesc := seo.OgDescription.String
		resp.OGDescription = &ogDesc
	}

	if seo.OgImageID.Valid {
		imageID := uuid.UUID(seo.OgImageID.Bytes)
		resp.OGImageID = &imageID
	}

	if seo.CanonicalUrl.Valid {
		url := seo.CanonicalUrl.String
		resp.CanonicalURL = &url
	}

	return resp
}

// Helper functions for converting between types

func stringToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func intToPgInt4(i *int) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*i), Valid: true}
}

func uuidToPgUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

func boolToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Bool: true, Valid: true} // Default to true
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func pointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func statusToNullPageStatus(s *string) sqlc.NullPageStatus {
	if s == nil {
		return sqlc.NullPageStatus{Valid: false}
	}
	return sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(*s), Valid: true}
}

func projectStatusToNullProjectStatus(s *string) sqlc.NullProjectStatus {
	if s == nil {
		return sqlc.NullProjectStatus{Valid: false}
	}
	return sqlc.NullProjectStatus{ProjectStatus: sqlc.ProjectStatus(*s), Valid: true}
}
