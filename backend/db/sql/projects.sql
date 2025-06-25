
-- name: CreateProject :one
INSERT INTO projects (
    title, slug, description, meta_description, canonical_url, cover_media_id, is_published, published_at
) VALUES (
    @title, @slug, @description, @meta_description, @canonical_url, @cover_media_id, @is_published, @published_at
)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = @id;

-- name: GetProjectBySlug :one
SELECT * FROM projects WHERE slug = @slug;

-- name: ListProjects :many
SELECT * FROM projects ORDER BY created_at DESC;

-- name: ListPublishedProjects :many
SELECT * FROM projects
WHERE is_published = true
ORDER BY created_at DESC;


-- name: UpdateProject :one
UPDATE projects
SET
    title = @title,
    slug = @slug,
    description = @description,
    meta_description = @meta_description,
    canonical_url = @canonical_url,
    cover_media_id = @cover_media_id,
    is_published = @is_published,
    published_at = @published_at,
    updated_at = @updated_at
WHERE id = @id
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = @id;
