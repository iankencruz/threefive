// backend/internal/projects/service.go
package projects

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/seo"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/utils"
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

	// 2. Create blocks using blocks service
	if len(req.Blocks) > 0 {
		if err := s.blockService.CreateBlocks(ctx, qtx, "project", project.ID, req.Blocks); err != nil {
			return nil, err
		}
	}

	// 3. Create SEO if provided
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
	resp.SEO, err = seo.Get(ctx, s.queries, "project", project.ID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
