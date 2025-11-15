// backend/internal/shared/seo/seo.go
package seo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/utils"
	"github.com/iankencruz/threefive/internal/shared/validation"
	"github.com/jackc/pgx/v5"
)

// ============================================
// SEO Models
// ============================================

// Request represents SEO data in requests
type Request struct {
	MetaTitle       *string    `json:"meta_title,omitempty"`
	MetaDescription *string    `json:"meta_description,omitempty"`
	OGTitle         *string    `json:"og_title,omitempty"`
	OGDescription   *string    `json:"og_description,omitempty"`
	OGImageID       *uuid.UUID `json:"og_image_id,omitempty"`
	CanonicalURL    *string    `json:"canonical_url,omitempty"`
	RobotsIndex     *bool      `json:"robots_index,omitempty"`
	RobotsFollow    *bool      `json:"robots_follow,omitempty"`
}

// Response represents SEO data in responses
type Response struct {
	ID              uuid.UUID  `json:"id"`
	MetaTitle       *string    `json:"meta_title,omitempty"`
	MetaDescription *string    `json:"meta_description,omitempty"`
	OGTitle         *string    `json:"og_title,omitempty"`
	OGDescription   *string    `json:"og_description,omitempty"`
	OGImageID       *uuid.UUID `json:"og_image_id,omitempty"`
	CanonicalURL    *string    `json:"canonical_url,omitempty"`
	RobotsIndex     bool       `json:"robots_index"`
	RobotsFollow    bool       `json:"robots_follow"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ============================================
// SEO Service Functions
// ============================================

// Create creates SEO data for an entity
func Create(ctx context.Context, qtx *sqlc.Queries, entityType string, entityID uuid.UUID, req *Request) error {
	_, err := qtx.CreateSEO(ctx, sqlc.CreateSEOParams{
		EntityType:      entityType,
		EntityID:        entityID,
		MetaTitle:       utils.StrToPg(req.MetaTitle),
		MetaDescription: utils.StrToPg(req.MetaDescription),
		OgTitle:         utils.StrToPg(req.OGTitle),
		OgDescription:   utils.StrToPg(req.OGDescription),
		OgImageID:       utils.UUIDToPg(req.OGImageID),
		CanonicalUrl:    utils.StrToPg(req.CanonicalURL),
		RobotsIndex:     utils.BoolToPg(req.RobotsIndex, true),
		RobotsFollow:    utils.BoolToPg(req.RobotsFollow, true),
	})
	if err != nil {
		return errors.Internal("Failed to create SEO", err)
	}
	return nil
}

// Upsert updates or creates SEO data
func Upsert(ctx context.Context, qtx *sqlc.Queries, entityType string, entityID uuid.UUID, req *Request) error {
	_, err := qtx.UpsertSEO(ctx, sqlc.UpsertSEOParams{
		EntityType:      entityType,
		EntityID:        entityID,
		MetaTitle:       utils.StrToPg(req.MetaTitle),
		MetaDescription: utils.StrToPg(req.MetaDescription),
		OgTitle:         utils.StrToPg(req.OGTitle),
		OgDescription:   utils.StrToPg(req.OGDescription),
		OgImageID:       utils.UUIDToPg(req.OGImageID),
		CanonicalUrl:    utils.StrToPg(req.CanonicalURL),
		RobotsIndex:     utils.BoolToPg(req.RobotsIndex, true),
		RobotsFollow:    utils.BoolToPg(req.RobotsFollow, true),
	})
	if err != nil {
		return errors.Internal("Failed to upsert SEO", err)
	}
	return nil
}

// Get retrieves SEO data for an entity (returns nil if not found)
func Get(ctx context.Context, queries *sqlc.Queries, entityType string, entityID uuid.UUID) (*Response, error) {
	seo, err := queries.GetSEO(ctx, sqlc.GetSEOParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No SEO data is not an error
		}
		return nil, errors.Internal("Failed to get SEO", err)
	}

	return BuildResponse(seo), nil
}

// BuildResponse converts database SEO model to response
func BuildResponse(s sqlc.Seo) *Response {
	return &Response{
		ID:              s.ID,
		MetaTitle:       utils.PgToStr(s.MetaTitle),
		MetaDescription: utils.PgToStr(s.MetaDescription),
		OGTitle:         utils.PgToStr(s.OgTitle),
		OGDescription:   utils.PgToStr(s.OgDescription),
		OGImageID:       utils.PgToUUID(s.OgImageID),
		CanonicalURL:    utils.PgToStr(s.CanonicalUrl),
		RobotsIndex:     s.RobotsIndex.Bool,
		RobotsFollow:    s.RobotsFollow.Bool,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

// ============================================
// Validation
// ============================================

// Validate validates SEO request data
func Validate(v *validation.Validator, req *Request) {
	// Meta title
	if req.MetaTitle != nil {
		v.MaxLength("seo.meta_title", *req.MetaTitle, 60)
	}

	// Meta description
	if req.MetaDescription != nil {
		v.MaxLength("seo.meta_description", *req.MetaDescription, 160)
	}

	// OG title
	if req.OGTitle != nil {
		v.MaxLength("seo.og_title", *req.OGTitle, 60)
	}

	// OG description
	if req.OGDescription != nil {
		v.MaxLength("seo.og_description", *req.OGDescription, 160)
	}

	// Canonical URL
	if req.CanonicalURL != nil && *req.CanonicalURL != "" {
		v.URL("seo.canonical_url", *req.CanonicalURL)
	}
}
