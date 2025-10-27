package gallery

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewService(db *pgxpool.Pool, queries *sqlc.Queries) *Service {
	return &Service{
		db:      db,
		queries: queries,
	}
}

func (s *Service) CreateGallery(ctx context.Context, req CreateGalleryRequest) (*GalleryWithMedia, error) {
	// Create gallery
	gallery, err := s.queries.CreateGallery(ctx, sqlc.CreateGalleryParams{
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		return nil, err
	}

	// Link media if provided
	if len(req.MediaIDs) > 0 {
		for i, mediaID := range req.MediaIDs {
			_, err := s.queries.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
				MediaID:    mediaID,
				EntityType: "gallery",
				EntityID:   gallery.ID,
				SortOrder:  pgtype.Int4{Int32: int32(i), Valid: true},
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return s.GetGallery(ctx, gallery.ID)
}

func (s *Service) GetGallery(ctx context.Context, id uuid.UUID) (*GalleryWithMedia, error) {
	gallery, err := s.queries.GetGalleryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get media for gallery using existing media_relations
	media, err := s.queries.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
		EntityType: "gallery",
		EntityID:   id,
	})
	if err != nil {
		return nil, err
	}

	// Get media count
	count, err := s.queries.GetGalleryMediaCount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GalleryWithMedia{
		ID:          gallery.ID,
		Title:       gallery.Title,
		Description: gallery.Description.String,
		Media:       media,
		MediaCount:  count,
		CreatedAt:   gallery.CreatedAt,
		UpdatedAt:   gallery.UpdatedAt,
	}, nil
}

func (s *Service) ListGalleries(ctx context.Context, limit, offset int32) ([]GalleryListItem, error) {
	galleries, err := s.queries.ListGalleries(ctx, sqlc.ListGalleriesParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, err
	}

	var result []GalleryListItem
	for _, g := range galleries {
		count, err := s.queries.GetGalleryMediaCount(ctx, g.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, GalleryListItem{
			ID:          g.ID,
			Title:       g.Title,
			Description: g.Description.String,
			MediaCount:  count,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		})
	}

	return result, nil
}

func (s *Service) UpdateGallery(ctx context.Context, id uuid.UUID, req UpdateGalleryRequest) (*GalleryWithMedia, error) {
	// Get existing gallery to preserve values if not provided
	existing, err := s.queries.GetGalleryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Use existing values if new ones are empty
	title := req.Title
	if title == "" {
		title = existing.Title
	}

	description := pgtype.Text{String: req.Description, Valid: req.Description != ""}
	if req.Description == "" && existing.Description.Valid {
		description = existing.Description
	}

	// Update gallery
	_, err = s.queries.UpdateGallery(ctx, sqlc.UpdateGalleryParams{
		Title:       title,
		Description: description,
		ID:          id,
	})
	if err != nil {
		return nil, err
	}

	// If media IDs provided, update the links
	if req.MediaIDs != nil {
		// Get existing media relations
		existingMedia, err := s.queries.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
			EntityType: "gallery",
			EntityID:   id,
		})
		if err != nil {
			return nil, err
		}

		// Remove all existing links
		for _, media := range existingMedia {
			err := s.queries.UnlinkMediaFromEntity(ctx, sqlc.UnlinkMediaFromEntityParams{
				MediaID:    media.ID,
				EntityType: "gallery",
				EntityID:   id,
			})
			if err != nil {
				return nil, err
			}
		}

		// Add new links with updated sort order
		for i, mediaID := range req.MediaIDs {
			_, err := s.queries.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
				MediaID:    mediaID,
				EntityType: "gallery",
				EntityID:   id,
				SortOrder:  pgtype.Int4{Int32: int32(i), Valid: true},
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return s.GetGallery(ctx, id)
}

func (s *Service) DeleteGallery(ctx context.Context, id uuid.UUID) error {
	// Get existing media relations
	existingMedia, err := s.queries.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
		EntityType: "gallery",
		EntityID:   id,
	})
	if err != nil {
		return err
	}

	// Unlink all media
	for _, media := range existingMedia {
		err := s.queries.UnlinkMediaFromEntity(ctx, sqlc.UnlinkMediaFromEntityParams{
			MediaID:    media.ID,
			EntityType: "gallery",
			EntityID:   id,
		})
		if err != nil {
			return err
		}
	}

	return s.queries.DeleteGallery(ctx, id)
}

func (s *Service) LinkMedia(ctx context.Context, galleryID, mediaID uuid.UUID, sortOrder int32) error {
	_, err := s.queries.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
		MediaID:    mediaID,
		EntityType: "gallery",
		EntityID:   galleryID,
		SortOrder:  pgtype.Int4{Int32: sortOrder, Valid: true},
	})
	return err
}

func (s *Service) UnlinkMedia(ctx context.Context, galleryID, mediaID uuid.UUID) error {
	return s.queries.UnlinkMediaFromEntity(ctx, sqlc.UnlinkMediaFromEntityParams{
		MediaID:    mediaID,
		EntityType: "gallery",
		EntityID:   galleryID,
	})
}
