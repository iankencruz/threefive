-- backend/sql/queries/blocks_richtext.sql

-- ============================================
-- Richtext Block Queries
-- ============================================

-- name: CreateRichtextBlock :one
INSERT INTO block_richtext (block_id, content)
VALUES (@block_id, @content)
RETURNING *;

-- name: GetRichtextBlockByBlockID :one
SELECT * FROM block_richtext
WHERE block_id = @block_id;

-- name: GetRichtextBlocksByEntity :many
SELECT br.*
FROM block_richtext br
INNER JOIN blocks b ON b.id = br.block_id
WHERE b.entity_type = @entity_type AND b.entity_id = @entity_id
ORDER BY b.sort_order;

-- name: UpdateRichtextBlock :one
UPDATE block_richtext
SET content = @content
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteRichtextBlock :exec
DELETE FROM block_richtext WHERE block_id = @block_id;
