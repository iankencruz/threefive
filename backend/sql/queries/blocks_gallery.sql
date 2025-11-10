-- backend/sql/queries/blocks_gallery.sql

-- ============================================
-- Gallery Block Queries
-- ============================================

-- name: CreateGalleryBlock :one
INSERT INTO block_gallery (
    block_id, 
    title
)
VALUES (@block_id, @title)
RETURNING *;

-- name: GetGalleryBlockByBlockID :one
SELECT * FROM block_gallery
WHERE block_id = @block_id;

-- name: GetGalleryBlocksByEntity :many
SELECT bg.*
FROM block_gallery bg
INNER JOIN blocks b ON b.id = bg.block_id
WHERE b.entity_type = @entity_type AND b.entity_id = @entity_id
ORDER BY b.sort_order;

-- name: UpdateGalleryBlock :one
UPDATE block_gallery
SET title = @title
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteGalleryBlock :exec
DELETE FROM block_gallery WHERE block_id = @block_id;
