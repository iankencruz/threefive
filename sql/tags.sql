-- Tags CRUD Operations

-- name: CreateTag :one
INSERT INTO tags (
    id,
    name,
    slug,
    created_at,
    updated_at
) VALUES (
    @id,
    @name,
    @slug,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetTagByID :one
SELECT * FROM tags
WHERE id = @id
LIMIT 1;

-- name: GetTagBySlug :one
SELECT * FROM tags
WHERE slug = @slug
LIMIT 1;

-- name: GetTagIDBySlug :one
SELECT id FROM tags
WHERE slug = @slug
LIMIT 1;

-- name: ListTags :many
SELECT * FROM tags
ORDER BY name ASC
LIMIT @limit_val
OFFSET @offset_val;

-- name: ListAllTags :many
SELECT * FROM tags
ORDER BY name ASC;

-- name: UpdateTag :one
UPDATE tags SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = @id;

-- name: CountTags :one
SELECT COUNT(*) FROM tags;

-- name: CheckTagSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM tags 
    WHERE slug = @slug 
      AND id != @tag_id
);

-- name: SearchTags :many
SELECT * FROM tags
WHERE name ILIKE @search_term
ORDER BY name ASC
LIMIT @limit_val
OFFSET @offset_val;

-- Tag Usage Statistics

-- name: GetTagUsageCount :one
SELECT COUNT(*) FROM project_tags
WHERE tag_id = @tag_id;

-- name: GetMostUsedTags :many
SELECT 
    t.*,
    COUNT(pt.project_id) as usage_count
FROM tags t
LEFT JOIN project_tags pt ON t.id = pt.tag_id
GROUP BY t.id
ORDER BY usage_count DESC, t.name ASC
LIMIT @limit_val;

-- name: GetUnusedTags :many
SELECT t.* FROM tags t
LEFT JOIN project_tags pt ON t.id = pt.tag_id
WHERE pt.tag_id IS NULL
ORDER BY t.name ASC;

-- Batch Operations

-- name: GetTagsByIDs :many
SELECT * FROM tags
WHERE id = ANY(@tag_ids::uuid[])
ORDER BY name ASC;

-- name: FindOrCreateTag :one
INSERT INTO tags (id, name, slug, created_at, updated_at)
VALUES (@id, @name, @slug, NOW(), NOW())
ON CONFLICT (slug) DO UPDATE SET updated_at = NOW()
RETURNING *;
