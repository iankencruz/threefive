-- backend/sql/queries/blocks_about.sql

-- ============================================
-- About Block Queries
-- ============================================

-- name: CreateAboutBlock :one
INSERT INTO block_about (
  block_id, 
  title,
  description,
  heading,
  subheading
)
VALUES (
  @block_id, 
  @title,
  @description,
  @heading,
  @subheading
)
returning *;


-- name: GetAboutBlockByBlockID :one
SELECT * FROM block_about
WHERE block_id = @block_id;

-- name: GetAboutBlocksByEntity :many
SELECT ba.*
FROM block_about ba
INNER JOIN blocks b ON b.id = ba.block_id
WHERE b.entity_type = @entity_type AND b.entity_id = @entity_id
ORDER BY b.sort_order;

-- name: UpdateAboutBlock :one
UPDATE block_about
SET 
  title = COALESCE(@title, title),
  description = COALESCE(@description, description),
  heading =  COALESCE(@heading, heading),
  subheading = COALESCE(@subheading, subheading)
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteAboutBlock :exec
DELETE FROM block_about WHERE block_id = @block_id;
