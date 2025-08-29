package blogs

import (
	"github.com/iankencruz/threefive/internal/generated"
)

// Wrapper struct for hydrated blogs

type BlogWithMedia struct {
	*generated.Blog
	Media []*generated.Media `json:"media"`
}
