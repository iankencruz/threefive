// backend/internal/blocks/types.go
package blocks

import "github.com/google/uuid"

// ============================================
// Block Type Constants
// ============================================

const (
	TypeHero     = "hero"
	TypeRichtext = "richtext"
	TypeHeader   = "header"
	TypeGallery  = "gallery"
	TypeAbout    = "about"
)

// ValidBlockTypes returns all valid block types
func ValidBlockTypes() []string {
	return []string{TypeHero, TypeRichtext, TypeHeader, TypeGallery, TypeAbout}
}

// ============================================
// Block-Specific Data Models
// ============================================

// HeroBlockData represents hero block data
type HeroBlockData struct {
	Title    string     `json:"title"`
	Subtitle *string    `json:"subtitle,omitempty"`
	ImageID  *uuid.UUID `json:"image_id,omitempty"`
	CtaText  *string    `json:"cta_text,omitempty"`
	CtaURL   *string    `json:"cta_url,omitempty"`
}

// RichtextBlockData represents richtext block data
type RichtextBlockData struct {
	Content string `json:"content"`
}

// HeaderBlockData represents header block data
type HeaderBlockData struct {
	Heading    string  `json:"heading"`
	Subheading *string `json:"subheading,omitempty"`
	Level      string  `json:"level"` // h1, h2, h3, h4, h5, h6
}

// GalleryBlockData represents gallery block data
type GalleryBlockData struct {
	Title    *string     `json:"title,omitempty"`
	MediaIDs []uuid.UUID `json:"media_ids"`
}

// AboutBlockData represents about me block data
type AboutBlockData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Heading     string `json:"heading"`
	Subheading  string `json:"subheading,omitempty"`
}

// ValidHeaderLevels returns all valid header levels...
func ValidHeaderLevels() []string {
	return []string{"h1", "h2", "h3", "h4", "h5", "h6"}
}
