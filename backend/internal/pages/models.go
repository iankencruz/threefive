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
	Title           *string    `json:"title,omitempty"`
	Slug            *string    `json:"slug,omitempty"`
	FeaturedImageID *uuid.UUID `json:"featured_image_id,omitempty"`
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
	IsFeatured  *bool   `json:"is_featured,omitempty"`
}

// ============================================
// Response Models
// ============================================

// PageResponse represents a page in the response
type PageResponse struct {
	ID              uuid.UUID              `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	PageType        string                 `json:"page_type"`
	Status          string                 `json:"status"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id,omitempty"`
	AuthorID        uuid.UUID              `json:"author_id"`
	Blocks          []blocks.BlockResponse `json:"blocks"`
	SEO             *SEOResponse           `json:"seo,omitempty"`
	ProjectData     *ProjectDataResponse   `json:"project_data,omitempty"`
	BlogData        *BlogDataResponse      `json:"blog_data,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	PublishedAt     *time.Time             `json:"published_at,omitempty"`
}

// PageListResponse represents a paginated list of pages
type PageListResponse struct {
	Pages      []PageSummary `json:"pages"`
	Pagination Pagination    `json:"pagination"`
}

// PageSummary represents a page in list view (without blocks)
type PageSummary struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	PageType        string     `json:"page_type"`
	Status          string     `json:"status"`
	FeaturedImageID *uuid.UUID `json:"featured_image_id,omitempty"`
	AuthorID        uuid.UUID  `json:"author_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
}

// SEOResponse represents SEO data in the response
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

// ProjectDataResponse represents project data in the response
type ProjectDataResponse struct {
	ClientName    *string  `json:"client_name,omitempty"`
	ProjectYear   *int     `json:"project_year,omitempty"`
	ProjectURL    *string  `json:"project_url,omitempty"`
	Technologies  []string `json:"technologies"`
	ProjectStatus string   `json:"project_status"`
}

// BlogDataResponse represents blog data in the response
type BlogDataResponse struct {
	Excerpt     *string `json:"excerpt,omitempty"`
	ReadingTime *int    `json:"reading_time,omitempty"`
	IsFeatured  bool    `json:"is_featured"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ListPagesFilters represents filters for listing pages
type ListPagesFilters struct {
	Status   *string
	PageType *string
	AuthorID *uuid.UUID
	SortBy   string // created_at_desc, updated_at_desc, title_asc, etc
	Page     int
	Limit    int
}
