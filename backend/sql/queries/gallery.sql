-- name: CreateGallery :one
INSERT INTO galleries (title, created_at, updated_at)
VALUES ($1, NOW(), NOW())
RETURNING *;

-- name: GetGalleryByID :one
SELECT * FROM galleries
WHERE id = $1;

-- name: ListGalleries :many
SELECT * FROM galleries
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateGallery :one
UPDATE galleries
SET title = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: DeleteGallery :exec
DELETE FROM galleries
WHERE id = $1;

-- name: CreateGalleryImage :one
INSERT INTO gallery_images (gallery_id, image_url, position, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: GetGalleryImages :many
SELECT * FROM gallery_images
WHERE gallery_id = $1
ORDER BY position ASC;

-- name: GetGalleryImageByID :one
SELECT * FROM gallery_images
WHERE id = $1;

-- name: DeleteGalleryImage :exec
DELETE FROM gallery_images
WHERE id = $1;

-- name: DeleteAllGalleryImages :exec
DELETE FROM gallery_images
WHERE gallery_id = $1;

-- name: UpdateImagePosition :exec
UPDATE gallery_images
SET position = $1
WHERE id = $2;
