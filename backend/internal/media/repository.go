package media

import (
	"context"

	"github.com/iankencruz/threefive/backend/internal/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Create(ctx context.Context, arg generated.CreateMediaParams) (*generated.Media, error)
	GetByID(ctx context.Context, id pgtype.UUID) (*generated.Media, error)
	List(ctx context.Context) ([]*generated.Media, error)
	UpdateSortOrder(ctx context.Context, id pgtype.UUID, sort int32) error
	Delete(ctx context.Context, id pgtype.UUID) error
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

func (r *MediaRepository) GetByID(ctx context.Context, id pgtype.UUID) (*generated.Media, error) {
	media, err := r.q.GetMediaByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) List(ctx context.Context) ([]*generated.Media, error) {
	medias, err := r.q.ListMedia(ctx)
	if err != nil {
		return nil, err
	}
	// convert []Media to []*Media
	out := make([]*generated.Media, len(medias))
	for i := range medias {
		out[i] = &medias[i]
	}
	return out, nil
}

func (r *MediaRepository) UpdateSortOrder(ctx context.Context, id pgtype.UUID, sort int32) error {
	return r.q.UpdateMediaSortOrder(ctx, generated.UpdateMediaSortOrderParams{
		ID:        id,
		SortOrder: sort,
	})
}

func (r *MediaRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteMedia(ctx, id)
}
