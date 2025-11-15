-- backend/sql/queries/blogs.sql

-- ============================================
-- Blogs Queries
-- ============================================

-- name: CreateBlog :one
INSERT INTO blogs (
    title,
    slug,
    status,
    excerpt,
    reading_time,
    is_featured,
    featured_image_id,
    published_at
)
VALUES (
    @title,
    @slug,
    @status,
    @excerpt,
    @reading_time,
    @is_featured,
    @featured_image_id,
    @published_at
)
RETURNING *;

-- name: GetBlogByID :one
SELECT * FROM blogs
WHERE id = @id AND deleted_at IS NULL;

-- name: GetBlogBySlug :one
SELECT * FROM blogs
WHERE slug = @slug AND deleted_at IS NULL;

-- name: ListBlogs :many
SELECT * FROM blogs
WHERE deleted_at IS NULL
  AND (@status::text = '' OR status = @status::page_status)
  AND (@is_featured::text = '' OR 
       (@is_featured = 'true' AND is_featured = true) OR 
       (@is_featured = 'false' AND is_featured = false))
ORDER BY 
  CASE WHEN @sort_by = 'created_at' AND @sort_order = 'desc' THEN created_at END DESC,
  CASE WHEN @sort_by = 'created_at' AND @sort_order = 'asc' THEN created_at END ASC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'desc' THEN published_at END DESC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'asc' THEN published_at END ASC,
  CASE WHEN @sort_by = 'title' AND @sort_order = 'desc' THEN title END DESC,
  CASE WHEN @sort_by = 'title' AND @sort_order = 'asc' THEN title END ASC,
  created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CountBlogs :one
SELECT COUNT(*) FROM blogs
WHERE deleted_at IS NULL
  AND (@status::text = '' OR status = @status::page_status)
  AND (@is_featured::text = '' OR 
       (@is_featured = 'true' AND is_featured = true) OR 
       (@is_featured = 'false' AND is_featured = false));

-- name: UpdateBlog :one
UPDATE blogs
SET 
    title = COALESCE(@title, title),
    slug = COALESCE(@slug, slug),
    status = COALESCE(@status, status),
    excerpt = @excerpt,
    reading_time = @reading_time,
    is_featured = COALESCE(@is_featured, is_featured),
    featured_image_id = @featured_image_id,
    published_at = @published_at,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UpdateBlogStatus :one
UPDATE blogs
SET status = @status, updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: SoftDeleteBlog :exec
UPDATE blogs
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id;

-- name: CheckBlogSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM blogs 
    WHERE slug = @slug 
    AND deleted_at IS NULL
    AND (@exclude_id::uuid IS NULL OR id != @exclude_id)
);
