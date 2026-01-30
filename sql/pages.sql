-- Basic page queries

-- name: GetPageBySlug :one
SELECT * FROM pages
WHERE slug = @slug AND deleted_at IS NULL
LIMIT 1;

-- name: GetPageByType :one
SELECT * FROM pages
WHERE page_type = @page_type AND deleted_at IS NULL
LIMIT 1;

-- name: ListPages :many
SELECT * FROM pages
WHERE deleted_at IS NULL
ORDER BY page_type ASC;

-- name: UpdatePage :one
UPDATE pages SET
    title = COALESCE(sqlc.narg('title'), title),
    slug = COALESCE(sqlc.narg('slug'), slug),
    hero_media_id = COALESCE(sqlc.narg('hero_media_id'), hero_media_id),
    header = COALESCE(sqlc.narg('header'), header),
    sub_header = COALESCE(sqlc.narg('sub_header'), sub_header),
    content = COALESCE(sqlc.narg('content'), content),
    content_image_id = COALESCE(sqlc.narg('content_image_id'), content_image_id),
    cta_text = COALESCE(sqlc.narg('cta_text'), cta_text),
    cta_link = COALESCE(sqlc.narg('cta_link'), cta_link),
    email = COALESCE(sqlc.narg('email'), email),
    social_links = COALESCE(sqlc.narg('social_links'), social_links),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- Featured Projects for About Page

-- name: AddFeaturedProject :one
INSERT INTO page_featured_projects (page_id, project_id, display_order)
VALUES (@page_id, @project_id, @display_order)
ON CONFLICT (page_id, project_id) 
DO UPDATE SET display_order = EXCLUDED.display_order
RETURNING *;

-- name: GetFeaturedProjects :many
SELECT 
    p.id,
    p.title,
    p.slug,
    p.description,
    p.featured_image_id,
    pfp.display_order
FROM page_featured_projects pfp
JOIN projects p ON pfp.project_id = p.id
WHERE pfp.page_id = @page_id 
  AND p.deleted_at IS NULL
  AND p.status = 'published'
ORDER BY pfp.display_order ASC;

-- name: RemoveFeaturedProject :exec
DELETE FROM page_featured_projects
WHERE page_id = @page_id AND project_id = @project_id;

-- name: ClearFeaturedProjects :exec
DELETE FROM page_featured_projects
WHERE page_id = @page_id;

-- name: CountFeaturedProjects :one
SELECT COUNT(*) FROM page_featured_projects
WHERE page_id = @page_id;
