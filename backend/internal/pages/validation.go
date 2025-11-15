// backend/internal/pages/validation.go
package pages

import (
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/seo"
	"github.com/iankencruz/threefive/internal/shared/validation"
)

// Valid page types
var ValidEntityTypes = []string{"generic", "blog"}

// Valid page statuses
var ValidPageStatuses = []string{"draft", "published", "archived"}

// Valid project statuses
var ValidProjectStatuses = []string{"completed", "ongoing", "archived"}

// Validate validates a CreatePageRequest
func (r *CreatePageRequest) Validate(v *validation.Validator) {
	// Title validation
	v.Required("title", r.Title)
	v.MinLength("title", r.Title, 1)
	v.MaxLength("title", r.Title, 200)

	// Slug validation
	v.Required("slug", r.Slug)
	v.Slug("slug", r.Slug)
	v.MinLength("slug", r.Slug, 1)
	v.MaxLength("slug", r.Slug, 200)

	// Status validation
	v.Required("status", r.Status)
	v.In("status", r.Status, ValidPageStatuses)

	// Blocks validation
	if len(r.Blocks) > 0 {
		blocks.ValidateBlocks(v, r.Blocks)
	}

	// SEO validation
	if r.SEO != nil {
		seo.Validate(v, r.SEO)
	}
}

// Validate validates an UpdatePageRequest
func (r *UpdatePageRequest) Validate(v *validation.Validator) {
	// At least one field must be provided
	if r.Title == nil && r.Slug == nil && r.FeaturedImageID == nil {
		v.AddError("update", "At least one field must be provided for update")
		return
	}

	// Title validation (if provided)
	if r.Title != nil {
		v.MinLength("title", *r.Title, 1)
		v.MaxLength("title", *r.Title, 200)
	}

	// Slug validation (if provided)
	if r.Slug != nil {
		v.Slug("slug", *r.Slug)
		v.MinLength("slug", *r.Slug, 1)
		v.MaxLength("slug", *r.Slug, 200)
	}

	// Status validation (if provided)
	if r.Status != nil {
		v.In("status", *r.Status, ValidPageStatuses)
	}
}

// Validate validates an UpdatePageStatusRequest
func (r *UpdatePageStatusRequest) Validate(v *validation.Validator) {
	v.Required("status", r.Status)
	v.In("status", r.Status, ValidPageStatuses)
}
