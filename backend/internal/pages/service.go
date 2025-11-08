// backend/internal/pages/service.go
package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
func (s *Service) CreatePage(ctx context.Context, req CreatePageRequest, userID uuid.UUID) (*PageResponse, error) {
	// Check slug uniqueness - pass nil UUID for new pages
	var nilUUID uuid.UUID

	pageSlug := req.Slug

	switch req.PageType {
	case "project":
		pageSlug = fmt.Sprintf("projects/%s", req.Slug)
	case "blog":
		pageSlug = fmt.Sprintf("blog/%s", req.Slug)
	default:
		pageSlug = req.Slug
	}

	fmt.Printf("Checking project type: %s\n", req.PageType)
	fmt.Printf("Computed slug: %s\n", pageSlug)

	exists, err := s.queries.CheckSlugExists(ctx, sqlc.CheckSlugExistsParams{
		Slug:      pageSlug,
		ExcludeID: nilUUID,
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
		Slug:            pageSlug,
		PageType:        sqlc.PageType(req.PageType),
		Status:          sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(req.Status), Valid: true},
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
func (s *Service) ListPages(ctx context.Context, pageType string, limit, offset int32) (*PageListResponse, error) {
	var pageTypeParam sqlc.NullPageType
	if pageType != "" {
		// Only set the value and Valid flag if a pageType was provided (non-empty string)
		pageTypeParam = sqlc.NullPageType{PageType: sqlc.PageType(pageType), Valid: true}
	} else {
		// If pageType is "", it remains {Valid: false}, which translates to NULL
		// in SQL, activating the "OR ... IS NULL" logic to skip the filter.
		pageTypeParam = sqlc.NullPageType{Valid: false}
	}
	// Get total count
	totalCount, err := s.queries.CountPages(ctx, sqlc.CountPagesParams{
		Status:   sqlc.NullPageStatus{Valid: false}, // NULL for all statuses
		PageType: pageTypeParam,                     // NULL for all types
		AuthorID: pgtype.UUID{Valid: false},
	})
	if err != nil {
		return nil, errors.Internal("Failed to count pages", err)
	}

	// Get pages
	pages, err := s.queries.ListPages(ctx, sqlc.ListPagesParams{
		Status:    sqlc.NullPageStatus{Valid: false},
		PageType:  pageTypeParam,
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
	page := int(offset)/int(limit) + 1

	return &PageListResponse{
		Pages: pageResponses,
		Pagination: Pagination{
			Page:       page,
			Limit:      int(limit),
			TotalPages: totalPages,
			TotalCount: int(totalCount),
		},
	}, nil
}

// UpdatePage updates a page and its related data in a transaction
func (s *Service) UpdatePage(ctx context.Context, pageID uuid.UUID, req UpdatePageRequest) (*PageResponse, error) {
	// Get existing page to know its current state
	existingPage, err := s.queries.GetPageByID(ctx, pageID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Page not found", "page_not_found")
		}
		return nil, errors.Internal("Failed to get existing page", err)
	}

	// Determine the page_type to use for slug processing
	pageType := string(existingPage.PageType) // Use existing by default
	if req.PageType != nil {
		pageType = *req.PageType // Use new page_type if provided
	}

	// Process slug - add prefix based on page_type
	var finalSlug *string
	if req.Slug != nil {
		var computedSlug string

		switch pageType {
		case "project":
			computedSlug = fmt.Sprintf("projects/%s", *req.Slug)
		case "blog":
			computedSlug = fmt.Sprintf("blog/%s", *req.Slug)
		default:
			computedSlug = *req.Slug
		}

		finalSlug = &computedSlug

		// Check slug uniqueness with the computed slug
		exists, err := s.queries.CheckSlugExists(ctx, sqlc.CheckSlugExistsParams{
			Slug:      computedSlug,
			ExcludeID: pageID,
		})
		if err != nil {
			return nil, errors.Internal("Failed to check slug", err)
		}
		if exists {
			return nil, errors.Conflict("Slug already exists", "slug_exists")
		}
	} else if req.PageType != nil && *req.PageType != string(existingPage.PageType) {
		// If page_type is changing but slug is not provided, update slug with prefix
		var computedSlug string

		// Strip existing prefix from slug if it exists
		currentSlug := existingPage.Slug
		currentSlug = strings.TrimPrefix(currentSlug, "projects/")
		currentSlug = strings.TrimPrefix(currentSlug, "blog/")

		switch pageType {
		case "project":
			computedSlug = fmt.Sprintf("projects/%s", currentSlug)
		case "blog":
			computedSlug = fmt.Sprintf("blog/%s", currentSlug)
		default:
			computedSlug = currentSlug
		}

		finalSlug = &computedSlug

		// Check slug uniqueness
		exists, err := s.queries.CheckSlugExists(ctx, sqlc.CheckSlugExistsParams{
			Slug:      computedSlug,
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

	fmt.Printf("Page type changed, slug: %s\n", pointerToString(finalSlug))

	// Update page
	page, err := qtx.UpdatePage(ctx, sqlc.UpdatePageParams{
		ID:              pageID,
		Title:           pointerToString(req.Title),
		Slug:            pointerToString(finalSlug),
		PageType:        pageTypeToPageType(req.PageType),
		Status:          statusToNullPageStatus(req.Status),
		FeaturedImageID: uuidToPgUUID(req.FeaturedImageID),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Page not found", "page_not_found")
		}
		return nil, errors.Internal("Failed to update page", err)
	}

	// Update blocks if provided
	if req.Blocks != nil {
		// Delete existing blocks
		if err := qtx.DeleteBlocksByPageID(ctx, pageID); err != nil {
			return nil, errors.Internal("Failed to delete existing blocks", err)
		}

		// Create new blocks
		if len(*req.Blocks) > 0 {
			if err := s.blockService.CreateBlocks(ctx, qtx, pageID, *req.Blocks); err != nil {
				return nil, err
			}
		}
	}

	// Update SEO if provided
	if req.SEO != nil {
		if err := s.upsertSEO(ctx, qtx, pageID, req.SEO); err != nil {
			return nil, err
		}
	}

	// Update project data if provided
	if req.ProjectData != nil {
		if err := s.upsertProjectData(ctx, qtx, pageID, req.ProjectData); err != nil {
			return nil, err
		}
	}

	// Update blog data if provided
	if req.BlogData != nil {
		if err := s.upsertBlogData(ctx, qtx, pageID, req.BlogData); err != nil {
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

// HardDeletePage permanently deletes a page
func (s *Service) HardDeletePage(ctx context.Context, pageID uuid.UUID) error {
	err := s.queries.HardDeletePage(ctx, pageID)
	if err != nil {
		return errors.Internal("Failed to permanently delete page", err)
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
		PageType:  string(page.PageType),
		Status:    string(page.Status.PageStatus),
		AuthorID:  page.AuthorID,
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
		resp.SEO = buildSEOResponse(seo)
	}

	// Get project data if project page
	if page.PageType == "project" {
		projectData, err := s.queries.GetProjectData(ctx, page.ID)
		if err != nil && err != pgx.ErrNoRows {
			return nil, errors.Internal("Failed to get project data", err)
		}
		if err == nil {
			resp.ProjectData = buildProjectDataResponse(projectData)
		}
	}

	// Get blog data if blog page
	if page.PageType == "blog" {
		blogData, err := s.queries.GetBlogData(ctx, page.ID)
		if err != nil && err != pgx.ErrNoRows {
			return nil, errors.Internal("Failed to get blog data", err)
		}
		if err == nil {
			resp.BlogData = buildBlogDataResponse(blogData)
		}
	}

	return resp, nil
}

// createSEO creates SEO data for a page
func (s *Service) createSEO(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *SEORequest) error {
	_, err := qtx.CreatePageSEO(ctx, sqlc.CreatePageSEOParams{
		PageID:          pageID,
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
func (s *Service) upsertSEO(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *SEORequest) error {
	_, err := qtx.UpsertPageSEO(ctx, sqlc.UpsertPageSEOParams{
		PageID:          pageID,
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

// createProjectData creates project data for a page
func (s *Service) createProjectData(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *ProjectDataRequest) error {
	technologies, _ := json.Marshal(req.Technologies)

	_, err := qtx.CreateProjectData(ctx, sqlc.CreateProjectDataParams{
		PageID:        pageID,
		ClientName:    stringToPgText(req.ClientName),
		ProjectYear:   intToPgInt4(req.ProjectYear),
		ProjectUrl:    stringToPgText(req.ProjectURL),
		Technologies:  technologies,
		ProjectStatus: stringToPgText(req.ProjectStatus),
	})
	if err != nil {
		return errors.Internal("Failed to create project data", err)
	}
	return nil
}

// upsertProjectData updates or creates project data using the upsert query
func (s *Service) upsertProjectData(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *ProjectDataRequest) error {
	technologies, _ := json.Marshal(req.Technologies)

	_, err := qtx.UpsertProjectData(ctx, sqlc.UpsertProjectDataParams{
		PageID:        pageID,
		ClientName:    stringToPgText(req.ClientName),
		ProjectYear:   intToPgInt4(req.ProjectYear),
		ProjectUrl:    stringToPgText(req.ProjectURL),
		Technologies:  technologies,
		ProjectStatus: stringToPgText(req.ProjectStatus),
	})
	if err != nil {
		return errors.Internal("Failed to upsert project data", err)
	}
	return nil
}

// createBlogData creates blog data for a page
func (s *Service) createBlogData(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *BlogDataRequest) error {
	_, err := qtx.CreateBlogData(ctx, sqlc.CreateBlogDataParams{
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

// upsertBlogData updates or creates blog data using the upsert query
func (s *Service) upsertBlogData(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, req *BlogDataRequest) error {
	_, err := qtx.UpsertBlogData(ctx, sqlc.UpsertBlogDataParams{
		PageID:      pageID,
		Excerpt:     stringToPgText(req.Excerpt),
		ReadingTime: intToPgInt4(req.ReadingTime),
		IsFeatured:  pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		return errors.Internal("Failed to upsert blog data", err)
	}
	return nil
}

// buildSEOResponse builds an SEO response
func buildSEOResponse(seo sqlc.PageSeo) *SEOResponse {
	resp := &SEOResponse{
		RobotsIndex:  seo.RobotsIndex.Bool,
		RobotsFollow: seo.RobotsFollow.Bool,
	}

	if seo.MetaTitle.Valid {
		resp.MetaTitle = &seo.MetaTitle.String
	}
	if seo.MetaDescription.Valid {
		resp.MetaDescription = &seo.MetaDescription.String
	}
	if seo.OgTitle.Valid {
		resp.OGTitle = &seo.OgTitle.String
	}
	if seo.OgDescription.Valid {
		resp.OGDescription = &seo.OgDescription.String
	}
	if seo.OgImageID.Valid {
		ogImageID := uuid.UUID(seo.OgImageID.Bytes)
		resp.OGImageID = &ogImageID
	}
	if seo.CanonicalUrl.Valid {
		resp.CanonicalURL = &seo.CanonicalUrl.String
	}

	return resp
}

// buildProjectDataResponse builds a project data response
func buildProjectDataResponse(data sqlc.PageProjectData) *ProjectDataResponse {
	resp := &ProjectDataResponse{}

	if data.ClientName.Valid {
		resp.ClientName = &data.ClientName.String
	}
	if data.ProjectYear.Valid {
		year := int(data.ProjectYear.Int32)
		resp.ProjectYear = &year
	}
	if data.ProjectUrl.Valid {
		resp.ProjectURL = &data.ProjectUrl.String
	}
	if data.ProjectStatus.Valid {
		resp.ProjectStatus = &data.ProjectStatus.String
	}

	// Parse technologies JSON
	var technologies []string
	if err := json.Unmarshal(data.Technologies, &technologies); err == nil {
		resp.Technologies = technologies
	}

	return resp
}

// buildBlogDataResponse builds a blog data response
func buildBlogDataResponse(data sqlc.PageBlogData) *BlogDataResponse {
	resp := &BlogDataResponse{}

	if data.Excerpt.Valid {
		resp.Excerpt = &data.Excerpt.String
	}
	if data.ReadingTime.Valid {
		readingTime := int(data.ReadingTime.Int32)
		resp.ReadingTime = &readingTime
	}

	return resp
}

// ============================================
// Type conversion helpers
// ============================================

func uuidToPgUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

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

func boolToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Bool: true, Valid: true} // Default to true
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func statusToNullPageStatus(status *string) sqlc.NullPageStatus {
	if status == nil {
		return sqlc.NullPageStatus{Valid: false}
	}
	return sqlc.NullPageStatus{
		PageStatus: sqlc.PageStatus(*status),
		Valid:      true,
	}
}

func pointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func pageTypeToPageType(pageType *string) sqlc.PageType {
	if pageType == nil {
		return "" // Return empty string for null
	}
	return sqlc.PageType(*pageType)
}
