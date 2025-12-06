// backend/internal/blogs/validation.go
package blogs

import (
	"strconv"

	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/shared/seo"
	"github.com/iankencruz/threefive/internal/shared/validation"
)

// Valid blog statuses (reuses page_status enum)
var ValidBlogStatuses = []string{"draft", "published", "archived"}

// Validate validates a CreateBlogRequest
func (r *CreateBlogRequest) Validate(v *validation.Validator) {
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
	v.In("status", r.Status, ValidBlogStatuses)

	// Description validation (optional)
	if r.Description != nil {
		v.MaxLength("description", *r.Description, 500)
	}

	// Reading time validation (optional, must be positive)
	if r.ReadingTime != nil {
		v.MinLength("reading_time", strconv.Itoa(*r.ReadingTime), 1)
		v.MaxLength("reading_time", strconv.Itoa(*r.ReadingTime), 999) // Max ~16 hours
	}

	// IsFeatured is a boolean, no validation needed (can be true/false)

	// Blocks validation
	if len(r.Blocks) > 0 {
		blocks.ValidateBlocks(v, r.Blocks)
	}

	// SEO validation
	if r.SEO != nil {
		seo.Validate(v, r.SEO)
	}
}

// Validate validates an UpdateBlogRequest
func (r *UpdateBlogRequest) Validate(v *validation.Validator) {
	// At least one field must be provided
	if r.Title == nil && r.Slug == nil && r.Status == nil &&
		r.Description == nil && r.ReadingTime == nil && r.IsFeatured == nil &&
		r.FeaturedImageID == nil && r.Blocks == nil && r.SEO == nil {
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
		v.In("status", *r.Status, ValidBlogStatuses)
	}

	// Description validation (if provided)
	if r.Description != nil {
		v.MaxLength("description", *r.Description, 500)
	}

	// Reading time validation (if provided)
	if r.ReadingTime != nil {
		v.MinLength("reading_time", strconv.Itoa(*r.ReadingTime), 1)
		v.MaxLength("reading_time", strconv.Itoa(*r.ReadingTime), 999)
	}

	// IsFeatured - no validation needed (boolean)

	// Blocks validation (if provided)
	if r.Blocks != nil && len(*r.Blocks) > 0 {
		blocks.ValidateBlocks(v, *r.Blocks)
	}

	// SEO validation (if provided)
	if r.SEO != nil {
		seo.Validate(v, r.SEO)
	}
}

// Validate validates an UpdateBlogStatusRequest
func (r *UpdateBlogStatusRequest) Validate(v *validation.Validator) {
	v.Required("status", r.Status)
	v.In("status", r.Status, ValidBlogStatuses)
}
