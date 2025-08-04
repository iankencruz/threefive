
-- name: CreateBlog :one
INSERT INTO blogs (
  slug, title, cover_image_id, seo_description, seo_title, canonical_url,
  is_draft, is_published, published_at
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetBlogBySlug :one
SELECT * FROM blogs
WHERE slug = $1;

-- name: ListBlogs :many
SELECT * FROM blogs
ORDER BY created_at DESC;

-- name: UpdateBlog :one
UPDATE blogs
SET
  title = $2,
  cover_image_id = $3,
  seo_description = $4,
  seo_title = $5,
  canonical_url = $6,
  is_draft = $7,
  is_published = $8,
  published_at = $9,
  updated_at = now()
WHERE slug = $1
RETURNING *;

-- name: DeleteBlog :exec
DELETE FROM blogs
WHERE slug = $1;




-- name: GetBlocksByOwner :many
SELECT id, type, sort_order 
FROM blocks
WHERE parent_id = $1
ORDER BY sort_order ASC;
