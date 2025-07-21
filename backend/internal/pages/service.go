package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/generated"
)

type PageService struct {
	Repo         Repository
	BlockRepo    *blocks.Repository
	BlockService *blocks.Service
}

func NewPageService(repo Repository, blockRepo *blocks.Repository, blockService blocks.Service) *PageService {
	return &PageService{
		Repo:         repo,
		BlockRepo:    blockRepo,
		BlockService: &blockService,
	}
}

// Create inserts a new page and auto-generates the slug from title
func (s *PageService) Create(ctx context.Context, arg generated.CreatePageParams) (*generated.Page, error) {
	arg.Slug = slug.Make(strings.ToLower(arg.Title))
	return s.Repo.CreatePage(ctx, arg)
}

// Update modifies a page and returns the fully updated page record
func (s *PageService) Update(ctx context.Context, arg generated.UpdatePageParams) (*generated.Page, error) {
	updatedPage, err := s.Repo.UpdatePage(ctx, arg)
	if err != nil {
		return nil, err
	}
	return updatedPage, nil
}

// GetPageWithBlocks fetches a page and its associated blocks by slug

func (s *PageService) GetPageBlocks(ctx context.Context, pageID uuid.UUID) ([]*blocks.BlockWithProps, error) {
	baseBlocks, err := s.BlockRepo.GetBaseBlocksFromPage(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("get base blocks failed: %w", err)
	}

	var result []*blocks.BlockWithProps
	for _, b := range baseBlocks {
		var props any

		switch b.Type {
		case "heading":
			h, err := s.BlockService.Repo.GetHeadingBlock(ctx, b.ID)
			if err != nil {
				return nil, fmt.Errorf("get heading block failed: %w", err)
			}
			props = map[string]any{
				"title":       h.Title,
				"description": h.Description,
			}

		case "richtext":
			rt, err := s.BlockService.Repo.GetRichTextBlockByID(ctx, b.ID)
			if err != nil {
				return nil, fmt.Errorf("get richtext block failed: %w", err)
			}
			props = map[string]any{
				"html": rt.Html,
			}

		case "image":
			img, err := s.BlockService.Repo.GetImageBlock(ctx, b.ID)
			if err != nil {
				return nil, fmt.Errorf("get image block failed: %w", err)
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
			Props:     props, // THIS is what makes title/description show up
		})
	}

	fmt.Printf("📦 Final BlocksWithProps: %+v\n", result)
	return result, nil
}

func (s *PageService) UpdateWithBlocks(ctx context.Context, req UpdatePageWithBlocksRequest) (*generated.Page, error) {
	page, err := s.Repo.UpdatePage(ctx, req.Page)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🔍 Update blocks received: %d\n", len(req.Blocks))
	for i, b := range req.Blocks {
		fmt.Printf("Block %d: %+v\n", i, b)
	}

	err = s.BlockService.SyncBlocks(ctx, "page", req.Page.ID, req.Blocks)
	if err != nil {
		return nil, err
	}

	return page, nil
}
