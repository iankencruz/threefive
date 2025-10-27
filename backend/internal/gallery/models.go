package gallery

import (
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
)

type CreateGalleryRequest struct {
	Title       string      `json:"title" validate:"required"`
	Description string      `json:"description"`
	MediaIDs    []uuid.UUID `json:"media_ids"`
}

type UpdateGalleryRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	MediaIDs    []uuid.UUID `json:"media_ids"`
}

type LinkMediaRequest struct {
	MediaID   uuid.UUID `json:"media_id" validate:"required"`
	SortOrder int32     `json:"sort_order"`
}

type GalleryWithMedia struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Media       []sqlc.Media `json:"media"`
	MediaCount  int64        `json:"media_count"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type GalleryListItem struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	MediaCount  int64     `json:"media_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
