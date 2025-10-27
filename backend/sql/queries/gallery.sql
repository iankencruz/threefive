-- name: CreateGallery :one
INSERT INTO galleries (title, description, created_at, updated_at)
VALUES (@title, @description, NOW(), NOW())
RETURNING *;

-- name: GetGalleryByID :one
SELECT * FROM galleries
WHERE id = @id;

-- name: ListGalleries :many
SELECT * FROM galleries
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: UpdateGallery :one
UPDATE galleries
SET 
    title = @title,
    description = @description,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteGallery :exec
DELETE FROM galleries
WHERE id = @id;

-- name: GetGalleryMediaCount :one
SELECT COUNT(*) as count
FROM media_relations
WHERE entity_type = 'gallery' AND entity_id = @entity_id;
