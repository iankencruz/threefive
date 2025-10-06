-- backend/sql/queries/pages.sql

-- name: CreatePage :one
INSERT INTO pages (
    title, 
    slug, 
    page_type, 
    status,
    featured_image_id,
    author_id
)
VALUES (@title, @slug, @page_type, @status, @featured_image_id, @author_id)
RETURNING *;

-- name: GetPageByID :one
SELECT * FROM pages
WHERE id = @id AND deleted_at IS NULL;

-- name: GetPageBySlug :one
SELECT * FROM pages
WHERE slug = @slug 
  AND status = 'published' 
  AND deleted_at IS NULL;

-- name: GetPageBySlugAdmin :one
SELECT * FROM pages
WHERE slug = @slug AND deleted_at IS NULL;

-- name: ListPages :many
SELECT * FROM pages
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::page_status IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('page_type')::page_type IS NULL OR page_type = sqlc.narg('page_type'))
  AND (sqlc.narg('author_id')::uuid IS NULL OR author_id = sqlc.narg('author_id'))
ORDER BY 
  CASE WHEN @sort_by::text = 'created_at_desc' THEN created_at END DESC,
  CASE WHEN @sort_by::text = 'created_at_asc' THEN created_at END ASC,
  CASE WHEN @sort_by::text = 'updated_at_desc' THEN updated_at END DESC,
  CASE WHEN @sort_by::text = 'updated_at_asc' THEN updated_at END ASC,
  CASE WHEN @sort_by::text = 'title_asc' THEN title END ASC,
  CASE WHEN @sort_by::text = 'title_desc' THEN title END DESC,
  created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: ListPublishedPages :many
SELECT * FROM pages
WHERE status = 'published' 
  AND deleted_at IS NULL
  AND (sqlc.narg('page_type')::page_type IS NULL OR page_type = sqlc.narg('page_type'))
ORDER BY published_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CountPages :one
SELECT COUNT(*) FROM pages
WHERE deleted_at IS NULL
  AND (sqlc.narg('status')::page_status IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('page_type')::page_type IS NULL OR page_type = sqlc.narg('page_type'))
  AND (sqlc.narg('author_id')::uuid IS NULL OR author_id = sqlc.narg('author_id'));

-- name: UpdatePage :one
UPDATE pages
SET 
    title = COALESCE(@title, title),
    slug = COALESCE(@slug, slug),
    status = COALESCE(@status, status),
    featured_image_id = @featured_image_id,
    published_at = CASE 
        WHEN @status = 'published' AND published_at IS NULL 
        THEN NOW() 
        ELSE published_at 
    END,
    updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING *;

-- name: UpdatePageStatus :one
UPDATE pages
SET 
    status = @status,
    published_at = CASE 
        WHEN @status = 'published' AND published_at IS NULL 
        THEN NOW() 
        ELSE published_at 
    END,
    updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeletePage :exec
UPDATE pages
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id;

-- name: HardDeletePage :exec
DELETE FROM pages WHERE id = @id;

-- name: CheckSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM pages 
    WHERE slug = @slug 
      AND deleted_at IS NULL
      AND (@exclude_id::uuid IS NULL OR id != @exclude_id)
) AS exists;
