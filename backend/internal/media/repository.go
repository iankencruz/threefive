package media

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/backend/internal/generated"
)

type Repository interface {
	Create(ctx context.Context, arg generated.CreateMediaParams) (*generated.Media, error)
	GetByID(ctx context.Context, id uuid.UUID) (*generated.Media, error)
	List(ctx context.Context) ([]generated.Media, error)
	ListPaginated(ctx context.Context, limit, offset int32) ([]generated.Media, error)
	UpdateSortOrder(ctx context.Context, id uuid.UUID, sort int32) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountMedia(ctx context.Context) (int, error)
	UpdateMedia(ctx context.Context, id uuid.UUID, title string, alt string) error
}

type MediaRepository struct {
	q *generated.Queries
}

func NewRepository(q *generated.Queries) Repository {
	return &MediaRepository{q: q}
}

func (r *MediaRepository) Create(ctx context.Context, arg generated.CreateMediaParams) (*generated.Media, error) {
	media, err := r.q.CreateMedia(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) GetByID(ctx context.Context, id uuid.UUID) (*generated.Media, error) {
	media, err := r.q.GetMediaByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) List(ctx context.Context) ([]generated.Media, error) {
	return r.q.ListMedia(ctx)
}

func (r *MediaRepository) UpdateSortOrder(ctx context.Context, id uuid.UUID, sort int32) error {
	return r.q.UpdateMediaSortOrder(ctx, generated.UpdateMediaSortOrderParams{
		ID:        id,
		SortOrder: sort,
	})
}

func (r *MediaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteMedia(ctx, id)
}

func (r *MediaRepository) ListPublic(ctx context.Context) ([]*generated.Media, error) {
	media, err := r.q.ListPublicMedia(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*generated.Media, len(media))
	for i := range media {
		out[i] = &media[i]
	}
	return out, nil
}
func (r *MediaRepository) ListPaginated(ctx context.Context, limit, offset int32) ([]generated.Media, error) {
	return r.q.ListMediaPaginated(ctx, generated.ListMediaPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *MediaRepository) CountMedia(ctx context.Context) (int, error) {
	count, err := r.q.CountMedia(ctx)
	return int(count), err
}

func (r *MediaRepository) UpdateMedia(ctx context.Context, id uuid.UUID, title, altText string) error {
	return r.q.UpdateMedia(ctx, generated.UpdateMediaParams{
		ID:      id,
		Title:   &title,
		AltText: &altText,
	})
}
