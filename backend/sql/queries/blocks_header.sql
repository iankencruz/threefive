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

-- name: GetHeaderBlocksByPageID :many
SELECT bh.*
FROM block_header bh
INNER JOIN blocks b ON b.id = bh.block_id
WHERE b.page_id = @page_id
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
