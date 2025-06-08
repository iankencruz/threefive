
-- name: CreateMedia :one
INSERT INTO media (
  url, thumbnail_url, type, title, alt_text, mime_type, file_size, sort_order
) VALUES (
  @url, @thumbnail_url, @type, @title, @alt_text, @mime_type, @file_size, @sort_order
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
