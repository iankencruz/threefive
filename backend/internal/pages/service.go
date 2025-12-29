// backend/internal/pages/service.go
package pages

import (
	"context"
	"log"
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

// Service handles pages business logic
type Service struct {
	db           *pgxpool.Pool
	queries      *sqlc.Queries
	blockService *blocks.Service
	config       ServiceConfig
}

type ServiceConfig struct {
	AutoPurgeRetentionDays int
}

// NewService creates a new pages service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries, blockService *blocks.Service, cfg ServiceConfig) *Service {
	return &Service{
		db:           db,
		queries:      queries,
		blockService: blockService,
		config:       cfg,
	}
}

// CreatePage creates a new page with blocks in a transaction
func (s *Service) CreatePage(ctx context.Context, req *CreatePageRequest) (*PageResponse, error) {
	// Check slug uniqueness - pass nil UUID for new pages
	var nilUUID uuid.UUID

	exists, err := s.queries.CheckSlugExists(ctx, sqlc.CheckSlugExistsParams{
		Slug:      req.Slug,
		ExcludeID: nilUUID,
	})

	// Check slug errors
	if err != nil {
		return nil, errors.Internal("Failed to check slug", err)
	}

	// check slug existence
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

	// 1. Create page
	page, err := qtx.CreatePage(ctx, sqlc.CreatePageParams{
		Title:           req.Title,
		Slug:            req.Slug,
		Status:          sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(req.Status), Valid: true},
		FeaturedImageID: utils.UUIDToPg(req.FeaturedImageID),
	})
	if err != nil {
		return nil, errors.Internal("Failed to create page", err)
	}

	// 2. Create blocks using blocks service
	if len(req.Blocks) > 0 {
		if err := s.blockService.CreateBlocks(ctx, qtx, "page", page.ID, req.Blocks); err != nil {
			return nil, err
		}
	}

	// 3. Create SEO if provided
	if req.SEO != nil {
		if err := seo.Create(ctx, qtx, "page", page.ID, req.SEO); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Internal("Failed to commit transaction", err)
	}

	// Get full page with all relations
	return s.GetPageByID(ctx, page.ID)
}

// GetPageByID retrieves a page by ID with all relations
func (s *Service) GetPageByID(ctx context.Context, pageID uuid.UUID) (*PageResponse, error) {
	// Get page
	page, err := s.queries.GetPageByID(ctx, pageID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Page not found", "page_not_found")
		}
		return nil, errors.Internal("Failed to get page", err)
	}

	return s.buildPageResponse(ctx, page)
}

// GetPageBySlug retrieves a page by slug with all relations
func (s *Service) GetPageBySlug(ctx context.Context, slug string) (*PageResponse, error) {
	// Get page
	page, err := s.queries.GetPageBySlug(ctx, slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Page not found", "page_not_found")
		}
		return nil, errors.Internal("Failed to get page", err)
	}

	return s.buildPageResponse(ctx, page)
}

// ListPages retrieves pages with pagination
func (s *Service) ListPages(ctx context.Context, params ListPagesParams) (*PageListResponse, error) {
	// Convert filters to strings
	statusStr := ""
	if params.StatusFilter != nil {
		statusStr = *params.StatusFilter
	}

	// Get Total count with filters
	totalCount, err := s.queries.CountPages(ctx, statusStr)
	if err != nil {
		return nil, errors.Internal("Failed to count pages", err)
	}

	pages, err := s.queries.ListPages(ctx, sqlc.ListPagesParams{
		Status:    statusStr,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
		OffsetVal: params.Offset,
		LimitVal:  params.Limit,
	})

	pageResponses := make([]PageResponse, 0, len(pages))
	for _, page := range pages {
		resp, err := s.buildPageResponse(ctx, page)
		if err != nil {
			return nil, err
		}
		pageResponses = append(pageResponses, *resp)
	}

	// Calculate pagination
	totalPages := int(totalCount) / int(params.Limit)
	if int(totalCount)%int(params.Limit) != 0 {
		totalPages++
	}

	currentPage := int(params.Offset)/int(params.Limit) + 1

	return &PageListResponse{
		Pages: pageResponses,
		Pagination: Pagination{
			Page:       currentPage,
			Limit:      int(params.Limit),
			TotalPages: totalPages,
			TotalCount: int(totalCount),
		},
	}, nil
}

func (s *Service) ListPublishedPages(ctx context.Context, params ListPagesParams) (*PageListResponse, error) {
	// Get Total count with filters
	totalCount, err := s.queries.CountPublishedPages(ctx)
	if err != nil {
		return nil, errors.Internal("Failed to count pages", err)
	}

	pages, err := s.queries.ListPublishedPages(ctx, sqlc.ListPublishedPagesParams{
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
		OffsetVal: params.Offset,
		LimitVal:  params.Limit,
	})
	if err != nil {
		return nil, errors.Internal("Failed to list published pages", err)
	}

	pageResponses := make([]PageResponse, 0, len(pages))
	for _, page := range pages {
		resp, err := s.buildPageResponse(ctx, page)
		if err != nil {
			return nil, err
		}
		pageResponses = append(pageResponses, *resp)
	}

	// Calculate pagination
	totalPages := int(totalCount) / int(params.Limit)
	if int(totalCount)%int(params.Limit) != 0 {
		totalPages++
	}

	currentPage := int(params.Offset)/int(params.Limit) + 1

	return &PageListResponse{
		Pages: pageResponses,
		Pagination: Pagination{
			Page:       currentPage,
			Limit:      int(params.Limit),
			TotalPages: totalPages,
			TotalCount: int(totalCount),
		},
	}, nil
}

