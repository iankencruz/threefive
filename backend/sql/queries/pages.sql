-- backend/sql/queries/pages.sql

-- ============================================
-- Pages Queries (Generic pages only)
-- ============================================

-- name: CreatePage :one
INSERT INTO pages (
    title,
    slug,
    status,
    featured_image_id,
    author_id,
    published_at
)
VALUES (
    @title,
    @slug,
    @status,
    @featured_image_id,
    @author_id,
    @published_at
)
RETURNING *;

-- name: GetPageByID :one
SELECT * FROM pages
WHERE id = @id AND deleted_at IS NULL;

-- name: GetPageBySlug :one
SELECT * FROM pages
WHERE slug = @slug AND deleted_at IS NULL;

-- name: ListPages :many
SELECT * FROM pages
WHERE deleted_at IS NULL
  AND (@status::page_status IS NULL OR status = @status)
  AND (@author_id::uuid IS NULL OR author_id = @author_id)
ORDER BY 
  CASE WHEN @sort_by = 'created_at_desc' THEN created_at END DESC,
  CASE WHEN @sort_by = 'created_at_asc' THEN created_at END ASC,
  CASE WHEN @sort_by = 'published_at_desc' THEN published_at END DESC,
  CASE WHEN @sort_by = 'published_at_asc' THEN published_at END ASC,
  created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CountPages :one
SELECT COUNT(*) FROM pages
WHERE deleted_at IS NULL
  AND (@status::page_status IS NULL OR status = @status)
  AND (@author_id::uuid IS NULL OR author_id = @author_id);

-- name: UpdatePage :one
UPDATE pages
SET 
    title = COALESCE(@title, title),
    slug = COALESCE(@slug, slug),
    status = COALESCE(@status, status),
    featured_image_id = @featured_image_id,
    published_at = @published_at,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UpdatePageStatus :one
UPDATE pages
SET status = @status, updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: SoftDeletePage :exec
UPDATE pages
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id;

-- name: CheckSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM pages 
    WHERE slug = @slug 
    AND deleted_at IS NULL
    AND (@exclude_id::uuid IS NULL OR id != @exclude_id)
);

-- name: PurgeOldDeletedPages :one
WITH deleted AS (
    DELETE FROM pages 
    WHERE deleted_at IS NOT NULL           -- Has been soft-deleted
      AND deleted_at < @cutoff_date        -- @cutoff_date is a PARAMETER (not a column)
    RETURNING id
)
SELECT COUNT(*) as count FROM deleted;
