// backend/internal/pages/models.go
package pages

import (
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/seo"
)

// ============================================
// Request Models
// ============================================

// CreatePageRequest represents the request to create a new page
type CreatePageRequest struct {
	Title           string                `json:"title"`
	Slug            string                `json:"slug"`
	Status          string                `json:"status"` // draft, published
	FeaturedImageID *uuid.UUID            `json:"featured_image_id,omitempty"`
	Blocks          []blocks.BlockRequest `json:"blocks"`
	SEO             *seo.Request          `json:"seo,omitempty"`
}

// UpdatePageRequest represents the request to update a page
type UpdatePageRequest struct {
	Title           *string                `json:"title,omitempty"`
	Slug            *string                `json:"slug,omitempty"`
	Status          *string                `json:"status,omitempty"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id,omitempty"`
	Blocks          *[]blocks.BlockRequest `json:"blocks,omitempty"`
	SEO             *seo.Request           `json:"seo,omitempty"`
}

// UpdatePageStatusRequest represents the request to change page status
type UpdatePageStatusRequest struct {
	Status string `json:"status"` // draft, published, archived
}

// ============================================
// Response Models
// ============================================

// PageResponse represents a page in API responses
type PageResponse struct {
	ID              uuid.UUID              `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	Status          string                 `json:"status"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id,omitempty"`
	Blocks          []blocks.BlockResponse `json:"blocks,omitempty"`
	SEO             *seo.Response          `json:"seo,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty"`
}

// PageListResponse represents a list of pages
type PageListResponse struct {
	Pages      []PageResponse `json:"pages"`
	Pagination Pagination     `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}

// ListPagesParams represents query parameters for listing/filtering pages
type ListPagesParams struct {
	StatusFilter *string
	SortBy       string
	SortOrder    string
	Limit        int32
	Offset       int32
}
