// backend/internal/blocks/service.go
package blocks

import (
	"context"
	"fmt"
	"log"

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
func (s *Service) CreateBlocks(ctx context.Context, qtx *sqlc.Queries, entityType string, entityID uuid.UUID, blocks []BlockRequest) error {
	for i, blockReq := range blocks {
		// Create base block
		block, err := qtx.CreateBlock(ctx, sqlc.CreateBlockParams{
			EntityType: entityType,
			EntityID:   entityID,
			Type:       blockReq.Type,
			SortOrder:  int32(i),
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
		case TypeGallery:
			if err := s.createGalleryBlock(ctx, qtx, block.ID, blockReq.Data); err != nil {
				return err
			}
		case TypeFeature:
			if err := s.createFeatureBlock(ctx, qtx, block.ID, blockReq.Data); err != nil {
				return err
			}
		default:
			return errors.BadRequest("Invalid block type", "invalid_block_type")
		}
	}

	return nil
}

// UpdateBlocks updates blocks for a page in a transaction
func (s *Service) UpdateBlocks(ctx context.Context, qtx *sqlc.Queries, entityType string, entityID uuid.UUID, blocks []BlockRequest) error {
	// Get existing blocks to determine what to update/delete/create
	// existingBlocks, err := qtx.GetBlocksByPageID(ctx, pageID)
	existingBlocks, err := qtx.GetBlocksByEntity(ctx, sqlc.GetBlocksByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return errors.Internal("Failed to get existing blocks", err)
	}

	// Create a map of existing block IDs
	existingBlockMap := make(map[uuid.UUID]bool)
	for _, eb := range existingBlocks {
		existingBlockMap[eb.ID] = true
	}

	// Track which blocks are being kept/updated
	updatedBlockIDs := make(map[uuid.UUID]bool)

	// Process each block in the request
	for i, block := range blocks {
		if block.ID != nil {
			// Update existing block
			updatedBlockIDs[*block.ID] = true

			if err := qtx.UpdateBlockOrder(ctx, sqlc.UpdateBlockOrderParams{
				SortOrder: int32(i),
				ID:        *block.ID,
			}); err != nil {
				return errors.Internal("Failed to update block order", err)
			}

			// Update type-specific data
			switch block.Type {
			case TypeHero:
				if err := s.updateHeroBlock(ctx, qtx, *block.ID, block.Data); err != nil {
					return err
				}
			case TypeRichtext:
				if err := s.updateRichtextBlock(ctx, qtx, *block.ID, block.Data); err != nil {
					return err
				}
			case TypeHeader:
				if err := s.updateHeaderBlock(ctx, qtx, *block.ID, block.Data); err != nil {
					return err
				}
			case TypeGallery:
				if err := s.updateGalleryBlock(ctx, qtx, *block.ID, block.Data); err != nil {
					return err
				}
			case TypeFeature:
				if err := s.updateFeatureBlock(ctx, qtx, *block.ID, block.Data); err != nil {
					return err
				}
			default:
				return errors.BadRequest("Invalid block type", "invalid_block_type")
			}
		} else {
			// Create new block
			newBlock, err := qtx.CreateBlock(ctx, sqlc.CreateBlockParams{
				EntityType: entityType,
				EntityID:   entityID,
				Type:       block.Type,
				SortOrder:  int32(i),
			})
			if err != nil {
				return errors.Internal("Failed to create block", err)
			}

			// Create type-specific block data
			switch block.Type {
			case TypeHero:
				if err := s.createHeroBlock(ctx, qtx, newBlock.ID, block.Data); err != nil {
					return err
				}
			case TypeRichtext:
				if err := s.createRichtextBlock(ctx, qtx, newBlock.ID, block.Data); err != nil {
					return err
				}
			case TypeHeader:
				if err := s.createHeaderBlock(ctx, qtx, newBlock.ID, block.Data); err != nil {
					return err
				}
			case TypeGallery:
				if err := s.createGalleryBlock(ctx, qtx, newBlock.ID, block.Data); err != nil {
					return err
				}
			default:
				return errors.BadRequest("Invalid block type", "invalid_block_type")
			}
		}
	}

	// Delete blocks that are no longer in the request
	for _, existingBlock := range existingBlocks {
		if !updatedBlockIDs[existingBlock.ID] {
			if err := qtx.DeleteBlock(ctx, existingBlock.ID); err != nil {
				return errors.Internal("Failed to delete block", err)
			}
		}
	}

	return nil
}

// GetPageBlocks retrieves all blocks for a page
func (s *Service) GetBlocksByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]BlockResponse, error) {
	// Get base blocks
	baseBlocks, err := s.queries.GetBlocksByEntity(ctx, sqlc.GetBlocksByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get blocks", err)
	}

	if len(baseBlocks) == 0 {
		return []BlockResponse{}, nil
	}

	// ***********************
	// Get all block type data

	// Hero Block
	heroBlocks, err := s.queries.GetHeroBlocksByEntity(ctx, sqlc.GetHeroBlocksByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get hero blocks", err)
	}

	// RichText Block
	richtextBlocks, err := s.queries.GetRichtextBlocksByEntity(ctx, sqlc.GetRichtextBlocksByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get richtext blocks", err)
	}

	// Header Block
	headerBlocks, err := s.queries.GetHeaderBlocksByEntity(ctx, sqlc.GetHeaderBlocksByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get header blocks", err)
	}

	// Gallery Block
	galleryBlocks, err := s.queries.GetGalleryBlocksByEntity(ctx, sqlc.GetGalleryBlocksByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get gallery blocks", err)
	}

	aboutBlocks, err := s.queries.GetFeatureBlocksByEntity(ctx, sqlc.GetFeatureBlocksByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get about me blocks", err)
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

	galleryMap := make(map[uuid.UUID]sqlc.BlockGallery)
	for _, gallery := range galleryBlocks {
		galleryMap[gallery.BlockID] = gallery
	}

	featureMap := make(map[uuid.UUID]sqlc.BlockFeature)
	for _, about := range aboutBlocks {
		featureMap[about.BlockID] = about
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
		// backend/internal/blocks/service.go - in GetBlocksByEntity method
		case TypeHero:
			if hero, ok := heroMap[block.ID]; ok {
				heroData := HeroBlockData{
					Title:    hero.Title,
					Subtitle: nullTextToPtr(hero.Subtitle),
					ImageID:  nullUUIDToPtr(hero.ImageID),
					CtaText:  nullTextToPtr(hero.CtaText),
					CtaURL:   nullTextToPtr(hero.CtaUrl),
				}

				// ✨ NEW: Fetch the actual media if ImageID exists
				if hero.ImageID.Valid {
					mediaID := uuid.UUID(hero.ImageID.Bytes)
					media, err := s.queries.GetMediaByID(ctx, mediaID)
					if err == nil {
						// Add media to the response as a map for flexibility
						blockResp.Data = map[string]any{
							"title":    heroData.Title,
							"subtitle": heroData.Subtitle,
							"image_id": heroData.ImageID,
							"cta_text": heroData.CtaText,
							"cta_url":  heroData.CtaURL,
							"media":    media, // ← Add the full media object!
						}
					} else {
						// If media fetch fails, use struct without media
						blockResp.Data = heroData
					}
				} else {
					blockResp.Data = heroData
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
					Level:      pgTextToString(header.Level),
				}
			}
		case TypeFeature:
			if feature, ok := featureMap[block.ID]; ok {
				featureData := FeatureBlockData{
					Title:       feature.Title,
					Description: feature.Description,
					Heading:     feature.Heading,
					Subheading:  feature.Subheading,
					ImageID:     nullUUIDToPtr(feature.ImageID),
				}

				// Fetch the actual media if ImageID exists (same pattern as Hero)
				if feature.ImageID.Valid {
					mediaID := uuid.UUID(feature.ImageID.Bytes)
					media, err := s.queries.GetMediaByID(ctx, mediaID)
					if err == nil {
						blockResp.Data = map[string]any{
							"title":       featureData.Title,
							"description": featureData.Description,
							"heading":     featureData.Heading,
							"subheading":  featureData.Subheading,
							"image_id":    featureData.ImageID,
							"media":       media,
						}
					} else {
						blockResp.Data = featureData
					}
				} else {
					blockResp.Data = featureData
				}
			}
		case TypeGallery:
			if gallery, ok := galleryMap[block.ID]; ok {
				// Fetch media for this gallery
				media, err := s.queries.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
					EntityType: "block_gallery",
					EntityID:   gallery.ID,
				})

				// Always set data structure, even if media fetch fails or is empty
				if err != nil {
					// Log error but continue with empty media
					log.Printf("Warning: Failed to fetch media for gallery block %s: %v", gallery.ID, err)
					blockResp.Data = map[string]any{
						"title": nullTextToPtr(gallery.Title),
						"media": []sqlc.Media{}, // Empty array instead of nil
					}
				} else {
					// Ensure media is never nil
					if media == nil {
						media = []sqlc.Media{}
					}
					blockResp.Data = map[string]any{
						"title": nullTextToPtr(gallery.Title),
						"media": media,
					}
				}
			}
		}

		blocks = append(blocks, blockResp)
	}

	return blocks, nil
}

