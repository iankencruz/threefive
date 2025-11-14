// backend/internal/blogs/models.go
package blogs

import (
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
)

// ============================================
// Request Models
// ============================================

// CreateBlogRequest represents the request body for creating a blog
type CreateBlogRequest struct {
	Title           string                `json:"title"`
	Slug            string                `json:"slug"`
	Status          string                `json:"status"` // draft, published, archived
	Excerpt         *string               `json:"excerpt"`
	ReadingTime     *int                  `json:"reading_time"`
	IsFeatured      bool                  `json:"is_featured"`
	FeaturedImageID *uuid.UUID            `json:"featured_image_id"`
	Blocks          []blocks.BlockRequest `json:"blocks"`
	SEO             *SEORequest           `json:"seo"`
}

// UpdateBlogRequest represents the request body for updating a blog
type UpdateBlogRequest struct {
	Title           *string                `json:"title"`
	Slug            *string                `json:"slug"`
	Status          *string                `json:"status"`
	Excerpt         *string                `json:"excerpt"`
	ReadingTime     *int                   `json:"reading_time"`
	IsFeatured      *bool                  `json:"is_featured"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id"`
	Blocks          *[]blocks.BlockRequest `json:"blocks"`
	SEO             *SEORequest            `json:"seo"`
}

// UpdateBlogStatusRequest represents the request body for updating blog status
type UpdateBlogStatusRequest struct {
	Status string `json:"status"` // draft, published, archived
}

// SEORequest represents SEO metadata in requests
type SEORequest struct {
	MetaTitle       *string    `json:"meta_title"`
	MetaDescription *string    `json:"meta_description"`
	OGTitle         *string    `json:"og_title"`
	OGDescription   *string    `json:"og_description"`
	OGImageID       *uuid.UUID `json:"og_image_id"`
	CanonicalURL    *string    `json:"canonical_url"`
	RobotsIndex     *bool      `json:"robots_index"`
	RobotsFollow    *bool      `json:"robots_follow"`
}

// ============================================
// Response Models
// ============================================

// BlogResponse represents a blog with all its data
type BlogResponse struct {
	ID              uuid.UUID              `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	Status          string                 `json:"status"`
	Excerpt         *string                `json:"excerpt"`
	ReadingTime     *int                   `json:"reading_time"`
	IsFeatured      bool                   `json:"is_featured"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	PublishedAt     *time.Time             `json:"published_at,omitempty"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty"`
	Blocks          []blocks.BlockResponse `json:"blocks"`
	SEO             *SEOResponse           `json:"seo,omitempty"`
}

// SEOResponse represents SEO metadata in responses
type SEOResponse struct {
	ID              uuid.UUID  `json:"id"`
	MetaTitle       *string    `json:"meta_title"`
	MetaDescription *string    `json:"meta_description"`
	OGTitle         *string    `json:"og_title"`
	OGDescription   *string    `json:"og_description"`
	OGImageID       *uuid.UUID `json:"og_image_id"`
	CanonicalURL    *string    `json:"canonical_url"`
	RobotsIndex     bool       `json:"robots_index"`
	RobotsFollow    bool       `json:"robots_follow"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// BlogListResponse represents a paginated list of blogs
type BlogListResponse struct {
	Blogs      []BlogResponse `json:"blogs"`
	Pagination Pagination     `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}
