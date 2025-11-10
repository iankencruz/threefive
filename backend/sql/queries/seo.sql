-- backend/sql/queries/seo.sql

-- ============================================
-- SEO Queries (Polymorphic)
-- ============================================

-- name: CreateSEO :one
INSERT INTO seo (
    entity_type,
    entity_id,
    meta_title,
    meta_description,
    og_title,
    og_description,
    og_image_id,
    canonical_url,
    robots_index,
    robots_follow
)
VALUES (
    @entity_type,
    @entity_id,
    @meta_title,
    @meta_description,
    @og_title,
    @og_description,
    @og_image_id,
    @canonical_url,
    @robots_index,
    @robots_follow
)
RETURNING *;

-- name: GetSEO :one
SELECT * FROM seo
WHERE entity_type = @entity_type AND entity_id = @entity_id;

-- name: UpdateSEO :one
UPDATE seo
SET 
    meta_title = @meta_title,
    meta_description = @meta_description,
    og_title = @og_title,
    og_description = @og_description,
    og_image_id = @og_image_id,
    canonical_url = @canonical_url,
    robots_index = @robots_index,
    robots_follow = @robots_follow,
    updated_at = NOW()
WHERE entity_type = @entity_type AND entity_id = @entity_id
RETURNING *;

-- name: UpsertSEO :one
INSERT INTO seo (
    entity_type,
    entity_id,
    meta_title,
    meta_description,
    og_title,
    og_description,
    og_image_id,
    canonical_url,
    robots_index,
    robots_follow
)
VALUES (
    @entity_type,
    @entity_id,
    @meta_title,
    @meta_description,
    @og_title,
    @og_description,
    @og_image_id,
    @canonical_url,
    @robots_index,
    @robots_follow
)
ON CONFLICT (entity_type, entity_id) 
DO UPDATE SET
    meta_title = EXCLUDED.meta_title,
    meta_description = EXCLUDED.meta_description,
    og_title = EXCLUDED.og_title,
    og_description = EXCLUDED.og_description,
    og_image_id = EXCLUDED.og_image_id,
    canonical_url = EXCLUDED.canonical_url,
    robots_index = EXCLUDED.robots_index,
    robots_follow = EXCLUDED.robots_follow,
    updated_at = NOW()
RETURNING *;

-- name: DeleteSEO :exec
DELETE FROM seo 
WHERE entity_type = @entity_type AND entity_id = @entity_id;
