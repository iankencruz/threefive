package blogs

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/generated"
)

type Repository interface {
	GetBySlug(ctx context.Context, slug string) (generated.Blog, error)
	List(ctx context.Context) ([]generated.Blog, error)
	Create(ctx context.Context, arg generated.CreateBlogParams) (generated.Blog, error)
	Update(ctx context.Context, arg generated.UpdateBlogParams) (generated.Blog, error)
	Delete(ctx context.Context, slug string) error

	ListMediaForGallery(ctx context.Context, galleryID uuid.UUID) ([]generated.Media, error)
}

type BlogRepository struct {
	queries *generated.Queries
}

func NewRepository(q *generated.Queries) *BlogRepository {
	return &BlogRepository{
		queries: q,
	}
}

func (r *BlogRepository) GetBySlug(ctx context.Context, slug string) (generated.Blog, error) {
	return r.queries.GetBlogBySlug(ctx, slug)
}

func (r *BlogRepository) List(ctx context.Context) ([]generated.Blog, error) {
	return r.queries.ListBlogs(ctx)
}

func (r *BlogRepository) Create(ctx context.Context, arg generated.CreateBlogParams) (generated.Blog, error) {
	return r.queries.CreateBlog(ctx, arg)
}

func (r *BlogRepository) Update(ctx context.Context, arg generated.UpdateBlogParams) (generated.Blog, error) {
	return r.queries.UpdateBlog(ctx, arg)
}

func (r *BlogRepository) Delete(ctx context.Context, slug string) error {
	return r.queries.DeleteBlog(ctx, slug)
}

func (r *BlogRepository) ListMediaForGallery(ctx context.Context, galleryID uuid.UUID) ([]generated.Media, error) {
	return r.queries.ListMediaForGallery(ctx, galleryID)
}
