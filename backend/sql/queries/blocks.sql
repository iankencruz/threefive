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

-- name: GetRichtextBlocksByPageID :many
SELECT br.*
FROM block_richtext br
INNER JOIN blocks b ON b.id = br.block_id
WHERE b.page_id = @page_id
ORDER BY b.sort_order;

-- name: UpdateRichtextBlock :one
UPDATE block_richtext
SET content = @content
WHERE block_id = @block_id
RETURNING *;

-- name: DeleteRichtextBlock :exec
DELETE FROM block_richtext WHERE block_id = @block_id;

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
