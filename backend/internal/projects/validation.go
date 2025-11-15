// backend/internal/projects/validation.go
package projects

import (
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/seo"
	"github.com/iankencruz/threefive/internal/shared/validation"
)

// Valid project statuses
var ValidProjectStatuses = []string{"completed", "ongoing", "archived"}

// Valid page statuses (reused for projects)
var ValidPageStatuses = []string{"draft", "published", "archived"}

// Validate validates a CreateProjectRequest
func (r *CreateProjectRequest) Validate(v *validation.Validator) {
	// Required fields
	v.Required("title", r.Title)
	v.Required("slug", r.Slug)
	v.Required("status", r.Status)
	v.Required("project_status", r.ProjectStatus)

	// Title validation
	if r.Title != "" {
		v.MinLength("title", r.Title, 3)
		v.MaxLength("title", r.Title, 200)
	}

	// Slug validation
	if r.Slug != "" {
		v.Slug("slug", r.Slug)
		v.MinLength("slug", r.Slug, 3)
		v.MaxLength("slug", r.Slug, 200)
	}

	// Status validation
	v.In("status", r.Status, ValidPageStatuses)

	// Project status validation
	v.In("project_status", r.ProjectStatus, ValidProjectStatuses)

	// Project year validation
	if r.ProjectYear != nil {
		if *r.ProjectYear < 1900 || *r.ProjectYear > 2100 {
			v.AddError("project_year", "Project year must be between 1900 and 2100")
		}
	}

	// Project URL validation
	if r.ProjectURL != nil && *r.ProjectURL != "" {
		v.URL("project_url", *r.ProjectURL)
	}

	// Validate blocks if provided
	if len(r.Blocks) > 0 {
		blocks.ValidateBlocks(v, r.Blocks)
	}

	// Validate SEO if provided
	if r.SEO != nil {
		seo.Validate(v, r.SEO)
	}
}

// Validate validates an UpdateProjectRequest
func (r *UpdateProjectRequest) Validate(v *validation.Validator) {
	// Title validation (if provided)
	if r.Title != nil {
		if *r.Title == "" {
			v.AddError("title", "Title cannot be empty")
		} else {
			v.MinLength("title", *r.Title, 3)
			v.MaxLength("title", *r.Title, 200)
		}
	}

	// Slug validation (if provided)
	if r.Slug != nil {
		if *r.Slug == "" {
			v.AddError("slug", "Slug cannot be empty")
		} else {
			v.Slug("slug", *r.Slug)
			v.MinLength("slug", *r.Slug, 3)
			v.MaxLength("slug", *r.Slug, 200)
		}
	}

	// Status validation (if provided)
	if r.Status != nil {
		v.In("status", *r.Status, ValidPageStatuses)
	}

	// Project status validation (if provided)
	if r.ProjectStatus != nil {
		v.In("project_status", *r.ProjectStatus, ValidProjectStatuses)
	}

	// Project year validation (if provided)
	if r.ProjectYear != nil {
		if *r.ProjectYear < 1900 || *r.ProjectYear > 2100 {
			v.AddError("project_year", "Project year must be between 1900 and 2100")
		}
	}

	// Project URL validation (if provided)
	if r.ProjectURL != nil && *r.ProjectURL != "" {
		v.URL("project_url", *r.ProjectURL)
	}

	// Validate blocks if provided
	if r.Blocks != nil {
		blocks.ValidateBlocks(v, *r.Blocks)
	}

	// Validate SEO if provided
	if r.SEO != nil {
		seo.Validate(v, r.SEO)
	}
}

// Validate validates an UpdateProjectStatusRequest
func (r *UpdateProjectStatusRequest) Validate(v *validation.Validator) {
	v.Required("status", r.Status)
	v.In("status", r.Status, ValidPageStatuses)
}