// UpdatePage updates a page and its related data in a transaction
func (s *Service) UpdatePage(ctx context.Context, pageID uuid.UUID, req *UpdatePageRequest) (*PageResponse, error) {
	// Check slug uniqueness if slug is being updated
	if req.Slug != nil {
		exists, err := s.queries.CheckSlugExists(ctx, sqlc.CheckSlugExistsParams{
			Slug:      *req.Slug,
			ExcludeID: pageID,
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

	// Update page
	page, err := qtx.UpdatePage(ctx, sqlc.UpdatePageParams{
		ID:              pageID,
		Title:           utils.PtrStr(req.Title),
		Slug:            utils.PtrStr(req.Slug),
		Status:          utils.StatusToPg(req.Status),
		FeaturedImageID: utils.UUIDToPg(req.FeaturedImageID),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Page not found", "page_not_found")
		}
		return nil, errors.Internal("Failed to update page", err)
	}

	// Update blocks if provided
	if req.Blocks != nil {
		// Use UpdateBlocks instead of deleting and recreating
		if err := s.blockService.UpdateBlocks(ctx, qtx, "page", pageID, *req.Blocks); err != nil {
			return nil, err
		}
	}
	// Update SEO if provided
	if req.SEO != nil {
		if err := seo.Upsert(ctx, qtx, "page", page.ID, req.SEO); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Internal("Failed to commit transaction", err)
	}

	return s.GetPageByID(ctx, page.ID)
}

// UpdatePageStatus updates only the status of a page
func (s *Service) UpdatePageStatus(ctx context.Context, pageID uuid.UUID, status string) error {
	_, err := s.queries.UpdatePageStatus(ctx, sqlc.UpdatePageStatusParams{
		ID:     pageID,
		Status: sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(status), Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("Page not found", "page_not_found")
		}
		return errors.Internal("Failed to update page status", err)
	}

	return nil
}

// DeletePage soft deletes a page by modifying its slug
func (s *Service) DeletePage(ctx context.Context, pageID uuid.UUID) error {
	// Verify page exists
	_, err := s.queries.GetPageByID(ctx, pageID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("Page not found", "page_not_found")
		}
		return errors.Internal("Failed to get page", err)
	}

	// Soft delete - the query will modify the slug
	err = s.queries.SoftDeletePage(ctx, pageID)
	if err != nil {
		return errors.Internal("Failed to delete page", err)
	}

	return nil
}

// PurgeOldDeletedPages permanently deletes pages that have been soft-deleted
// for longer than the configured retention period
func (s *Service) PurgeOldDeletedPages(ctx context.Context) error {
	cutoffDate := time.Now().AddDate(0, 0, -s.config.AutoPurgeRetentionDays)

	log.Printf("[Pages] Purging pages deleted before %s (retention: %d days)",
		cutoffDate.Format("2006-01-02 15:04:05"),
		s.config.AutoPurgeRetentionDays)

	cutoffTimestamp := pgtype.Timestamptz{
		Time:  cutoffDate,
		Valid: true,
	}

	rowsDeleted, err := s.queries.PurgeOldDeletedPages(ctx, cutoffTimestamp)
	if err != nil {
		log.Printf("[Pages] Purge failed: %v", err)
		return errors.Internal("Failed to purge old deleted pages", err)
	}

	log.Printf("[Pages] Successfully purged %d page(s)", rowsDeleted)

	return nil
}

// ============================================
// Helper functions
// ============================================

// buildPageResponse builds a complete page response with all relations
func (s *Service) buildPageResponse(ctx context.Context, page sqlc.Pages) (*PageResponse, error) {
	resp := &PageResponse{
		ID:        page.ID,
		Title:     page.Title,
		Slug:      page.Slug,
		Status:    string(page.Status.PageStatus),
		CreatedAt: page.CreatedAt,
		UpdatedAt: page.UpdatedAt,
	}

	// Featured image
	if page.FeaturedImageID.Valid {
		featuredImageID := uuid.UUID(page.FeaturedImageID.Bytes)
		resp.FeaturedImageID = &featuredImageID
	}

	// Deleted at
	if page.DeletedAt.Valid {
		resp.DeletedAt = &page.DeletedAt.Time
	}

	// Get blocks
	blockResponses, err := s.blockService.GetBlocksByEntity(ctx, "page", page.ID)
	if err != nil {
		return nil, err
	}
	resp.Blocks = blockResponses

	// Get SEO
	resp.SEO, err = seo.Get(ctx, s.queries, "page", page.ID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