// ============================================
// Private Helper Methods - Create
// ============================================

func (s *Service) createHeroBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	heroData, err := ParseBlockData(TypeHero, data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
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

func (s *Service) createRichtextBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
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

func (s *Service) createHeaderBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	headerData, err := ParseBlockData(TypeHeader, data)
	if err != nil {
		return errors.BadRequest("Invalid header block data", "invalid_block_data")
	}

	header := headerData.(*HeaderBlockData)

	_, err = qtx.CreateHeaderBlock(ctx, sqlc.CreateHeaderBlockParams{
		BlockID:    blockID,
		Heading:    header.Heading,
		Subheading: strToNullText(header.Subheading),
		Level:      stringToPgText(header.Level),
	})
	if err != nil {
		return errors.Internal("Failed to create header block", err)
	}

	return nil
}

func (s *Service) createGalleryBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	galleryData, err := ParseBlockData(TypeGallery, data)
	if err != nil {
		return errors.BadRequest("Invalid gallery block data", "invalid_block_data")
	}

	gallery := galleryData.(*GalleryBlockData)

	// Create gallery block
	galleryBlock, err := qtx.CreateGalleryBlock(ctx, sqlc.CreateGalleryBlockParams{
		BlockID: blockID,
		Title:   strToNullText(gallery.Title),
	})
	if err != nil {
		return errors.Internal("Failed to create gallery block", err)
	}

	// Link media to gallery block
	for i, mediaID := range gallery.MediaIDs {
		_, err := qtx.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
			MediaID:    mediaID,
			EntityType: "block_gallery",
			EntityID:   galleryBlock.ID,
			SortOrder:  pgtype.Int4{Int32: int32(i), Valid: true},
		})
		if err != nil {
			return errors.Internal("Failed to link media to gallery block", err)
		}
	}

	return nil
}

