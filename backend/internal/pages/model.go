package pages

import "github.com/iankencruz/threefive/internal/generated"

type PageWithGalleries struct {
	Page      generated.Page     `json:"page"`
	Galleries []GalleryWithMedia `json:"galleries"`
}

// Gallery with media embedded inside the gallery object
type GalleryWithMedia struct {
	Gallery GalleryWithEmbeddedMedia `json:"gallery"`
}

// Extended gallery struct that includes media as a property
type GalleryWithEmbeddedMedia struct {
	generated.Gallery                   // Embed all original gallery fields
	Media             []generated.Media `json:"media"` // Media nested inside gallery
}
