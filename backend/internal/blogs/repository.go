package blogs

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/generated"
)

type Repository interface {
	GetBySlug(ctx context.Context, slug string) (generated.Blog, error)
	List(ctx context.Context) ([]generated.Blog, error)
	Create(ctx context.Context, arg generated.CreateBlogParams) (generated.Blog, error)
	Update(ctx context.Context, arg generated.UpdateBlogParams) (generated.Blog, error)
	Delete(ctx context.Context, slug string) error
	GetBlocksForBlog(ctx context.Context, blogID uuid.UUID) ([]generated.Block, error)
}

type PgxBlogRepository struct {
	queries   *generated.Queries
	blockRepo *blocks.Repository
}

func NewRepository(q *generated.Queries, blockRepo *blocks.Repository) *PgxBlogRepository {
	return &PgxBlogRepository{
		queries:   q,
		blockRepo: blockRepo,
	}
}

func (r *PgxBlogRepository) GetBySlug(ctx context.Context, slug string) (generated.Blog, error) {
	return r.queries.GetBlogBySlug(ctx, slug)
}

func (r *PgxBlogRepository) List(ctx context.Context) ([]generated.Blog, error) {
	return r.queries.ListBlogs(ctx)
}

func (r *PgxBlogRepository) Create(ctx context.Context, arg generated.CreateBlogParams) (generated.Blog, error) {
	return r.queries.CreateBlog(ctx, arg)
}

func (r *PgxBlogRepository) Update(ctx context.Context, arg generated.UpdateBlogParams) (generated.Blog, error) {
	return r.queries.UpdateBlog(ctx, arg)
}

func (r *PgxBlogRepository) Delete(ctx context.Context, slug string) error {
	return r.queries.DeleteBlog(ctx, slug)
}

func (r *PgxBlogRepository) GetBlocksForBlog(ctx context.Context, blogID uuid.UUID) ([]blocks.BlockWithProps, error) {
	generatedBlocks, err := r.blockRepo.GetBlogBlocks(ctx, blogID)
	if err != nil {
		return nil, err
	}

	var result []blocks.BlockWithProps
	for _, b := range generatedBlocks {
		result = append(result, blocks.BlockWithProps{
			ID:        &b.ID,
			Type:      b.Type,
			SortOrder: b.SortOrder,
			Props:     b.Props,
		})
	}

	return result, nil
}
