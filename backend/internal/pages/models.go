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
	PageType        string                `json:"page_type"` // generic, project, blog
	Status          string                `json:"status"`    // draft, published
	FeaturedImageID *uuid.UUID            `json:"featured_image_id,omitempty"`
	Blocks          []blocks.BlockRequest `json:"blocks"`
	SEO             *SEORequest           `json:"seo,omitempty"`
	ProjectData     *ProjectDataRequest   `json:"project_data,omitempty"`
	BlogData        *BlogDataRequest      `json:"blog_data,omitempty"`
}

// UpdatePageRequest represents the request to update a page
type UpdatePageRequest struct {
	Title           *string                `json:"title,omitempty"`
	Slug            *string                `json:"slug,omitempty"`
	Status          *string                `json:"status,omitempty"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id,omitempty"`
	Blocks          *[]blocks.BlockRequest `json:"blocks,omitempty"`
	SEO             *SEORequest            `json:"seo,omitempty"`
	ProjectData     *ProjectDataRequest    `json:"project_data,omitempty"`
	BlogData        *BlogDataRequest       `json:"blog_data,omitempty"`
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

// ProjectDataRequest represents project-specific data in the request
type ProjectDataRequest struct {
	ClientName    *string  `json:"client_name,omitempty"`
	ProjectYear   *int     `json:"project_year,omitempty"`
	ProjectURL    *string  `json:"project_url,omitempty"`
	Technologies  []string `json:"technologies,omitempty"`
	ProjectStatus *string  `json:"project_status,omitempty"` // completed, ongoing, archived
}

// BlogDataRequest represents blog-specific data in the request
type BlogDataRequest struct {
	Excerpt     *string `json:"excerpt,omitempty"`
	ReadingTime *int    `json:"reading_time,omitempty"`
}

// ============================================
// Response Models
// ============================================

// PageResponse represents a page in API responses
type PageResponse struct {
	ID              uuid.UUID              `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	PageType        string                 `json:"page_type"`
	Status          string                 `json:"status"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id,omitempty"`
	AuthorID        uuid.UUID              `json:"author_id"`
	Blocks          []blocks.BlockResponse `json:"blocks,omitempty"`
	SEO             *SEOResponse           `json:"seo,omitempty"`
	ProjectData     *ProjectDataResponse   `json:"project_data,omitempty"`
	BlogData        *BlogDataResponse      `json:"blog_data,omitempty"`
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

// ProjectDataResponse represents project data in API responses
type ProjectDataResponse struct {
	ClientName    *string  `json:"client_name,omitempty"`
	ProjectYear   *int     `json:"project_year,omitempty"`
	ProjectURL    *string  `json:"project_url,omitempty"`
	Technologies  []string `json:"technologies,omitempty"`
	ProjectStatus *string  `json:"project_status,omitempty"`
}

// BlogDataResponse represents blog data in API responses
type BlogDataResponse struct {
	Excerpt     *string    `json:"excerpt,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	ReadingTime *int       `json:"reading_time,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
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
