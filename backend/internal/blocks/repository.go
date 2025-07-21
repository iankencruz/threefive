package blocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/generated"
)

type Repository struct {
	q *generated.Queries
}

func NewRepository(q *generated.Queries) *Repository {
	return &Repository{q: q}
}

func (r *Repository) UpdateBlockSortOrder(ctx context.Context, arg generated.UpdateBlockSortOrderParams) error {
	return r.q.UpdateBlockSortOrder(ctx, arg)
}

func (r *Repository) GetBaseBlocksFromPage(ctx context.Context, pageID uuid.UUID) ([]generated.Block, error) {
	return r.q.GetBlocksForPage(ctx, pageID)
}

func (r *Repository) GetBlogBlocks(ctx context.Context, blogID uuid.UUID) ([]generated.Block, error) {
	return r.q.GetBlocksForBlog(ctx, blogID)
}

func (r *Repository) GetBlockByID(ctx context.Context, id uuid.UUID) (generated.Block, error) {
	return r.q.GetBlockByID(ctx, id)
}

func (r *Repository) CreateBlock(ctx context.Context, arg generated.CreateBlockParams) (generated.Block, error) {
	return r.q.CreateBlock(ctx, arg)
}

func (r *Repository) UpdateBlock(ctx context.Context, arg generated.UpdateBlockParams) error {
	return r.q.UpdateBlock(ctx, arg)
}

func (r *Repository) DeleteBlockByID(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteBlock(ctx, id)
}

// Image Blocks
func (r *Repository) CreateImageBlock(ctx context.Context, arg generated.CreateImageBlockParams) error {
	return r.q.CreateImageBlock(ctx, arg)
}

func (r *Repository) UpdateImageBlock(ctx context.Context, arg generated.UpdateImageBlockParams) error {
	return r.q.UpdateImageBlock(ctx, arg)
}

func (r *Repository) DeleteImageBlock(ctx context.Context, blockID uuid.UUID) error {
	return r.q.DeleteImageBlock(ctx, blockID)
}

func (r *Repository) GetImageBlock(ctx context.Context, blockID uuid.UUID) (*generated.ImageBlock, error) {
	block, err := r.q.GetImageBlock(ctx, blockID)
	if err != nil {
		return nil, err
	}
	return &block, err
}

func (r *Repository) CreateRichTextBlock(ctx context.Context, arg generated.CreateRichTextBlockParams) error {
	return r.q.CreateRichTextBlock(ctx, arg)
}

func (r *Repository) UpdateRichTextBlock(ctx context.Context, arg generated.UpdateRichTextBlockParams) error {
	return r.q.UpdateRichTextBlock(ctx, arg)
}

func (r *Repository) CreateHeadingBlock(ctx context.Context, arg generated.CreateHeadingBlockParams) error {
	return r.q.CreateHeadingBlock(ctx, arg)
}

func (r *Repository) UpdateHeadingBlock(ctx context.Context, arg generated.UpdateHeadingBlockParams) error {
	return r.q.UpdateHeadingBlock(ctx, arg)
}

func (r *Repository) GetHeadingBlock(ctx context.Context, id uuid.UUID) (*generated.HeadingBlock, error) {
	block, err := r.q.GetHeadingBlock(ctx, id)
	if err != nil {
		return nil, err
	}
	return &block, err
}

func (r *Repository) GetRichTextBlockByID(ctx context.Context, id uuid.UUID) (*generated.RichtextBlock, error) {
	block, err := r.q.GetRichTextBlock(ctx, id)
	if err != nil {
		return nil, err
	}
	return &block, err
}

func (r *Repository) DeleteRichTextBlockByID(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteRichTextBlock(ctx, id)
}

func (r *Repository) GetHeadingBlockByID(ctx context.Context, id uuid.UUID) (*generated.HeadingBlock, error) {
	block, err := r.q.GetHeadingBlock(ctx, id)
	if err != nil {
		return nil, err
	}
	return &block, err
}

func (r *Repository) DeleteHeadingBlockByID(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteHeadingBlock(ctx, id)
}
