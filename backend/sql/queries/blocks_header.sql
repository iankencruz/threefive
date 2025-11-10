-- backend/sql/queries/blocks_header.sql

-- ============================================
-- Header Block Queries
-- ============================================

-- name: CreateHeaderBlock :one
INSERT INTO block_header (
    block_id, 
    heading, 
    subheading, 
    level
)
VALUES (@block_id, @heading, @subheading, @level)
RETURNING *;

-- name: GetHeaderBlockByBlockID :one
SELECT * FROM block_header
WHERE block_id = @block_id;

-- name: GetHeaderBlocksByEntity :many
SELECT bh.*
FROM block_header bh
INNER JOIN blocks b ON b.id = bh.block_id
WHERE b.entity_type = @entity_type AND b.entity_id = @entity_id
ORDER BY b.sort_order;

-- name: UpdateHeaderBlock :one
UPDATE block_header
SET 
    heading = COALESCE(@heading, heading),
    subheading = @subheading,
    level = COALESCE(@level, level)
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteHeaderBlock :exec
DELETE FROM block_header WHERE block_id = @block_id;
