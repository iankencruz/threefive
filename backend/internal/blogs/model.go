package blogs

import (
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/generated"
)

// Wrapper struct for hydrated blogs

type BlogWithBlocks struct {
	Blog   generated.Blog           `json:"blog"`
	Blocks []*blocks.BlockWithProps `json:"blocks"` // ✅ updated
}
