// internal/services/tags.go
package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type TagService struct {
	queries *generated.Queries
}

func NewTagService(queries *generated.Queries) *TagService {
	return &TagService{
		queries: queries,
	}
}

// CreateTagRequest represents the data needed to create a tag
type CreateTagRequest struct {
	Name string
	Slug string // Optional, will be auto-generated if empty
}

// UpdateTagRequest represents the data needed to update a tag
type UpdateTagRequest struct {
	Name *string
	Slug *string
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, req *CreateTagRequest) (*generated.Tag, error) {
	// Generate slug if not provided
	if req.Slug == "" {
		req.Slug = GenerateSlug(req.Name)
	}

	if !IsValidSlug(req.Slug) {
		return nil, fmt.Errorf("invalid slug format: must be lowercase, alphanumeric with hyphens only")
	}

	// Check slug uniqueness
	exists, err := s.queries.CheckTagSlugExists(ctx, generated.CheckTagSlugExistsParams{
		Slug:  req.Slug,
		TagID: pgtype.UUID{Valid: false},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("slug already exists: %s", req.Slug)
	}

	tagID := uuid.New()

	tag, err := s.queries.CreateTag(ctx, generated.CreateTagParams{
		ID: pgtype.UUID{
			Bytes: tagID,
			Valid: true,
		},
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return &tag, nil
}

// GetTagByID retrieves a tag by ID
func (s *TagService) GetTagByID(ctx context.Context, tagID uuid.UUID) (*generated.Tag, error) {
	tag, err := s.queries.GetTagByID(ctx, pgtype.UUID{
		Bytes: tagID,
		Valid: true,
	})
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	return &tag, nil
}

// GetTagBySlug retrieves a tag by slug
func (s *TagService) GetTagBySlug(ctx context.Context, slug string) (*generated.Tag, error) {
	tag, err := s.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	return &tag, nil
}

// ListTags retrieves a paginated list of tags
func (s *TagService) ListTags(ctx context.Context, limit, offset int32) ([]generated.Tag, error) {
	tags, err := s.queries.ListTags(ctx, generated.ListTagsParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// ListAllTags retrieves all tags without pagination
func (s *TagService) ListAllTags(ctx context.Context) ([]generated.Tag, error) {
	tags, err := s.queries.ListAllTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all tags: %w", err)
	}

	return tags, nil
}

// UpdateTag updates a tag
func (s *TagService) UpdateTag(ctx context.Context, tagID uuid.UUID, req *UpdateTagRequest) (*generated.Tag, error) {
	// Validate new slug if provided
	if req.Slug != nil {
		if !IsValidSlug(*req.Slug) {
			return nil, fmt.Errorf("invalid slug format")
		}

		exists, err := s.queries.CheckTagSlugExists(ctx, generated.CheckTagSlugExistsParams{
			Slug: *req.Slug,
			TagID: pgtype.UUID{
				Bytes: tagID,
				Valid: true,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("slug already exists: %s", *req.Slug)
		}
	}

	// Build update params
	updateParams := generated.UpdateTagParams{
		ID: pgtype.UUID{
			Bytes: tagID,
			Valid: true,
		},
	}

	if req.Name != nil {
		updateParams.Name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if req.Slug != nil {
		updateParams.Slug = pgtype.Text{String: *req.Slug, Valid: true}
	}

	tag, err := s.queries.UpdateTag(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return &tag, nil
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	// Check if tag is in use
	usageCount, err := s.queries.GetTagUsageCount(ctx, pgtype.UUID{
		Bytes: tagID,
		Valid: true,
	})
	if err != nil {
		return fmt.Errorf("failed to check tag usage: %w", err)
	}

	if usageCount > 0 {
		return fmt.Errorf("cannot delete tag: it is used by %d projects", usageCount)
	}

	if err := s.queries.DeleteTag(ctx, pgtype.UUID{
		Bytes: tagID,
		Valid: true,
	}); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// CountTags returns the total count of tags
func (s *TagService) CountTags(ctx context.Context) (int64, error) {
	return s.queries.CountTags(ctx)
}

// SearchTags searches for tags by name
func (s *TagService) SearchTags(ctx context.Context, searchTerm string, limit, offset int32) ([]generated.Tag, error) {
	tags, err := s.queries.SearchTags(ctx, generated.SearchTagsParams{
		SearchTerm: "%" + searchTerm + "%",
		LimitVal:   limit,
		OffsetVal:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	return tags, nil
}

// GetMostUsedTags retrieves the most used tags
func (s *TagService) GetMostUsedTags(ctx context.Context, limit int32) ([]TagResponse, error) {
	rows, err := s.queries.GetMostUsedTags(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most used tags: %w", err)
	}

	responses := make([]TagResponse, len(rows))
	for i, row := range rows {
		// Reconstruct the Tag struct from individual fields
		responses[i] = TagResponse{
			Tag: generated.Tag{
				ID:        row.ID,
				Name:      row.Name,
				Slug:      row.Slug,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			UsageCount: row.UsageCount,
		}
	}

	return responses, nil
}

// GetUnusedTags retrieves tags that aren't used by any projects
func (s *TagService) GetUnusedTags(ctx context.Context) ([]generated.Tag, error) {
	tags, err := s.queries.GetUnusedTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unused tags: %w", err)
	}

	return tags, nil
}
