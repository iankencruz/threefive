// backend/internal/pages/service.go
package pages

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

// Service handles pages business logic
type Service struct {
	db           *pgxpool.Pool
	queries      *sqlc.Queries
	blockService *blocks.Service
}

// NewService creates a new pages service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries, blockService *blocks.Service) *Service {
	return &Service{
		db:           db,
		queries:      queries,
		blockService: blockService,
	}
}

// CreatePage creates a new page with blocks in a transaction
func (s *Service) CreatePage(ctx context.Context, req CreatePageRequest, userID uuid.UUID) (*PageResponse, error) {
	// Check slug uniqueness
	exists, err := s.queries.CheckSlugExists(ctx, sqlc.CheckSlugExistsParams{
		Slug:      req.Slug,
		ExcludeID: uuid.UUID{}, // Empty UUID for new pages
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

	// 1. Create page
	page, err := qtx.CreatePage(ctx, sqlc.CreatePageParams{
		Title:           req.Title,
		Slug:            req.Slug,
		PageType:        sqlc.PageType(req.PageType),
		Status:          statusToNullPageStatus(req.Status),
		FeaturedImageID: uuidToPgUUID(req.FeaturedImageID),
		AuthorID:        userID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to create page", err)
	}

	// 2. Create blocks using blocks service
	if len(req.Blocks) > 0 {
		if err := s.blockService.CreateBlocks(ctx, qtx, page.ID, req.Blocks); err != nil {
			return nil, err
		}
	}

	// 3. Create SEO if provided
	if req.SEO != nil {
		if err := s.createSEO(ctx, qtx, page.ID, req.SEO); err != nil {
			return nil, err
		}
	}

	// 4. Create project data if project page
	if req.PageType == "project" && req.ProjectData != nil {
		if err := s.createProjectData(ctx, qtx, page.ID, req.ProjectData); err != nil {
			return nil, err
		}
	}

	// 5. Create blog data if blog page
	if req.PageType == "blog" && req.BlogData != nil {
		if err := s.createBlogData(ctx, qtx, page.ID, req.BlogData); err != nil {
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
func (s *Service) ListPages(ctx context.Context, limit, offset int32) (*PageListResponse, error) {
	// Get total count
	totalCount, err := s.queries.CountPages(ctx, sqlc.CountPagesParams{
		Status:   sqlc.NullPageStatus{Valid: false}, // NULL for all statuses
		PageType: sqlc.NullPageType{Valid: false},   // NULL for all types
		AuthorID: pgtype.UUID{Valid: false},
	})
	if err != nil {
		return nil, errors.Internal("Failed to count pages", err)
	}

	// Get pages
	pages, err := s.queries.ListPages(ctx, sqlc.ListPagesParams{
		Status:    sqlc.NullPageStatus{Valid: false},
		PageType:  sqlc.NullPageType{Valid: false},
		AuthorID:  pgtype.UUID{Valid: false},
		SortBy:    "created_at_desc",
		OffsetVal: offset,
		LimitVal:  limit,
	})
	if err != nil {
		return nil, errors.Internal("Failed to list pages", err)
	}

	// Build responses
	pageResponses := make([]PageResponse, 0, len(pages))
	for _, page := range pages {
		resp, err := s.buildPageResponse(ctx, page)
		if err != nil {
			return nil, err
		}
		pageResponses = append(pageResponses, *resp)
	}

	// Calculate pagination
	totalPages := int(totalCount) / int(limit)
	if int(totalCount)%int(limit) > 0 {
		totalPages++
	}

	return &PageListResponse{
		Pages: pageResponses,
		Pagination: Pagination{
			Page:       int(offset/limit) + 1,
			Limit:      int(limit),
			TotalPages: totalPages,
			TotalCount: int(totalCount),
		},
	}, nil
}

// UpdatePage updates a page
func (s *Service) UpdatePage(ctx context.Context, pageID uuid.UUID, req UpdatePageRequest) (*PageResponse, error) {
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

	// Update page
	_, err := s.queries.UpdatePage(ctx, sqlc.UpdatePageParams{
		ID:              pageID,
		Title:           pointerToString(req.Title),
		Slug:            pointerToString(req.Slug),
		FeaturedImageID: uuidToPgUUID(req.FeaturedImageID),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Page not found", "page_not_found")
		}
		return nil, errors.Internal("Failed to update page", err)
	}

	return s.GetPageByID(ctx, pageID)
}

// UpdatePageStatus updates the page status
func (s *Service) UpdatePageStatus(ctx context.Context, pageID uuid.UUID, status string) error {
	_, err := s.queries.UpdatePageStatus(ctx, sqlc.UpdatePageStatusParams{
		ID:     pageID,
		Status: statusToNullPageStatus(status),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("Page not found", "page_not_found")
		}
		return errors.Internal("Failed to update page status", err)
	}
	return nil
}

// DeletePage soft deletes a page
func (s *Service) DeletePage(ctx context.Context, pageID uuid.UUID) error {
	err := s.queries.SoftDeletePage(ctx, pageID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("Page not found", "page_not_found")
		}
		return errors.Internal("Failed to delete page", err)
	}
	return nil
}

// buildPageResponse builds a complete PageResponse from a page
func (s *Service) buildPageResponse(ctx context.Context, page sqlc.Pages) (*PageResponse, error) {
	resp := &PageResponse{
		ID:              page.ID,
		Title:           page.Title,
		Slug:            page.Slug,
		PageType:        string(page.PageType),
		Status:          string(page.Status.PageStatus),
		FeaturedImageID: nullUUIDToPointer(page.FeaturedImageID),
		AuthorID:        page.AuthorID,
		CreatedAt:       page.CreatedAt,
		UpdatedAt:       page.UpdatedAt,
	}

	// Get blocks - note: blocks.Service doesn't have GetBlocksByPageID, use GetPageBlocks instead
	blockResponses, err := s.blockService.GetPageBlocks(ctx, page.ID)
	if err != nil {
		return nil, err
	}
	resp.Blocks = blockResponses

	// Get SEO
	seo, err := s.queries.GetPageSEO(ctx, page.ID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, errors.Internal("Failed to get page SEO", err)
	}
	if err == nil {
		resp.SEO = &SEOResponse{
			MetaTitle:       pgTextToPointer(seo.MetaTitle),
			MetaDescription: pgTextToPointer(seo.MetaDescription),
			OGTitle:         pgTextToPointer(seo.OgTitle),
			OGDescription:   pgTextToPointer(seo.OgDescription),
			OGImageID:       nullUUIDToPointer(seo.OgImageID),
			CanonicalURL:    pgTextToPointer(seo.CanonicalUrl),
			RobotsIndex:     seo.RobotsIndex.Bool,
			RobotsFollow:    seo.RobotsFollow.Bool,
		}
	}

	// Get project data if project page
	if page.PageType == "project" {
		projectData, err := s.queries.GetProjectData(ctx, page.ID)
		if err != nil && err != pgx.ErrNoRows {
			return nil, errors.Internal("Failed to get project data", err)
		}
		if err == nil {
			resp.ProjectData = &ProjectDataResponse{
				ClientName:    pgTextToPointer(projectData.ClientName),
				ProjectYear:   pgInt4ToPointer(projectData.ProjectYear),
				ProjectURL:    pgTextToPointer(projectData.ProjectUrl),
				Technologies:  bytesToStringSlice(projectData.Technologies),
				ProjectStatus: pgTextToPointer(projectData.ProjectStatus),
			}
		}
	}

	// Get blog data if blog page
	if page.PageType == "blog" {
		blogData, err := s.queries.GetBlogData(ctx, page.ID)
		if err != nil && err != pgx.ErrNoRows {
			return nil, errors.Internal("Failed to get blog data", err)
		}
		if err == nil {
			resp.BlogData = &BlogDataResponse{
				Excerpt:     pgTextToPointer(blogData.Excerpt),
				ReadingTime: pgInt4ToPointer(blogData.ReadingTime),
			}
		}
	}

	return resp, nil
}

// Helper functions for creating metadata

func (s *Service) createSEO(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *SEORequest) error {
	_, err := qtx.UpsertPageSEO(ctx, sqlc.UpsertPageSEOParams{
		PageID:          pageID,
		MetaTitle:       stringToPgText(req.MetaTitle),
		MetaDescription: stringToPgText(req.MetaDescription),
		OgTitle:         stringToPgText(req.OGTitle),
		OgDescription:   stringToPgText(req.OGDescription),
		OgImageID:       uuidToNullUUID(req.OGImageID),
		CanonicalUrl:    stringToPgText(req.CanonicalURL),
		RobotsIndex:     boolToPgBool(req.RobotsIndex),
		RobotsFollow:    boolToPgBool(req.RobotsFollow),
	})
	if err != nil {
		return errors.Internal("Failed to create SEO data", err)
	}
	return nil
}

func (s *Service) createProjectData(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *ProjectDataRequest) error {
	_, err := qtx.UpsertProjectData(ctx, sqlc.UpsertProjectDataParams{
		PageID:        pageID,
		ClientName:    stringToPgText(req.ClientName),
		ProjectYear:   intToPgInt4(req.ProjectYear),
		ProjectUrl:    stringToPgText(req.ProjectURL),
		Technologies:  stringSliceToBytes(req.Technologies),
		ProjectStatus: stringToPgText(req.ProjectStatus),
	})
	if err != nil {
		return errors.Internal("Failed to create project data", err)
	}
	return nil
}

func (s *Service) createBlogData(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *BlogDataRequest) error {
	_, err := qtx.UpsertBlogData(ctx, sqlc.UpsertBlogDataParams{
		PageID:      pageID,
		Excerpt:     stringToPgText(req.Excerpt),
		ReadingTime: intToPgInt4(req.ReadingTime),
		IsFeatured:  pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		return errors.Internal("Failed to create blog data", err)
	}
	return nil
}

// Helper functions for type conversions

func uuidToNullUUID(val *uuid.UUID) pgtype.UUID {
	if val == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *val, Valid: true}
}

func uuidToPgUUID(val *uuid.UUID) pgtype.UUID {
	if val == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *val, Valid: true}
}

func stringToPgText(val *string) pgtype.Text {
	if val == nil || *val == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *val, Valid: true}
}

func pointerToString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func intToPgInt4(val *int) pgtype.Int4 {
	if val == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*val), Valid: true}
}

func boolToPgBool(val *bool) pgtype.Bool {
	if val == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *val, Valid: true}
}

