package blocks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/generated"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func decodeMap[T any](input any) (T, error) {
	var out T
	data, err := json.Marshal(input)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(data, &out)
	return out, err
}

// ========== CRUD ==========

func (s *Service) CreateBlock(ctx context.Context, block generated.Block, props any) error {
	fmt.Printf("➕ CreateBlock: type=%s id=%s\n", block.Type, block.ID)

	_, err := s.Repo.q.CreateBlock(ctx, generated.CreateBlockParams{
		ID:         block.ID,
		ParentType: block.ParentType,
		ParentID:   block.ParentID,
		Type:       block.Type,
		SortOrder:  block.SortOrder,
	})
	if err != nil {
		fmt.Printf("❌ failed to insert base block: %v\n", err)
		return fmt.Errorf("failed to insert base block: %w", err)
	}
	switch block.Type {

	case "heading":
		params, err := decodeMap[generated.CreateHeadingBlockParams](props)
		if err != nil {
			return fmt.Errorf("invalid heading props payload format: %w", err)
		}
		params.BlockID = block.ID
		return s.Repo.CreateHeadingBlock(ctx, params)

	case "richtext":
		params, err := decodeMap[generated.CreateRichTextBlockParams](props)
		if err != nil {
			return fmt.Errorf("invalid richtext props payload format: %w", err)
		}
		params.BlockID = block.ID
		return s.Repo.CreateRichTextBlock(ctx, params)

	case "image":
		params, err := decodeMap[generated.CreateImageBlockParams](props)
		if err != nil {
			return fmt.Errorf("invalid richtext props payload format: %w", err)
		}
		params.BlockID = block.ID
		return s.Repo.CreateImageBlock(ctx, params)

	default:
		return fmt.Errorf("unsupported block type: %s", block.Type)
	}
}

func (s *Service) UpdateBlock(ctx context.Context, block generated.Block, props any) error {
	fmt.Printf("✏️ UpdateBlock: type=%s id=%s\n", block.Type, block.ID)

	err := s.Repo.q.UpdateBlock(ctx, generated.UpdateBlockParams{
		ID:         block.ID,
		ParentType: block.ParentType,
		ParentID:   block.ParentID,
		Type:       block.Type,
		SortOrder:  int32(block.SortOrder),
	})
	if err != nil {
		fmt.Printf("❌ update base block failed: %v\n", err)
		return err
	}

	switch block.Type {
	case "heading":
		params, err := decodeMap[generated.UpdateHeadingBlockParams](props)
		if err != nil {
			return fmt.Errorf("invalid heading props payload format: %w", err)
		}
		params.BlockID = block.ID
		return s.Repo.UpdateHeadingBlock(ctx, params)

	case "richtext":
		params, err := decodeMap[generated.UpdateRichTextBlockParams](props)
		if err != nil {
			return fmt.Errorf("invalid richtext props payload format: %w", err)
		}
		params.BlockID = block.ID
		return s.Repo.UpdateRichTextBlock(ctx, params)

	case "image":
		params, err := decodeMap[generated.UpdateImageBlockParams](props)
		if err != nil {
			return fmt.Errorf("invalid image props payload format: %w", err)
		}
		params.BlockID = block.ID
		// if params.media_id == "" then convert it to nil
		return s.Repo.UpdateImageBlock(ctx, params)

	default:
		return fmt.Errorf("unsupported block type update: %s", block.Type)
	}
}

func (s *Service) DeleteBlockByID(ctx context.Context, id uuid.UUID, blockType string) error {
	// delete sub-table content first
	switch blockType {
	case "heading":
		if err := s.Repo.DeleteHeadingBlockByID(ctx, id); err != nil {
			return err
		}
	case "richtext":
		if err := s.Repo.DeleteRichTextBlockByID(ctx, id); err != nil {
			return err
		}
	case "image":
		if err := s.Repo.DeleteImageBlock(ctx, id); err != nil {
			return err
		}
	default:
		return errors.New("unsupported block delete: " + blockType)
	}

	// delete base block last
	return s.Repo.q.DeleteBlock(ctx, id)
}

func (s *Service) DeleteBlocksByParent(ctx context.Context, parentType string, parentID uuid.UUID) error {
	return s.Repo.q.DeleteBlocksByParent(ctx, generated.DeleteBlocksByParentParams{
		ParentType: parentType,
		ParentID:   parentID,
	})
}

func (s *Service) UpdateBlockSortOrder(ctx context.Context, id uuid.UUID, sortOrder int32) error {
	return s.Repo.q.UpdateBlockSortOrder(ctx, generated.UpdateBlockSortOrderParams{
		ID:        id,
		SortOrder: sortOrder,
	})
}

func (s *Service) SyncBlocks(ctx context.Context, parentType string, parentID uuid.UUID, incoming []BlockWithProps) error {
	existing, err := s.Repo.GetBaseBlocksFromPage(ctx, parentID)
	if err != nil {
		return err
	}

	existingMap := make(map[uuid.UUID]bool)
	for _, b := range existing {
		existingMap[b.ID] = true
	}

	incomingMap := make(map[uuid.UUID]BlockWithProps)

	fmt.Printf("🔁 SyncBlocks: incoming blocks (%d)\n", len(incoming))

	for _, b := range incoming {
		if b.ID == nil || *b.ID == uuid.Nil {
			newID := uuid.New()
			b.ID = &newID
		}

		incomingMap[*b.ID] = b

		block := generated.Block{
			ID:         *b.ID,
			ParentType: parentType,
			ParentID:   parentID,
			Type:       b.Type,
			SortOrder:  b.SortOrder, // ✅ Use latest
		}

		fmt.Printf("📦 Processing block: ID=%s Type=%s SortOrder=%d Props=%#v\n", b.ID.String(), b.Type, b.SortOrder, b.Props)

		// Base block exists?
		_, err := s.Repo.q.GetBlockByID(ctx, block.ID)
		if err != nil {
			// Create base + sub-table
			if err := s.CreateBlock(ctx, block, b.Props); err != nil {
				return fmt.Errorf("create block failed: %w", err)
			}
			continue
		}

		// Sub-table existence check
		switch b.Type {
		case "heading":
			if _, err := s.Repo.q.GetHeadingBlock(ctx, block.ID); err != nil {
				params, err := decodeMap[generated.CreateHeadingBlockParams](b.Props)
				if err != nil {
					return fmt.Errorf("invalid heading props payload: %w", err)
				}
				params.BlockID = block.ID
				if err := s.Repo.CreateHeadingBlock(ctx, params); err != nil {
					return fmt.Errorf("create heading block failed: %w", err)
				}
			}
		case "richtext":
			if _, err := s.Repo.q.GetRichTextBlock(ctx, block.ID); err != nil {
				params, err := decodeMap[generated.CreateRichTextBlockParams](b.Props)
				if err != nil {
					return fmt.Errorf("invalid richtext props payload: %w", err)
				}
				params.BlockID = block.ID
				if err := s.Repo.CreateRichTextBlock(ctx, params); err != nil {
					return fmt.Errorf("create richtext block failed: %w", err)
				}
			}
		}

		// ✅ Always update the full block (including latest SortOrder)
		if err := s.UpdateBlock(ctx, block, b.Props); err != nil {
			return fmt.Errorf("update block failed: %w", err)
		}
	}

	// Delete removed blocks
	for _, b := range existing {
		if _, found := incomingMap[b.ID]; !found {
			if err := s.DeleteBlockByID(ctx, b.ID, b.Type); err != nil {
				return fmt.Errorf("delete block failed: %w", err)
			}
		}
	}

	return nil
}
