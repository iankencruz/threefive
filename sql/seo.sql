-- name: GetSEO :one
SELECT * FROM seo
WHERE entity_type = @entity_type
  AND entity_id = @entity_id
LIMIT 1;

-- name: UpsertSEO :one
INSERT INTO seo (
    entity_type,
    entity_id,
    seo_title,
    seo_description,
    og_title,
    og_description,
    og_image_id,
    canonical_url,
    robots_index,
    robots_follow
) VALUES (
    @entity_type,
    @entity_id,
    @seo_title,
    @seo_description,
    @og_title,
    @og_description,
    @og_image_id,
    @canonical_url,
    @robots_index,
    @robots_follow
)
ON CONFLICT (entity_type, entity_id) DO UPDATE SET
    seo_title       = EXCLUDED.seo_title,
    seo_description = EXCLUDED.seo_description,
    og_title        = EXCLUDED.og_title,
    og_description  = EXCLUDED.og_description,
    og_image_id     = EXCLUDED.og_image_id,
    canonical_url   = EXCLUDED.canonical_url,
    robots_index    = EXCLUDED.robots_index,
    robots_follow   = EXCLUDED.robots_follow,
    updated_at      = NOW()
RETURNING *;

-- name: DeleteSEO :exec
DELETE FROM seo
WHERE entity_type = @entity_type
  AND entity_id = @entity_id;
