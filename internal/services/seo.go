package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SEOData is the view model passed to templates for rendering <meta> tags.
type SEOData struct {
	Title         string
	Description   string
	OGTitle       string
	OGDescription string
	OGImageURL    string
	CanonicalURL  string
	RobotsIndex   bool
	RobotsFollow  bool
	SiteName      string
	PageURL       string
}

// SEOResponse wraps the DB record with resolved fields ready for templates
type SEOResponse struct {
	Seo        *generated.Seo
	OGImageURL string
}

func (s *SEOResponse) HasSEO() bool {
	return s != nil && s.Seo != nil
}

func (s *SEOData) RobotsContent() string {
	index := "noindex"
	if s.RobotsIndex {
		index = "index"
	}
	follow := "nofollow"
	if s.RobotsFollow {
		follow = "follow"
	}
	return index + ", " + follow
}

func (s *SEOData) EffectiveOGTitle() string {
	if s.OGTitle != "" {
		return s.OGTitle
	}
	return s.Title
}

func (s *SEOData) EffectiveOGDescription() string {
	if s.OGDescription != "" {
		return s.OGDescription
	}
	return s.Description
}

// UpsertSEORequest is used by handlers to save SEO data for any entity.
type UpsertSEORequest struct {
	EntityType   string // "page", "project", "blog"
	EntityID     uuid.UUID
	SEOTitle     string
	SEODesc      string
	OGTitle      string
	OGDesc       string
	OGImageID    *uuid.UUID
	CanonicalURL string
	RobotsIndex  bool
	RobotsFollow bool
}

type SEOService struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	logger  *slog.Logger
}

func NewSEOService(db *pgxpool.Pool, queries *generated.Queries, logger *slog.Logger) *SEOService {
	return &SEOService{
		db:      db,
		queries: queries,
		logger:  logger.With("component", "seo_service"),
	}
}

// UpsertSEO creates or updates SEO data for any entity.
func (s *SEOService) UpsertSEO(ctx context.Context, req UpsertSEORequest) (*generated.Seo, error) {
	var ogImageID pgtype.UUID
	if req.OGImageID != nil {
		ogImageID = pgtype.UUID{Bytes: *(*[16]byte)(req.OGImageID[:]), Valid: true}
	}

	seo, err := s.queries.UpsertSEO(ctx, generated.UpsertSEOParams{
		EntityType:     req.EntityType,
		EntityID:       pgtype.UUID{Bytes: req.EntityID, Valid: true},
		SeoTitle:       pgtype.Text{String: req.SEOTitle, Valid: req.SEOTitle != ""},
		SeoDescription: pgtype.Text{String: req.SEODesc, Valid: req.SEODesc != ""},
		OgTitle:        pgtype.Text{String: req.OGTitle, Valid: req.OGTitle != ""},
		OgDescription:  pgtype.Text{String: req.OGDesc, Valid: req.OGDesc != ""},
		OgImageID:      ogImageID,
		CanonicalUrl:   pgtype.Text{String: req.CanonicalURL, Valid: req.CanonicalURL != ""},
		RobotsIndex:    req.RobotsIndex,
		RobotsFollow:   req.RobotsFollow,
	})
	if err != nil {
		s.logger.Error("failed to upsert SEO", "entity_type", req.EntityType, "entity_id", req.EntityID, "error", err)
		return nil, err
	}

	return &seo, nil
}

// GetSEO retrieves SEO data for any entity. Returns nil (not an error) if not found.
func (s *SEOService) GetSEO(ctx context.Context, entityType string, entityID uuid.UUID) (*generated.Seo, error) {
	seo, err := s.queries.GetSEO(ctx, generated.GetSEOParams{
		EntityType: entityType,
		EntityID:   pgtype.UUID{Bytes: entityID, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		s.logger.Error("failed to get SEO", "entity_type", entityType, "entity_id", entityID, "error", err)
		return nil, err
	}

	return &seo, nil
}

// DeleteSEO removes SEO data for an entity (call when entity is deleted).
func (s *SEOService) DeleteSEO(ctx context.Context, entityType string, entityID uuid.UUID) error {
	return s.queries.DeleteSEO(ctx, generated.DeleteSEOParams{
		EntityType: entityType,
		EntityID:   pgtype.UUID{Bytes: entityID, Valid: true},
	})
}

// ToSEOData converts a DB SEO record into the view model used in templates.
// ogImageURL should be pre-resolved by the caller from og_image_id.
// If seo is nil, sensible defaults are returned.
func ToSEOData(seo *generated.Seo, fallbackTitle, pageURL, siteName, ogImageURL string) *SEOData {
	if seo == nil {
		return &SEOData{
			Title:        fallbackTitle,
			SiteName:     siteName,
			PageURL:      pageURL,
			RobotsIndex:  true,
			RobotsFollow: true,
		}
	}

	title := fallbackTitle
	if seo.SeoTitle.Valid && seo.SeoTitle.String != "" {
		title = seo.SeoTitle.String
	}

	return &SEOData{
		Title:         title,
		Description:   seo.SeoDescription.String,
		OGTitle:       seo.OgTitle.String,
		OGDescription: seo.OgDescription.String,
		OGImageURL:    ogImageURL,
		CanonicalURL:  seo.CanonicalUrl.String,
		RobotsIndex:   seo.RobotsIndex,
		RobotsFollow:  seo.RobotsFollow,
		SiteName:      siteName,
		PageURL:       pageURL,
	}
}

func (s *SEOService) GetSEOResponse(ctx context.Context, entityType string, entityID uuid.UUID, mediaService *MediaService) (*SEOResponse, error) {
	seo, err := s.GetSEO(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}

	if seo == nil {
		return nil, nil // ← return nil, not &SEOResponse{Seo: nil}
	}

	resp := &SEOResponse{Seo: seo}

	if seo != nil && seo.OgImageID.Valid {
		if media, err := mediaService.GetMediaByID(ctx, seo.OgImageID); err == nil {
			mediaResp := mediaService.ToMediaResponse(media)
			resp.OGImageURL = mediaResp.ThumbnailURL
		}
	}

	return resp, nil
}
