package gallery

import "github.com/iankencruz/threefive/internal/generated"

type GalleryWithMedia struct {
	*generated.Gallery
	Media []*generated.Media `json:"media"`
}