func timeToPgTimestamp(val *time.Time) pgtype.Timestamp {
	if val == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *val, Valid: true}
}

func statusToNullPageStatus(status string) sqlc.NullPageStatus {
	return sqlc.NullPageStatus{
		PageStatus: sqlc.PageStatus(status),
		Valid:      true,
	}
}

func nullUUIDToPointer(val pgtype.UUID) *uuid.UUID {
	if !val.Valid {
		return nil
	}
	id := uuid.UUID(val.Bytes)
	return &id
}

func pgTextToPointer(val pgtype.Text) *string {
	if !val.Valid {
		return nil
	}
	return &val.String
}

func pgInt4ToPointer(val pgtype.Int4) *int {
	if !val.Valid {
		return nil
	}
	intVal := int(val.Int32)
	return &intVal
}

func nullTimeToPointer(val pgtype.Timestamp) *time.Time {
	if !val.Valid {
		return nil
	}
	return &val.Time
}

// Helper to convert []string to []byte (JSON array)
func stringSliceToBytes(val []string) []byte {
	if val == nil {
		return nil
	}
	// Use json.Marshal for proper JSON encoding
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return []byte("[]")
	}
	return jsonBytes
}

// Helper to convert []byte (JSON array) to []string
func bytesToStringSlice(val []byte) []string {
	if val == nil || len(val) == 0 {
		return []string{}
	}

	var result []string
	if err := json.Unmarshal(val, &result); err != nil {
		return []string{}
	}
	return result
}
