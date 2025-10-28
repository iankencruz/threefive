package blocks

import (
	"github.com/google/uuid"
)

// ============================================
// Request Models
// ============================================

// BlockRequest represents a block in create/update requests
type BlockRequest struct {
	ID   *uuid.UUID     `json:"id,omitempty"` // For updates
	Type string         `json:"type"`         // hero, richtext, header
	Data map[string]any `json:"data"`
}

// ============================================
// Response Models
// ============================================

// BlockResponse represents a block in responses
type BlockResponse struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	SortOrder int       `json:"sort_order"`
	Data      any       `json:"data"` // HeroBlockData, RichtextBlockData, or HeaderBlockData
}
