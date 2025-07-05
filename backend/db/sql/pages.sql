
-- name: CreatePage :one
INSERT INTO pages (
  slug, title, banner_image_id, seo_title, seo_description, seo_canonical,
  content, is_draft, is_published
)
VALUES (
  @slug, @title, @banner_image_id, @seo_title, @seo_description, @seo_canonical,
  @content, @is_draft, @is_published
)
RETURNING *;

-- name: GetPageByID :one
SELECT * FROM pages WHERE id = @id;

-- name: GetPageBySlug :one
SELECT * FROM pages WHERE slug = @slug;

-- name: GetPublishedPageBySlug :one
SELECT * FROM pages WHERE slug = @slug AND is_published = true;


-- name: ListPages :many
SELECT * FROM pages ORDER BY updated_at DESC;

-- name: UpdatePage :one
UPDATE pages
SET title = @title,
    slug = @slug,
    banner_image_id = @banner_image_id,
    seo_title = @seo_title,
    seo_description = @seo_description,
    seo_canonical = @seo_canonical,
    content = @content,
    is_draft = @is_draft,
    is_published = @is_published,
    updated_at = now()
WHERE id = @id
RETURNING *;

-- name: DeletePage :exec
DELETE FROM pages WHERE id = @id;
