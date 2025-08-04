package blogs

import (
	"context"
	"fmt"

	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/generated"
)

type BlogService struct {
	repo         Repository
	BlockRepo    *blocks.Repository
	BlockService *blocks.Service
}

func NewBlogService(repo Repository, blockRepo *blocks.Repository, blockServices blocks.Service) *BlogService {
	return &BlogService{
		repo:         repo,
		BlockRepo:    blockRepo,
		BlockService: &blockServices,
	}
}

func (s *BlogService) GetBySlug(ctx context.Context, slug string) (BlogWithBlocks, error) {
	blog, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return BlogWithBlocks{}, err
	}

	rawBlocks, err := s.repo.GetBlocksForBlog(ctx, blog.ID)
	if err != nil {
		return BlogWithBlocks{}, err
	}

	var result []*blocks.BlockWithProps
	for _, b := range rawBlocks {
		var props any

		switch b.Type {
		case "heading":
			h, err := s.BlockRepo.GetHeadingBlock(ctx, b.ID)
			if err != nil {
				return BlogWithBlocks{}, fmt.Errorf("get heading block failed: %w", err)
			}
			props = map[string]any{
				"title":       h.Title,
				"description": h.Description,
			}

		case "richtext":
			rt, err := s.BlockRepo.GetRichTextBlockByID(ctx, b.ID)
			if err != nil {
				return BlogWithBlocks{}, fmt.Errorf("get richtext block failed: %w", err)
			}
			props = map[string]any{
				"html": rt.Html,
			}

		case "image":
			img, err := s.BlockRepo.GetImageBlock(ctx, b.ID)
			if err != nil {
				return BlogWithBlocks{}, fmt.Errorf("get image block failed: %w", err)
			}
			props = map[string]any{
				"media_id":   img.MediaID,
				"alt_text":   img.AltText,
				"align":      img.Align,
				"object_fit": img.ObjectFit,
			}

		default:
			props = map[string]any{}
		}

		result = append(result, &blocks.BlockWithProps{
			ID:        &b.ID,
			Type:      b.Type,
			SortOrder: b.SortOrder,
			Props:     props,
		})
	}

	return BlogWithBlocks{
		Blog:   blog,
		Blocks: result,
	}, nil
}

func (s *BlogService) List(ctx context.Context) ([]generated.Blog, error) {
	return s.repo.List(ctx)
}
