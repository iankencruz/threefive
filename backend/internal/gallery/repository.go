package gallery

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/generated"
)

type Repository interface {
	CreateGallery(ctx context.Context, arg generated.CreateGalleryParams) (*generated.Gallery, error)
	GetGalleryByID(ctx context.Context, id uuid.UUID) (*generated.Gallery, error)
	GetGalleryBySlug(ctx context.Context, slug string) (*generated.Gallery, error)
	ListGalleries(ctx context.Context) ([]generated.Gallery, error)
	DeleteGallery(ctx context.Context, id uuid.UUID) error
	UpdateGallery(ctx context.Context, arg generated.UpdateGalleryParams) (*generated.Gallery, error)
	// AddMediaToGallery(ctx context.Context, arg generated.AddMediaToGalleryParams) error
	// RemoveMediaFromGallery(ctx context.Context, arg generated.RemoveMediaFromGalleryParams) error
	// ListMediaFromGallery(ctx context.Context, galleryID uuid.UUID) ([]generated.Gallery, error)
	// UpdateGalleryMediaSortOrder(ctx context.Context, arg generated.UpdateGalleryMediaSortOrderParams) error
}

type GalleryRepository struct {
	q *generated.Queries
}

func NewRepository(q *generated.Queries) Repository {
	return &GalleryRepository{q: q}
}

func (r *GalleryRepository) CreateGallery(ctx context.Context, arg generated.CreateGalleryParams) (*generated.Gallery, error) {
	gallery, err := r.q.CreateGallery(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &gallery, err
}

func (r *GalleryRepository) GetGalleryByID(ctx context.Context, id uuid.UUID) (*generated.Gallery, error) {
	gallery, err := r.q.GetGalleryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *GalleryRepository) GetGalleryBySlug(ctx context.Context, slug string) (*generated.Gallery, error) {
	gallery, err := r.q.GetGalleryBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *GalleryRepository) ListGalleries(ctx context.Context) ([]generated.Gallery, error) {
	return r.q.ListGalleries(ctx)
}

func (r *GalleryRepository) DeleteGallery(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteGallery(ctx, id)
}

func (r *GalleryRepository) UpdateGallery(ctx context.Context, arg generated.UpdateGalleryParams) (*generated.Gallery, error) {
	gallery, err := r.q.UpdateGallery(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}
