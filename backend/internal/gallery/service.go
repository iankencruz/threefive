package gallery

import (
	"context"
	"time"

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
