package gallery

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type GalleryService struct {
	repo Repository
}

func NewGalleryService(repo Repository) *GalleryService {
	return &GalleryService{repo: repo}
}

func (s *GalleryService) Create(ctx context.Context, arg generated.CreateGalleryParams) (*generated.Gallery, error) {
	if arg.PublishedAt.Valid {
		arg.PublishedAt = pgtype.Timestamptz{Valid: false}
	}
	return s.repo.CreateGallery(ctx, arg)
}

func (s *GalleryService) Update(ctx context.Context, arg generated.UpdateGalleryParams) (*generated.Gallery, error) {
	if arg.PublishedAt.Valid && !arg.IsPublished {
		arg.PublishedAt = pgtype.Timestamptz{Valid: false}
	}

	var ts pgtype.Timestamptz
	_ = ts.Scan(time.Now())
	arg.UpdatedAt = ts

	return s.repo.UpdateGallery(ctx, arg)
}

func (s *GalleryService) ListByPage(ctx context.Context, pageID uuid.UUID) ([]generated.Gallery, error) {
	return s.repo.GetByPage(ctx, pageID)
}

func (s *GalleryService) LinkToPage(ctx context.Context, galleryID, pageID uuid.UUID) error {
	return s.repo.LinkToPage(ctx, galleryID, pageID)
}

func (s *GalleryService) UnlinkFromPage(ctx context.Context, galleryID, pageID uuid.UUID) error {
	return s.repo.UnlinkFromPage(ctx, galleryID, pageID)
}

func (s *GalleryService) GetGallery(ctx context.Context, slug string) (*GalleryWithMedia, error) {
	return s.repo.GetGalleryBySlug(ctx, slug)
}
