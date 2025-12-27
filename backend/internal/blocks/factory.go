// backend/internal/blocks/factory.go
package blocks

import (
	"encoding/json"
	"fmt"
)

// ParseBlockData parses block data map into specific block type
func ParseBlockData(blockType string, data map[string]any) (any, error) {
	switch blockType {
	case TypeHero:
		return parseHeroBlockData(data)
	case TypeRichtext:
		return parseRichtextBlockData(data)
	case TypeHeader:
		return parseHeaderBlockData(data)
	case TypeGallery:
		return parseGalleryBlockData(data)
	case TypeFeature:
		return parseFeatureBlockData(data)
	default:
		return nil, fmt.Errorf("unknown block type: %s", blockType)
	}
}

func parseHeroBlockData(data map[string]any) (*HeroBlockData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal hero block data: %w", err)
	}

	var heroData HeroBlockData
	if err := json.Unmarshal(jsonData, &heroData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hero block data: %w", err)
	}

	return &heroData, nil
}

func parseRichtextBlockData(data map[string]any) (*RichtextBlockData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal richtext block data: %w", err)
	}

	var richtextData RichtextBlockData
	if err := json.Unmarshal(jsonData, &richtextData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal richtext block data: %w", err)
	}

	return &richtextData, nil
}

func parseHeaderBlockData(data map[string]any) (*HeaderBlockData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal header block data: %w", err)
	}

	var headerData HeaderBlockData
	if err := json.Unmarshal(jsonData, &headerData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal header block data: %w", err)
	}

	// Default level to h2 if not provided
	if headerData.Level == "" {
		headerData.Level = "h2"
	}

	return &headerData, nil
}

func parseGalleryBlockData(data map[string]any) (*GalleryBlockData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gallery block data: %w", err)
	}

	var galleryData GalleryBlockData
	if err := json.Unmarshal(jsonData, &galleryData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gallery block data: %w", err)
	}

	return &galleryData, nil
}

func parseFeatureBlockData(data map[string]any) (*FeatureBlockData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal about me block data: %w", err)
	}
	var aboutData FeatureBlockData
	if err := json.Unmarshal(jsonData, &aboutData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal about me block data: %w", err)
	}
	return &aboutData, nil
}