func (s *Service) createFeatureBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	featureData, err := ParseBlockData(TypeFeature, data)
	if err != nil {
		return errors.BadRequest("Invalid about me block data", "invalid_block_data")
	}

	feature := featureData.(*FeatureBlockData)

	// create feature block
	_, err = qtx.CreateFeatureBlock(ctx, sqlc.CreateFeatureBlockParams{
		BlockID:     blockID,
		Title:       feature.Title,
		Description: feature.Description,
		Heading:     feature.Heading,
		Subheading:  feature.Subheading,
		ImageID:     uuidToNullUUID(feature.ImageID),
	})
	if err != nil {
		return errors.Internal("Failed to create Feature block", err)
	}
	return nil
}

// ============================================
// Private Helper Methods - Update
// ============================================

func (s *Service) updateHeroBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	heroData, err := ParseBlockData(TypeHero, data)
	if err != nil {
		return errors.BadRequest("Invalid hero block data", "invalid_block_data")
	}

	hero := heroData.(*HeroBlockData)

	_, err = qtx.UpdateHeroBlock(ctx, sqlc.UpdateHeroBlockParams{
		BlockID:  blockID,
		Title:    hero.Title,
		Subtitle: strToNullText(hero.Subtitle),
		ImageID:  uuidToNullUUID(hero.ImageID),
		CtaText:  strToNullText(hero.CtaText),
		CtaUrl:   strToNullText(hero.CtaURL),
	})
	if err != nil {
		return errors.Internal("Failed to update hero block", err)
	}

	return nil
}

