package pages

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/generated"
)

type Repository interface {
	CreatePage(ctx context.Context, arg generated.CreatePageParams) (*generated.Page, error)
	GetPageByID(ctx context.Context, id uuid.UUID) (*generated.Page, error)
	GetPageBySlug(ctx context.Context, slug string) (*generated.Page, error)
	ListPages(ctx context.Context, sort string) ([]generated.Page, error)
	UpdatePage(ctx context.Context, arg generated.UpdatePageParams) (*generated.Page, error)
	DeletePage(ctx context.Context, id uuid.UUID) error
}

type PageRepository struct {
	q *generated.Queries
}

func NewRepository(q *generated.Queries) Repository {
	return &PageRepository{q: q}
}

func (r *PageRepository) CreatePage(ctx context.Context, arg generated.CreatePageParams) (*generated.Page, error) {
	page, err := r.q.CreatePage(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) GetPageByID(ctx context.Context, id uuid.UUID) (*generated.Page, error) {
	page, err := r.q.GetPageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) GetPageBySlug(ctx context.Context, slug string) (*generated.Page, error) {
	page, err := r.q.GetPageBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) ListPages(ctx context.Context, sort string) ([]generated.Page, error) {
	switch sort {
	case "desc":
		return r.q.ListPagesByUpdatedAsc(ctx)
	case "asc":
		return r.q.ListPagesByTitleAsc(ctx)
	default:
		return r.q.ListPagesByUpdatedDesc(ctx) // fallback
	}
}

func (r *PageRepository) UpdatePage(ctx context.Context, arg generated.UpdatePageParams) (*generated.Page, error) {
	page, err := r.q.UpdatePage(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepository) DeletePage(ctx context.Context, id uuid.UUID) error {
	return r.q.DeletePage(ctx, id)
}
