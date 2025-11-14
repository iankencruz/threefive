// backend/internal/blogs/service.go
package blogs

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service handles blogs business logic
type Service struct {
	db           *pgxpool.Pool
	queries      *sqlc.Queries
	blockService *blocks.Service
	config       ServiceConfig
}

type ServiceConfig struct {
	AutoPurgeRetentionDays int
}

// NewService creates a new blogs service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries, blockService *blocks.Service, cfg ServiceConfig) *Service {
	return &Service{
		db:           db,
		queries:      queries,
		blockService: blockService,
		config:       cfg,
	}
}

// CreateBlog creates a new blog with blocks in a transaction
func (s *Service) CreateBlog(ctx context.Context, req CreateBlogRequest) (*BlogResponse, error) {
	// Check slug uniqueness
	var nilUUID uuid.UUID

	exists, err := s.queries.CheckBlogSlugExists(ctx, sqlc.CheckBlogSlugExistsParams{
		Slug:      req.Slug,
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

	// 1. Create blog
	blog, err := qtx.CreateBlog(ctx, sqlc.CreateBlogParams{
		Title:           req.Title,
		Slug:            req.Slug,
		Status:          sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(req.Status), Valid: true},
		Excerpt:         stringToPgText(req.Excerpt),
		ReadingTime:     intToPgInt4(req.ReadingTime),
		IsFeatured:      pgtype.Bool{Bool: req.IsFeatured, Valid: true},
		FeaturedImageID: uuidToPgUUID(req.FeaturedImageID),
		PublishedAt:     pgtype.Timestamptz{Valid: false}, // Will be set when status changes to published
	})

	if err != nil {
		return nil, errors.Internal("Failed to create blog", err)
	}

	// 2. Create blocks using blocks service
	if len(req.Blocks) > 0 {
		if err := s.blockService.CreateBlocks(ctx, qtx, "blog", blog.ID, req.Blocks); err != nil {
			return nil, err
		}
	}

	// 3. Create SEO if provided
	if req.SEO != nil {
		if err := s.createSEO(ctx, qtx, blog.ID, req.SEO); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Internal("Failed to commit transaction", err)
	}

	// Get full blog with all relations
	return s.GetBlogByID(ctx, blog.ID)
}

// GetBlogByID retrieves a blog by ID with all relations
func (s *Service) GetBlogByID(ctx context.Context, blogID uuid.UUID) (*BlogResponse, error) {
	// Get blog
	blog, err := s.queries.GetBlogByID(ctx, blogID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Blog not found", "blog_not_found")
		}
		return nil, errors.Internal("Failed to get blog", err)
	}

	return s.buildBlogResponse(ctx, blog)
}

// GetBlogBySlug retrieves a blog by slug with all relations
func (s *Service) GetBlogBySlug(ctx context.Context, slug string) (*BlogResponse, error) {
	// Get blog
	blog, err := s.queries.GetBlogBySlug(ctx, slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Blog not found", "blog_not_found")
		}
		return nil, errors.Internal("Failed to get blog", err)
	}

	return s.buildBlogResponse(ctx, blog)
}

// ListBlogs retrieves blogs with pagination
func (s *Service) ListBlogs(ctx context.Context, limit, offset int32) (*BlogListResponse, error) {
	// Get total count
	totalCount, err := s.queries.CountBlogs(ctx, sqlc.CountBlogsParams{
		Status:     sqlc.PageStatus(""),
		IsFeatured: false, // plain bool for CountBlogs
	})

	if err != nil {
		return nil, errors.Internal("Failed to count blogs", err)
	}

	// Get blogs
	blogs, err := s.queries.ListBlogs(ctx, sqlc.ListBlogsParams{
		Status:     sqlc.PageStatus(""),
		IsFeatured: false, // plain bool for ListBlogs
		SortBy:     "created_at_desc",
		OffsetVal:  offset,
		LimitVal:   limit,
	})
	if err != nil {
		return nil, errors.Internal("Failed to list blogs", err)
	}

	// Build responses
	blogResponses := make([]BlogResponse, 0, len(blogs))
	for _, blog := range blogs {
		resp, err := s.buildBlogResponse(ctx, blog)
		if err != nil {
			return nil, err
		}
		blogResponses = append(blogResponses, *resp)
	}

	// Calculate pagination
	totalPages := int(totalCount) / int(limit)
	if int(totalCount)%int(limit) > 0 {
		totalPages++
	}
	currentPage := int(offset)/int(limit) + 1

	return &BlogListResponse{
		Blogs: blogResponses,
		Pagination: Pagination{
			Page:       currentPage,
			Limit:      int(limit),
			TotalPages: totalPages,
			TotalCount: int(totalCount),
		},
	}, nil
}

