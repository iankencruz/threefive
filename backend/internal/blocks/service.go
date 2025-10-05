// backend/internal/blocks/service.go
package blocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service handles block operations
type Service struct {
	queries *sqlc.Queries
}

// NewService creates a new blocks service
func NewService(queries *sqlc.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// CreateBlocks creates all blocks for a page in a transaction
func (s *Service) CreateBlocks(ctx context.Context, qtx *sqlc.Queries, pageID uuid.UUID, blocks []BlockRequest) error {
	for i, blockReq := range blocks {
		// Create base block
		block, err := qtx.CreateBlock(ctx, sqlc.CreateBlockParams{
			PageID:    pageID,
			Type:      blockReq.Type,
			SortOrder: int32(i),
		})
		if err != nil {
			return errors.Internal("Failed to create block", err)
		}

		// Create type-specific block data
		switch blockReq.Type {
		case TypeHero:
			if err := s.createHeroBlock(ctx, qtx, block.ID, blockReq.Data); err != nil {
				return err
			}
		case TypeRichtext:
			if err := s.createRichtextBlock(ctx, qtx, block.ID, blockReq.Data); err != nil {
				return err
			}
		case TypeHeader:
			if err := s.createHeaderBlock(ctx, qtx, block.ID, blockReq.Data); err != nil {
				return err
			}
		default:
			return errors.BadRequest("Invalid block type", "invalid_block_type")
		}
	}

	return nil
}

// GetPageBlocks retrieves all blocks for a page
func (s *Service) GetPageBlocks(ctx context.Context, pageID uuid.UUID) ([]BlockResponse, error) {
	// Get base blocks
	baseBlocks, err := s.queries.GetBlocksByPageID(ctx, pageID)
	if err != nil {
		return nil, errors.Internal("Failed to get blocks", err)
	}

	if len(baseBlocks) == 0 {
		return []BlockResponse{}, nil
	}

	// Get all block type data
	heroBlocks, err := s.queries.GetHeroBlocksByPageID(ctx, pageID)
	if err != nil {
		return nil, errors.Internal("Failed to get hero blocks", err)
	}

	richtextBlocks, err := s.queries.GetRichtextBlocksByPageID(ctx, pageID)
	if err != nil {
		return nil, errors.Internal("Failed to get richtext blocks", err)
	}

	headerBlocks, err := s.queries.GetHeaderBlocksByPageID(ctx, pageID)
	if err != nil {
		return nil, errors.Internal("Failed to get header blocks", err)
	}

	// Build lookup maps
	heroMap := make(map[uuid.UUID]sqlc.BlockHero)
	for _, h := range heroBlocks {
		heroMap[h.BlockID] = h
	}

	richtextMap := make(map[uuid.UUID]sqlc.BlockRichtext)
	for _, r := range richtextBlocks {
		richtextMap[r.BlockID] = r
	}

	headerMap := make(map[uuid.UUID]sqlc.BlockHeader)
	for _, h := range headerBlocks {
		headerMap[h.BlockID] = h
	}

	// Assemble response
	blocks := make([]BlockResponse, 0, len(baseBlocks))

	for _, block := range baseBlocks {
		blockResp := BlockResponse{
			ID:        block.ID,
			Type:      block.Type,
			SortOrder: int(block.SortOrder),
		}

		// Attach type-specific data
		switch block.Type {
		case TypeHero:
			if hero, ok := heroMap[block.ID]; ok {
				blockResp.Data = HeroBlockData{
					Title:    hero.Title,
					Subtitle: nullTextToPtr(hero.Subtitle),
					ImageID:  nullUUIDToPtr(hero.ImageID),
					CtaText:  nullTextToPtr(hero.CtaText),
					CtaURL:   nullTextToPtr(hero.CtaUrl),
				}
			}
		case TypeRichtext:
			if richtext, ok := richtextMap[block.ID]; ok {
				blockResp.Data = RichtextBlockData{
					Content: richtext.Content,
				}
			}
		case TypeHeader:
			if header, ok := headerMap[block.ID]; ok {
				blockResp.Data = HeaderBlockData{
					Heading:    header.Heading,
					Subheading: nullTextToPtr(header.Subheading),
					Level:      header.Level,
				}
			}
		}

		blocks = append(blocks, blockResp)
	}

	return blocks, nil
}

// ============================================
// Private Helper Methods
// ============================================

func (s *Service) createHeroBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]interface{}) error {
	heroData, err := ParseBlockData(TypeHero, data)
	if err != nil {
		return errors.BadRequest("Invalid hero block data", "invalid_block_data")
	}

	hero := heroData.(*HeroBlockData)

	_, err = qtx.CreateHeroBlock(ctx, sqlc.CreateHeroBlockParams{
		BlockID:  blockID,
		Title:    hero.Title,
		Subtitle: strToNullText(hero.Subtitle),
		ImageID:  uuidToNullUUID(hero.ImageID),
		CtaText:  strToNullText(hero.CtaText),
		CtaUrl:   strToNullText(hero.CtaURL),
	})
	if err != nil {
		return errors.Internal("Failed to create hero block", err)
	}

	return nil
}

func (s *Service) createRichtextBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]interface{}) error {
	richtextData, err := ParseBlockData(TypeRichtext, data)
	if err != nil {
		return errors.BadRequest("Invalid richtext block data", "invalid_block_data")
	}

	richtext := richtextData.(*RichtextBlockData)

	_, err = qtx.CreateRichtextBlock(ctx, sqlc.CreateRichtextBlockParams{
		BlockID: blockID,
		Content: richtext.Content,
	})
	if err != nil {
		return errors.Internal("Failed to create richtext block", err)
	}

	return nil
}

func (s *Service) createHeaderBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]interface{}) error {
	headerData, err := ParseBlockData(TypeHeader, data)
	if err != nil {
		return errors.BadRequest("Invalid header block data", "invalid_block_data")
	}

	header := headerData.(*HeaderBlockData)

	_, err = qtx.CreateHeaderBlock(ctx, sqlc.CreateHeaderBlockParams{
		BlockID:    blockID,
		Heading:    header.Heading,
		Subheading: strToNullText(header.Subheading),
		Level:      header.Level,
	})
	if err != nil {
		return errors.Internal("Failed to create header block", err)
	}

	return nil
}

// ============================================
// Helper Functions for Nullable Types
// ============================================

func strToNullText(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func uuidToNullUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func nullTextToPtr(nt pgtype.Text) *string {
	if !nt.Valid {
		return nil
	}
	return &nt.String
}

func nullUUIDToPtr(nu pgtype.UUID) *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	return &nu.Bytes
}
