-- backend/sql/queries/blocks.sql

-- ============================================
-- Base Blocks Queries (Polymorphic)
-- ============================================

-- name: CreateBlock :one
INSERT INTO blocks (entity_type, entity_id, type, sort_order)
VALUES (@entity_type, @entity_id, @type, @sort_order)
RETURNING *;

-- name: GetBlocksByEntity :many
SELECT * FROM blocks
WHERE entity_type = @entity_type AND entity_id = @entity_id
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

-- name: DeleteBlocksByEntity :exec
DELETE FROM blocks 
WHERE entity_type = @entity_type AND entity_id = @entity_id;
