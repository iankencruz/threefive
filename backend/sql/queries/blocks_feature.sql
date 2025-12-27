-- backend/sql/queries/blocks_about.sql

-- ============================================
-- Feature Block Queries
-- ============================================

-- name: CreateFeatureBlock :one
INSERT INTO block_feature (
  block_id, 
  title,
  description,
  heading,
  subheading,
  image_id
)
VALUES (
  @block_id, 
  @title,
  @description,
  @heading,
  @subheading,
  @image_id
)
returning *;


-- name: GetFeatureBlockByBlockID :one
SELECT * FROM block_feature
WHERE block_id = @block_id;

-- name: GetFeatureBlocksByEntity :many
SELECT bf.*
FROM block_feature bf
INNER JOIN blocks b ON b.id = bf.block_id
WHERE b.entity_type = @entity_type AND b.entity_id = @entity_id
ORDER BY b.sort_order;

-- name: UpdateFeatureBlock :one
UPDATE block_feature
SET 
  title = COALESCE(@title, title),
  description = COALESCE(@description, description),
  heading =  COALESCE(@heading, heading),
  subheading = COALESCE(@subheading, subheading),
  image_id = COALESCE(@image_id, image_id)
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteFeatureBlock :exec
DELETE FROM block_feature WHERE block_id = @block_id;
