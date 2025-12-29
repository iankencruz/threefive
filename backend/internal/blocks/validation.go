// backend/internal/blocks/validation.go
package blocks

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/validation"
)

// ============================================
// Block Request Validation
// ============================================

// ValidateBlockRequest validates a block request
func ValidateBlockRequest(v *validation.Validator, block *BlockRequest, fieldPrefix string) {
	// Type validation
	v.Required(fieldPrefix+".type", block.Type)
	v.In(fieldPrefix+".type", block.Type, ValidBlockTypes())

	// Data validation
	if block.Data == nil {
		v.AddError(fieldPrefix+".data", "Block data is required")
		return
	}

	// Type-specific validation
	switch block.Type {
	case TypeHero:
		validateHeroBlockData(v, block.Data, fieldPrefix)
	case TypeRichtext:
		validateRichtextBlockData(v, block.Data, fieldPrefix)
	case TypeHeader:
		validateHeaderBlockData(v, block.Data, fieldPrefix)
	case TypeGallery:
		validateGalleryBlockData(v, block.Data, fieldPrefix)
	case TypeFeature:
		validateFeatureBlockData(v, block.Data, fieldPrefix)
	}
}

// ValidateBlocks validates an array of block requests
func ValidateBlocks(v *validation.Validator, blocks []BlockRequest) {
	for i, block := range blocks {
		fieldPrefix := fmt.Sprintf("blocks[%d]", i)
		ValidateBlockRequest(v, &block, fieldPrefix)
	}
}

// ============================================
// Block Type-Specific Validation
// ============================================

func validateHeroBlockData(v *validation.Validator, data map[string]interface{}, fieldPrefix string) {
	// Title is required
	title, ok := data["title"].(string)
	if !ok || title == "" {
		v.AddError(fieldPrefix+".data.title", "Hero block title is required")
		return
	}
	v.MinLength(fieldPrefix+".data.title", title, 1)
	v.MaxLength(fieldPrefix+".data.title", title, 200)

	// CTA validation (both or neither)
	_, hasCtaText := data["cta_text"].(string)
	ctaURL, hasCtaURL := data["cta_url"].(string)

	if hasCtaText && !hasCtaURL {
		v.AddError(fieldPrefix+".data.cta_url", "CTA URL required when CTA text is provided")
	}
	if hasCtaURL && !hasCtaText {
		v.AddError(fieldPrefix+".data.cta_text", "CTA text required when CTA URL is provided")
	}

	// Validate URL format if provided
	if hasCtaURL && ctaURL != "" {
		v.URL(fieldPrefix+".data.cta_url", ctaURL)
	}
}

func validateRichtextBlockData(v *validation.Validator, data map[string]interface{}, fieldPrefix string) {
	content, ok := data["content"].(string)
	if !ok || content == "" {
		v.AddError(fieldPrefix+".data.content", "Richtext block content is required")
		return
	}
	v.MinLength(fieldPrefix+".data.content", content, 1)
}

func validateHeaderBlockData(v *validation.Validator, data map[string]interface{}, fieldPrefix string) {
	// Heading is required
	heading, ok := data["heading"].(string)
	if !ok || heading == "" {
		v.AddError(fieldPrefix+".data.heading", "Header block heading is required")
		return
	}
	v.MinLength(fieldPrefix+".data.heading", heading, 1)
	v.MaxLength(fieldPrefix+".data.heading", heading, 200)

	// Level validation (optional, defaults to h2)
	if level, ok := data["level"].(string); ok {
		v.In(fieldPrefix+".data.level", level, ValidHeaderLevels())
	}
}

func validateGalleryBlockData(v *validation.Validator, data map[string]any, fieldPrefix string) {
	// Media IDs validation
	mediaIDs, ok := data["media_ids"].([]any)
	if !ok {
		v.AddError(fieldPrefix+".data.media_ids", "Gallery block media_ids must be an array")
		return
	}

	if len(mediaIDs) == 0 {
		v.AddError(fieldPrefix+".data.media_ids", "Gallery block must have at least one media ID")
		return
	}

	// Validate each media ID is a valid UUID
	for i, id := range mediaIDs {
		idStr, ok := id.(string)
		if !ok {
			v.AddError(fmt.Sprintf("%s.data.media_ids[%d]", fieldPrefix, i), "Media ID must be a string")
			continue
		}
		if _, err := uuid.Parse(idStr); err != nil {
			v.AddError(fmt.Sprintf("%s.data.media_ids[%d]", fieldPrefix, i), "Invalid media ID format")
		}
	}
}

func validateFeatureBlockData(v *validation.Validator, data map[string]any, fieldPrefix string) {
	// Title is required
	title, ok := data["title"].(string)
	if !ok || title == "" {
		v.AddError(fieldPrefix+".data.title", "Feature block title is required")
		return
	}
	v.MinLength(fieldPrefix+".data.title", title, 1)
	v.MaxLength(fieldPrefix+".data.title", title, 200)

	// Description is required
	description, ok := data["description"].(string)
	if !ok || description == "" {
		v.AddError(fieldPrefix+".data.description", "Feature block description is required")
		return
	}
	v.MinLength(fieldPrefix+".data.description", description, 1)
	v.MaxLength(fieldPrefix+".data.description", description, 1000)

	// Heading is required
	heading, ok := data["heading"].(string)
	if !ok || heading == "" {
		v.AddError(fieldPrefix+".data.heading", "Feature block heading is required")
		return
	}
	v.MinLength(fieldPrefix+".data.heading", heading, 1)
	v.MaxLength(fieldPrefix+".data.heading", heading, 200)

	// Subheading is optional
	if subheading, ok := data["subheading"].(string); ok && subheading != "" {
		v.MaxLength(fieldPrefix+".data.subheading", subheading, 200)
	}

	// Media IDs validation
	mediaIDs, ok := data["media_ids"].([]any)
	if !ok {
		v.AddError(fieldPrefix+".data.media_ids", "Feature block media_ids must be an array")
		return
	}

	if len(mediaIDs) == 0 {
		v.AddError(fieldPrefix+".data.media_ids", "Feature block must have at least one media ID")
		return
	}

	// Validate each media ID is a valid UUID
	for i, id := range mediaIDs {
		idStr, ok := id.(string)
		if !ok {
			v.AddError(fmt.Sprintf("%s.data.media_ids[%d]", fieldPrefix, i), "Media ID must be a string")
			continue
		}
		if _, err := uuid.Parse(idStr); err != nil {
			v.AddError(fmt.Sprintf("%s.data.media_ids[%d]", fieldPrefix, i), "Invalid media ID format")
		}
	}
}
