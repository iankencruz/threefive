// backend/internal/blogs/service.go
package blogs

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
		Excerpt:         utils.StrToPg(req.Excerpt),
		ReadingTime:     utils.IntToPg(req.ReadingTime),
		IsFeatured:      pgtype.Bool{Bool: req.IsFeatured, Valid: true},
		FeaturedImageID: utils.UUIDToPg(req.FeaturedImageID),
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
		if err := seo.Create(ctx, qtx, "blog", blog.ID, req.SEO); err != nil {
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

// ListBlogs retrieves blogs with optional filtering and pagination
func (s *Service) ListBlogs(ctx context.Context, params ListBlogsParams) (*BlogListResponse, error) {
	// Convert filters to strings (empty string = no filter)
	statusStr := ""
	if params.StatusFilter != nil {
		statusStr = *params.StatusFilter
	}

	featuredStr := ""
	if params.FeaturedFilter != nil {
		if *params.FeaturedFilter {
			featuredStr = "true"
		} else {
			featuredStr = "false"
		}
	}

	// Get total count with filters
	totalCount, err := s.queries.CountBlogs(ctx, sqlc.CountBlogsParams{
		Status:     statusStr,
		IsFeatured: featuredStr,
	})
	if err != nil {
		return nil, errors.Internal("Failed to count blogs", err)
	}

	// Get blogs with filters
	blogs, err := s.queries.ListBlogs(ctx, sqlc.ListBlogsParams{
		Status:     statusStr,
		IsFeatured: featuredStr,
		SortBy:     params.SortBy,
		SortOrder:  params.SortOrder,
		OffsetVal:  params.Offset,
		LimitVal:   params.Limit,
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
	totalPages := int(totalCount) / int(params.Limit)
	if int(totalCount)%int(params.Limit) > 0 {
		totalPages++
	}
	currentPage := int(params.Offset)/int(params.Limit) + 1

	return &BlogListResponse{
		Blogs: blogResponses,
		Pagination: Pagination{
			Page:       currentPage,
			Limit:      int(params.Limit),
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
		Title:           utils.PtrStr(req.Title),
		Slug:            utils.PtrStr(req.Slug),
		Status:          utils.StatusToPg(req.Status),
		Excerpt:         utils.StrToPg(req.Excerpt),
		ReadingTime:     utils.IntToPg(req.ReadingTime),
		IsFeatured:      utils.BoolToPg(req.IsFeatured, true),
		FeaturedImageID: utils.UUIDToPg(req.FeaturedImageID),
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
		if err := seo.Upsert(ctx, qtx, "blog", blogID, req.SEO); err != nil {
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
	resp.SEO, err = seo.Get(ctx, s.queries, "blog", blog.ID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
