// backend/internal/pages/models.go
package pages

import (
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
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
	SEO             *SEORequest           `json:"seo,omitempty"`
}

// UpdatePageRequest represents the request to update a page
type UpdatePageRequest struct {
	Title           *string                `json:"title,omitempty"`
	Slug            *string                `json:"slug,omitempty"`
	Status          *string                `json:"status,omitempty"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id,omitempty"`
	Blocks          *[]blocks.BlockRequest `json:"blocks,omitempty"`
	SEO             *SEORequest            `json:"seo,omitempty"`
}

// UpdatePageStatusRequest represents the request to change page status
type UpdatePageStatusRequest struct {
	Status string `json:"status"` // draft, published, archived
}

// SEORequest represents SEO data in the request
type SEORequest struct {
	MetaTitle       *string    `json:"meta_title,omitempty"`
	MetaDescription *string    `json:"meta_description,omitempty"`
	OGTitle         *string    `json:"og_title,omitempty"`
	OGDescription   *string    `json:"og_description,omitempty"`
	OGImageID       *uuid.UUID `json:"og_image_id,omitempty"`
	CanonicalURL    *string    `json:"canonical_url,omitempty"`
	RobotsIndex     *bool      `json:"robots_index,omitempty"`
	RobotsFollow    *bool      `json:"robots_follow,omitempty"`
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
	AuthorID        uuid.UUID              `json:"author_id"`
	Blocks          []blocks.BlockResponse `json:"blocks,omitempty"`
	SEO             *SEOResponse           `json:"seo,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty"`
}

// SEOResponse represents SEO data in API responses
type SEOResponse struct {
	MetaTitle       *string    `json:"meta_title,omitempty"`
	MetaDescription *string    `json:"meta_description,omitempty"`
	OGTitle         *string    `json:"og_title,omitempty"`
	OGDescription   *string    `json:"og_description,omitempty"`
	OGImageID       *uuid.UUID `json:"og_image_id,omitempty"`
	CanonicalURL    *string    `json:"canonical_url,omitempty"`
	RobotsIndex     bool       `json:"robots_index"`
	RobotsFollow    bool       `json:"robots_follow"`
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
