
-- name: CreatePage :one
INSERT INTO pages (
  slug, title, cover_image_id, seo_title, seo_description, seo_canonical, content,
  is_draft, is_published
)
VALUES (
  @slug, @title, @cover_image_id, @seo_title, @seo_description, @seo_canonical, @content,
  @is_draft, @is_published
)
RETURNING *;

-- name: GetPageByID :one
SELECT * FROM pages WHERE id = @id;

-- name: GetPageBySlug :one
SELECT * FROM pages WHERE slug = @slug;

-- name: GetPublishedPageBySlug :one
SELECT * FROM pages WHERE slug = @slug AND is_published = true;



  

-- name: UpdatePage :one
UPDATE pages
SET title = @title,
    slug = @slug,
    cover_image_id = @cover_image_id,
    seo_title = @seo_title,
    seo_description = @seo_description,
    seo_canonical = @seo_canonical,
    is_draft = @is_draft,
    is_published = @is_published,
    content = @content,
    updated_at = now()
WHERE id = @id
RETURNING *;

-- name: DeletePage :exec
DELETE FROM pages WHERE id = @id;




-- Table Sorting

-- name: ListPagesByTitleAsc :many
SELECT * FROM pages ORDER BY title ASC;

-- name: ListPagesByTitleDesc :many
SELECT * FROM pages ORDER BY title DESC;

-- name: ListPagesByCreatedAsc :many
SELECT * FROM pages ORDER BY created_at ASC;

-- name: ListPagesByCreatedDesc :many
SELECT * FROM pages ORDER BY created_at DESC;

-- name: ListPagesByUpdatedAsc :many
SELECT * FROM pages ORDER BY updated_at ASC;

-- name: ListPagesByUpdatedDesc :many
SELECT * FROM pages ORDER BY updated_at DESC;


-- name: ListPagesByStatusAsc :many
-- Drafts first, then Published
SELECT * FROM pages
ORDER BY is_published ASC, updated_at DESC;

-- name: ListPagesByStatusDesc :many
-- Published first, then Drafts
SELECT * FROM pages
ORDER BY is_published DESC, updated_at DESC;




-- name: GetPublicPageWithGalleries :many
SELECT g.*
FROM galleries g
JOIN gallery_page gp ON gp.gallery_id = g.id
JOIN pages p ON gp.page_id = p.id
WHERE p.slug = @slug
ORDER BY gp.sort_order ASC;

-- name: GetMediaForGallery :many
SELECT m.*
FROM gallery_media gm
JOIN media m ON m.id = gm.media_id
WHERE gm.gallery_id = @gallery_id
ORDER BY gm.sort_order ASC;
