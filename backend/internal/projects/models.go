// backend/internal/projects/models.go
package projects

import (
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/seo"
)

// ============================================
// Request Models
// ============================================

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	Title           string                `json:"title"`
	Slug            string                `json:"slug"`
	Description     *string               `json:"description"`
	ProjectDate     *string               `json:"project_date"` // ISO 8601 date format
	Status          string                `json:"status"`
	ClientName      *string               `json:"client_name"`
	ProjectYear     *int                  `json:"project_year"`
	ProjectURL      *string               `json:"project_url"`
	Technologies    []string              `json:"technologies"`
	ProjectStatus   string                `json:"project_status"`
	FeaturedImageID *uuid.UUID            `json:"featured_image_id"`
	Blocks          []blocks.BlockRequest `json:"blocks"`
	SEO             *seo.Request          `json:"seo"`
}

// UpdateProjectRequest represents the request body for updating a project
type UpdateProjectRequest struct {
	Title           *string                `json:"title"`
	Slug            *string                `json:"slug"`
	Description     *string                `json:"description"`
	ProjectDate     *string                `json:"project_date"`
	Status          *string                `json:"status"`
	ClientName      *string                `json:"client_name"`
	ProjectYear     *int                   `json:"project_year"`
	ProjectURL      *string                `json:"project_url"`
	Technologies    *[]string              `json:"technologies"`
	ProjectStatus   *string                `json:"project_status"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id"`
	Blocks          *[]blocks.BlockRequest `json:"blocks"`
	SEO             *seo.Request           `json:"seo"`
}

// UpdateProjectStatusRequest represents the request body for updating project status
type UpdateProjectStatusRequest struct {
	Status string `json:"status"`
}

// ============================================
// Response Models
// ============================================

// ProjectResponse represents a project with all its data
type ProjectResponse struct {
	ID              uuid.UUID              `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	Description     *string                `json:"description"`
	ProjectDate     *string                `json:"project_date"`
	Status          string                 `json:"status"`
	ClientName      *string                `json:"client_name"`
	ProjectYear     *int                   `json:"project_year"`
	ProjectURL      *string                `json:"project_url"`
	Technologies    []string               `json:"technologies"`
	ProjectStatus   string                 `json:"project_status"`
	FeaturedImageID *uuid.UUID             `json:"featured_image_id"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	PublishedAt     *time.Time             `json:"published_at,omitempty"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty"`
	Blocks          []blocks.BlockResponse `json:"blocks"`
	SEO             *seo.Response          `json:"seo,omitempty"`
}

// ProjectListResponse represents a paginated list of projects
type ProjectListResponse struct {
	Projects   []ProjectResponse `json:"projects"`
	Pagination Pagination        `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}

// ListProjectsParams represents query parameters for listing/filtering projects
type ListProjectsParams struct {
	StatusFilter *string
	SortBy       string
	SortOrder    string
	Limit        int32
	Offset       int32
}
