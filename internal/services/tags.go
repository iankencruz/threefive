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
	Name string
	Slug string
}

// Validate validates the create request
func (r *CreateTagRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Name == "" {
		errors["name"] = "Tag name is required"
	}

	if r.Slug != "" && !IsValidSlug(r.Slug) {
		errors["slug"] = "Slug must contain only lowercase letters, numbers, and hyphens"
	}

	return errors
}

// Validate validates the update request
func (r *UpdateTagRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Name == "" {
		errors["name"] = "Tag name is required"
	}

	if r.Slug == "" {
		errors["slug"] = "Slug is required"
	} else if !IsValidSlug(r.Slug) {
		errors["slug"] = "Slug must contain only lowercase letters, numbers, and hyphens"
	}

	return errors
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

// GetTagBySlugWithUsage retrieves a tag by slug with usage count
func (s *TagService) GetTagBySlugWithUsage(ctx context.Context, slug string) (*TagResponse, error) {
	tag, err := s.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	// Get usage count
	usageCount, err := s.queries.GetTagUsageCount(ctx, tag.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage count: %w", err)
	}

	return &TagResponse{
		Tag:        tag,
		UsageCount: usageCount,
	}, nil
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

// ListTagsWithUsage retrieves a paginated list of tags with usage counts
func (s *TagService) ListTagsWithUsage(ctx context.Context, limit, offset int32) ([]TagResponse, error) {
	tags, err := s.queries.ListTags(ctx, generated.ListTagsParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	results := make([]TagResponse, len(tags))
	for i, tag := range tags {
		usageCount, err := s.queries.GetTagUsageCount(ctx, tag.ID)
		if err != nil {
			usageCount = 0
		}
		results[i] = TagResponse{
			Tag:        tag,
			UsageCount: usageCount,
		}
	}

	return results, nil
}

// ListAllTags retrieves all tags without pagination
func (s *TagService) ListAllTags(ctx context.Context) ([]generated.Tag, error) {
	tags, err := s.queries.ListAllTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all tags: %w", err)
	}

	return tags, nil
}

// UpdateTag updates a tag by ID (for internal use)
func (s *TagService) UpdateTag(ctx context.Context, tagID uuid.UUID, req *UpdateTagRequest) (*generated.Tag, error) {
	// Validate new slug
	if !IsValidSlug(req.Slug) {
		return nil, fmt.Errorf("invalid slug format")
	}

	exists, err := s.queries.CheckTagSlugExists(ctx, generated.CheckTagSlugExistsParams{
		Slug: req.Slug,
		TagID: pgtype.UUID{
			Bytes: tagID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("slug already exists: %s", req.Slug)
	}

	// Update tag
	tag, err := s.queries.UpdateTag(ctx, generated.UpdateTagParams{
		ID: pgtype.UUID{
			Bytes: tagID,
			Valid: true,
		},
		Name: pgtype.Text{String: req.Name, Valid: true},
		Slug: pgtype.Text{String: req.Slug, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return &tag, nil
}

// UpdateTagBySlug updates a tag by slug (used by handlers)
func (s *TagService) UpdateTagBySlug(ctx context.Context, slug string, req *UpdateTagRequest) (*generated.Tag, error) {
	// Get existing tag to find ID
	existingTag, err := s.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	// Validate new slug
	if !IsValidSlug(req.Slug) {
		return nil, fmt.Errorf("invalid slug format")
	}

	// Check if new slug conflicts with another tag
	if req.Slug != slug {
		exists, err := s.queries.CheckTagSlugExists(ctx, generated.CheckTagSlugExistsParams{
			Slug:  req.Slug,
			TagID: existingTag.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("slug already exists: %s", req.Slug)
		}
	}

	// Update tag
	tag, err := s.queries.UpdateTag(ctx, generated.UpdateTagParams{
		ID:   existingTag.ID,
		Name: pgtype.Text{String: req.Name, Valid: true},
		Slug: pgtype.Text{String: req.Slug, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return &tag, nil
}

// DeleteTag deletes a tag by ID
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

// DeleteTagBySlug deletes a tag by slug (used by handlers)
func (s *TagService) DeleteTagBySlug(ctx context.Context, slug string) error {
	// Get tag ID from slug
	tag, err := s.queries.GetTagBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	// Check if tag is in use
	usageCount, err := s.queries.GetTagUsageCount(ctx, tag.ID)
	if err != nil {
		return fmt.Errorf("failed to check tag usage: %w", err)
	}

	if usageCount > 0 {
		return fmt.Errorf("cannot delete tag '%s': it is used by %d project(s)", tag.Name, usageCount)
	}

	// Delete the tag
	if err := s.queries.DeleteTag(ctx, tag.ID); err != nil {
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

// SearchTagsWithUsage searches for tags by name with usage counts
func (s *TagService) SearchTagsWithUsage(ctx context.Context, searchTerm string, limit, offset int32) ([]TagResponse, error) {
	tags, err := s.queries.SearchTags(ctx, generated.SearchTagsParams{
		SearchTerm: "%" + searchTerm + "%",
		LimitVal:   limit,
		OffsetVal:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	results := make([]TagResponse, len(tags))
	for i, tag := range tags {
		usageCount, err := s.queries.GetTagUsageCount(ctx, tag.ID)
		if err != nil {
			usageCount = 0
		}
		results[i] = TagResponse{
			Tag:        tag,
			UsageCount: usageCount,
		}
	}

	return results, nil
}

// GetMostUsedTags retrieves the most used tags
func (s *TagService) GetMostUsedTags(ctx context.Context, limit int32) ([]TagResponse, error) {
	rows, err := s.queries.GetMostUsedTags(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most used tags: %w", err)
	}

	responses := make([]TagResponse, len(rows))
	for i, row := range rows {
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

// DeleteUnusedTags bulk deletes all unused tags
func (s *TagService) DeleteUnusedTags(ctx context.Context) (int, error) {
	unusedTags, err := s.queries.GetUnusedTags(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get unused tags: %w", err)
	}

	deleted := 0
	for _, tag := range unusedTags {
		if err := s.queries.DeleteTag(ctx, tag.ID); err != nil {
			continue
		}
		deleted++
	}

	return deleted, nil
}
