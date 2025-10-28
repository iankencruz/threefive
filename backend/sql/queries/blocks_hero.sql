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
    cta_url
)
VALUES (@block_id, @title, @subtitle, @image_id, @cta_text, @cta_url)
RETURNING *;

-- name: GetHeroBlockByBlockID :one
SELECT * FROM block_hero
WHERE block_id = @block_id;

-- name: GetHeroBlocksByPageID :many
SELECT bh.*
FROM block_hero bh
INNER JOIN blocks b ON b.id = bh.block_id
WHERE b.page_id = @page_id
ORDER BY b.sort_order;

-- name: UpdateHeroBlock :one
UPDATE block_hero
SET 
    title = COALESCE(@title, title),
    subtitle = @subtitle,
    image_id = @image_id,
    cta_text = @cta_text,
    cta_url = @cta_url
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteHeroBlock :exec
DELETE FROM block_hero WHERE block_id = @block_id;


