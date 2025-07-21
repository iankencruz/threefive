package blocks

import (
	"github.com/google/uuid"
)

type BlockWithProps struct {
	ID        *uuid.UUID `json:"id"`
	Type      string     `json:"type"`
	SortOrder int32      `json:"sort_order"`
	Props     any        `json:"props"` // raw props for subtable
}
