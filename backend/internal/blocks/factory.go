// backend/internal/blocks/factory.go
package blocks

import (
	"encoding/json"
	"fmt"
)

// ParseBlockData parses block data map into specific block type
func ParseBlockData(blockType string, data map[string]interface{}) (interface{}, error) {
	switch blockType {
	case TypeHero:
		return parseHeroBlockData(data)
	case TypeRichtext:
		return parseRichtextBlockData(data)
	case TypeHeader:
		return parseHeaderBlockData(data)
	case TypeGallery:
		return parseGalleryBlockData(data)
	default:
		return nil, fmt.Errorf("unknown block type: %s", blockType)
	}
}

func parseHeroBlockData(data map[string]interface{}) (*HeroBlockData, error) {
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

func parseRichtextBlockData(data map[string]interface{}) (*RichtextBlockData, error) {
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

func parseHeaderBlockData(data map[string]interface{}) (*HeaderBlockData, error) {
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

func parseGalleryBlockData(data map[string]interface{}) (*GalleryBlockData, error) {
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
