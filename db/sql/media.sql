
-- name: CreateMedia :one
INSERT INTO media (
  url, thumbnail_url, medium_url, type, title, alt_text, mime_type, file_size, sort_order
) VALUES (
  @url, @thumbnail_url, @medium_url, @type, @title, @alt_text, @mime_type, @file_size, @sort_order
)
RETURNING *;

-- name: GetMediaByID :one
SELECT * FROM media WHERE id = @id;

-- name: ListMedia :many
SELECT * FROM media ORDER BY sort_order ASC;

-- name: UpdateMediaSortOrder :exec
UPDATE media SET sort_order = @sort_order WHERE id = @id;

-- name: DeleteMedia :exec
DELETE FROM media WHERE id = @id;


-- name: ListPublicMedia :many
SELECT * FROM media
WHERE is_public = true
ORDER BY sort_order ASC;


-- name: ListMediaPaginated :many
SELECT * FROM media
ORDER BY sort_order ASC
LIMIT $1 OFFSET $2;


-- name: CountMedia :one
SELECT COUNT(*) FROM media;


-- name: UpdateMedia :exec
UPDATE media
SET alt_text = @alt_text, title = @title
WHERE id = @id;


-- name: SetPublicStatus :exec
UPDATE media
set is_public = @is_public
WHERE id = @id;