func (s *Service) updateRichtextBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	richtextData, err := ParseBlockData(TypeRichtext, data)
	if err != nil {
		return errors.BadRequest("Invalid richtext block data", "invalid_block_data")
	}

	richtext := richtextData.(*RichtextBlockData)

	_, err = qtx.UpdateRichtextBlock(ctx, sqlc.UpdateRichtextBlockParams{
		BlockID: blockID,
		Content: richtext.Content,
	})
	if err != nil {
		return errors.Internal("Failed to update richtext block", err)
	}

	return nil
}

func (s *Service) updateHeaderBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	headerData, err := ParseBlockData(TypeHeader, data)
	if err != nil {
		return errors.BadRequest("Invalid header block data", "invalid_block_data")
	}

	header := headerData.(*HeaderBlockData)

	_, err = qtx.UpdateHeaderBlock(ctx, sqlc.UpdateHeaderBlockParams{
		BlockID:    blockID,
		Heading:    header.Heading,
		Subheading: strToNullText(header.Subheading),
		Level:      stringToPgText(header.Level),
	})
	if err != nil {
		return errors.Internal("Failed to update header block", err)
	}

	return nil
}

func (s *Service) updateGalleryBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	galleryData, err := ParseBlockData(TypeGallery, data)
	if err != nil {
		return errors.BadRequest("Invalid gallery block data", "invalid_block_data")
	}

	gallery := galleryData.(*GalleryBlockData)

	// Get existing gallery block
	existingGallery, err := qtx.GetGalleryBlockByBlockID(ctx, blockID)
	if err != nil {
		return errors.Internal("Failed to get gallery block", err)
	}

	// Update gallery block
	_, err = qtx.UpdateGalleryBlock(ctx, sqlc.UpdateGalleryBlockParams{
		BlockID: blockID,
		Title:   strToNullText(gallery.Title),
	})
	if err != nil {
		return errors.Internal("Failed to update gallery block", err)
	}

	// Update media links - remove old ones
	existingMedia, err := qtx.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
		EntityType: "block_gallery",
		EntityID:   existingGallery.ID,
	})
	if err != nil {
		return errors.Internal("Failed to get existing media", err)
	}

	// Unlink all existing media
	for _, media := range existingMedia {
		err := qtx.UnlinkMediaFromEntity(ctx, sqlc.UnlinkMediaFromEntityParams{
			MediaID:    media.ID,
			EntityType: "block_gallery",
			EntityID:   existingGallery.ID,
		})
		if err != nil {
			return errors.Internal("Failed to unlink media", err)
		}
	}

	// Link new media
	for i, mediaID := range gallery.MediaIDs {
		_, err := qtx.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
			MediaID:    mediaID,
			EntityType: "block_gallery",
			EntityID:   existingGallery.ID,
			SortOrder:  pgtype.Int4{Int32: int32(i), Valid: true},
		})
		if err != nil {
			return errors.Internal("Failed to link media to gallery block", err)
		}
	}

	return nil
}

// updateFeatureBlock()
func (s *Service) updateFeatureBlock(ctx context.Context, qtx *sqlc.Queries, blockID uuid.UUID, data map[string]any) error {
	featureData, err := ParseBlockData(TypeFeature, data)
	if err != nil {
		return errors.BadRequest("Invalid about me block data", "invalid_block_data")
	}

	feature := featureData.(*FeatureBlockData)
	_, err = qtx.UpdateFeatureBlock(ctx, sqlc.UpdateFeatureBlockParams{
		BlockID:     blockID,
		Title:       feature.Title,
		Description: feature.Description,
		Heading:     feature.Heading,
		Subheading:  feature.Subheading,
		ImageID:     uuidToNullUUID(feature.ImageID),
	})
	if err != nil {
		return errors.Internal("Failed to update about me block", err)
	}
	return nil
}

// ============================================
// Helper Functions for Nullable Types
// ============================================

func stringToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func pgTextToString(pt pgtype.Text) string {
	if !pt.Valid {
		return ""
	}
	return pt.String
}

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
	id := uuid.UUID(nu.Bytes)
	return &id
}
