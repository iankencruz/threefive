// backend/internal/pages/validation.go
package pages

import (
	"github.com/iankencruz/threefive/internal/blocks"
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
		r.SEO.Validate(v)
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

// Validate validates SEORequest
func (r *SEORequest) Validate(v *validation.Validator) {
	// Meta title
	if r.MetaTitle != nil {
		v.MaxLength("seo.meta_title", *r.MetaTitle, 60)
	}

	// Meta description
	if r.MetaDescription != nil {
		v.MaxLength("seo.meta_description", *r.MetaDescription, 160)
	}

	// OG title
	if r.OGTitle != nil {
		v.MaxLength("seo.og_title", *r.OGTitle, 60)
	}

	// OG description
	if r.OGDescription != nil {
		v.MaxLength("seo.og_description", *r.OGDescription, 160)
	}

	// Canonical URL
	if r.CanonicalURL != nil && *r.CanonicalURL != "" {
		v.URL("seo.canonical_url", *r.CanonicalURL)
	}
}
