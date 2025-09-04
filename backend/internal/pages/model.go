package pages

import "github.com/iankencruz/threefive/internal/generated"

type PageWithGalleries struct {
	Page      generated.Page
	Galleries []GalleryWithMedia
}

type GalleryWithMedia struct {
	Gallery generated.Gallery
	Media   []generated.Media
}