// UpdateBlog updates a blog and its related data in a transaction
func (s *Service) UpdateBlog(ctx context.Context, blogID uuid.UUID, req UpdateBlogRequest) (*BlogResponse, error) {

	// Check slug uniqueness if slug is being updated
	if req.Slug != nil {
		exists, err := s.queries.CheckBlogSlugExists(ctx, sqlc.CheckBlogSlugExistsParams{
			Slug:      *req.Slug,
			ExcludeID: blogID,
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

	// Update blog
	_, err = qtx.UpdateBlog(ctx, sqlc.UpdateBlogParams{
		ID:              blogID,
		Title:           pointerToString(req.Title),
		Slug:            pointerToString(req.Slug),
		Status:          statusToNullPageStatus(req.Status),
		Excerpt:         stringToPgText(req.Excerpt),
		ReadingTime:     intToPgInt4(req.ReadingTime),
		IsFeatured:      boolToPgBool(req.IsFeatured),
		FeaturedImageID: uuidToPgUUID(req.FeaturedImageID),
		PublishedAt:     pgtype.Timestamptz{Valid: false}, // Handled by status change
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Blog not found", "blog_not_found")
		}
		return nil, errors.Internal("Failed to update blog", err)
	}

	// Update blocks if provided
	if req.Blocks != nil {
		// Delete existing blocks
		if err := qtx.DeleteBlocksByEntity(ctx, sqlc.DeleteBlocksByEntityParams{
			EntityType: "blog",
			EntityID:   blogID,
		}); err != nil {
			return nil, errors.Internal("Failed to delete existing blocks", err)
		}

		// Create new blocks
		if len(*req.Blocks) > 0 {
			if err := s.blockService.CreateBlocks(ctx, qtx, "blog", blogID, *req.Blocks); err != nil {
				return nil, err
			}
		}
	}

	// Update SEO if provided
	if req.SEO != nil {
		if err := s.upsertSEO(ctx, qtx, blogID, req.SEO); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Internal("Failed to commit transaction", err)
	}

	// Get updated blog
	return s.GetBlogByID(ctx, blogID)
}

// UpdateBlogStatus updates only the blog status
func (s *Service) UpdateBlogStatus(ctx context.Context, blogID uuid.UUID, status string) (*BlogResponse, error) {
	_, err := s.queries.UpdateBlogStatus(ctx, sqlc.UpdateBlogStatusParams{
		ID:     blogID,
		Status: sqlc.NullPageStatus{PageStatus: sqlc.PageStatus(status), Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("Blog not found", "blog_not_found")
		}
		return nil, errors.Internal("Failed to update blog status", err)
	}

	return s.GetBlogByID(ctx, blogID)
}

// SoftDeleteBlog soft deletes a blog
func (s *Service) SoftDeleteBlog(ctx context.Context, blogID uuid.UUID) error {
	err := s.queries.SoftDeleteBlog(ctx, blogID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("Blog not found", "blog_not_found")
		}
		return errors.Internal("Failed to delete blog", err)
	}
	return nil
}

// PurgeOldDeletedBlogs permanently deletes blogs that were soft-deleted before the cutoff date
func (s *Service) PurgeOldDeletedBlogs(ctx context.Context) error {
	cutoffDate := time.Now().AddDate(0, 0, -s.config.AutoPurgeRetentionDays)

	// Note: Would need to implement PurgeOldDeletedBlogs query in SQLC similar to pages
	// For now, just log
	log.Printf("Would purge blogs deleted before %v (not implemented yet)", cutoffDate)

	return nil
}

// ============================================
// Private Helper Methods
// ============================================

// buildBlogResponse builds a complete blog response with all relations
func (s *Service) buildBlogResponse(ctx context.Context, blog sqlc.Blogs) (*BlogResponse, error) {
	resp := &BlogResponse{
		ID:         blog.ID,
		Title:      blog.Title,
		Slug:       blog.Slug,
		Status:     string(blog.Status.PageStatus),
		IsFeatured: blog.IsFeatured.Bool,
		CreatedAt:  blog.CreatedAt,
		UpdatedAt:  blog.UpdatedAt,
	}

	// Excerpt
	if blog.Excerpt.Valid {
		resp.Excerpt = &blog.Excerpt.String
	}

	// Reading time
	if blog.ReadingTime.Valid {
		readingTime := int(blog.ReadingTime.Int32)
		resp.ReadingTime = &readingTime
	}

	// Featured image
	if blog.FeaturedImageID.Valid {
		featuredImageID := uuid.UUID(blog.FeaturedImageID.Bytes)
		resp.FeaturedImageID = &featuredImageID
	}

	// Published at
	if blog.PublishedAt.Valid {
		resp.PublishedAt = &blog.PublishedAt.Time
	}

	// Deleted at
	if blog.DeletedAt.Valid {
		resp.DeletedAt = &blog.DeletedAt.Time
	}

	// Get blocks
	blockResponses, err := s.blockService.GetBlocksByEntity(ctx, "blog", blog.ID)
	if err != nil {
		return nil, err
	}
	resp.Blocks = blockResponses

	// Get SEO
	seo, err := s.queries.GetSEO(ctx, sqlc.GetSEOParams{
		EntityType: "blog",
		EntityID:   blog.ID,
	})
	if err != nil && err != pgx.ErrNoRows {
		return nil, errors.Internal("Failed to get blog SEO", err)
	}
	if err == nil {
		resp.SEO = buildSEOResponse(seo)
	}

	return resp, nil
}

// createSEO creates SEO data for a blog
func (s *Service) createSEO(ctx context.Context, qtx *sqlc.Queries, blogID uuid.UUID, req *SEORequest) error {
	_, err := qtx.CreateSEO(ctx, sqlc.CreateSEOParams{
		EntityType:      "blog",
		EntityID:        blogID,
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
func (s *Service) upsertSEO(ctx context.Context, qtx *sqlc.Queries, blogID uuid.UUID, req *SEORequest) error {
	_, err := qtx.UpsertSEO(ctx, sqlc.UpsertSEOParams{
		EntityType:      "blog",
		EntityID:        blogID,
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

// buildSEOResponse builds an SEO response
func buildSEOResponse(seo sqlc.Seo) *SEOResponse {
	resp := &SEOResponse{
		ID:           seo.ID,
		RobotsIndex:  seo.RobotsIndex.Bool,
		RobotsFollow: seo.RobotsFollow.Bool,
		CreatedAt:    seo.CreatedAt,
		UpdatedAt:    seo.UpdatedAt,
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
		return pgtype.Bool{Bool: false, Valid: true} // Default to false for is_featured in updates
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
