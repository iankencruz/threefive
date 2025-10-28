-- backend/sql/queries/blocks.sql

-- ============================================
-- Base Blocks Queries
-- ============================================

-- name: CreateBlock :one
INSERT INTO blocks (page_id, type, sort_order)
VALUES (@page_id, @type, @sort_order)
RETURNING *;

-- name: GetBlocksByPageID :many
SELECT * FROM blocks
WHERE page_id = @page_id
ORDER BY sort_order;

-- name: GetBlockByID :one
SELECT * FROM blocks
WHERE id = @id;

-- name: UpdateBlockOrder :exec
UPDATE blocks
SET sort_order = @sort_order, updated_at = NOW()
WHERE id = @id;

-- name: DeleteBlock :exec
DELETE FROM blocks WHERE id = @id;

-- name: DeleteBlocksByPageID :exec
DELETE FROM blocks WHERE page_id = @page_id;


