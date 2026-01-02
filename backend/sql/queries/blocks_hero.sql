-- backend/sql/queries/blocks_hero.sql

-- ============================================
-- Hero Block Queries
-- ============================================

-- name: CreateHeroBlock :one
INSERT INTO block_hero (
    block_id, 
    title, 
    subtitle, 
    image_id, 
    cta_text, 
    cta_url,
    navigation_bar
)
VALUES (@block_id, @title, @subtitle, @image_id, @cta_text, @cta_url, @navigation_bar)
RETURNING *;

-- name: GetHeroBlockByBlockID :one
SELECT * FROM block_hero
WHERE block_id = @block_id;

-- name: GetHeroBlocksByEntity :many
SELECT bh.*
FROM block_hero bh
INNER JOIN blocks b ON b.id = bh.block_id
WHERE b.entity_type = @entity_type AND b.entity_id = @entity_id
ORDER BY b.sort_order;

-- name: UpdateHeroBlock :one
UPDATE block_hero
SET 
    title = COALESCE(@title, title),
    subtitle = @subtitle,
    image_id = @image_id,
    cta_text = @cta_text,
    cta_url = @cta_url,
    navigation_bar = @navigation_bar
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteHeroBlock :exec
DELETE FROM block_hero WHERE block_id = @block_id;
