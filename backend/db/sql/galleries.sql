-- name: CreateGallery:one
INSERT INTO galleries (
  title, slug, description, is_published, created_at, updated_at
) VALUES (
  @title, @slug, @description, @is_published, @created_at, @published_at
)
  RETURNING *;

-- name: GetGalleryByID :one
SELECT * FROM galleries WHERE id = @id;

-- name: GetGalleryBySlug :one
SELECT * FROM galleries WHERE slug = @slug;

-- name: ListGalleries :many
SELECT * FROM galleries ORDER BY created_at DESC;

-- name: UpdateGallery: one
UPDATE galleries
SET 
    title = @title,
    slug = @slug,
    description = @description,
    is_published = @is_published,
    created_at = @created_at,
    updated_at = @updated_at 
WHERE id = @id
RETURNING *;


-- name: UpdateGalleryBySlug: one
UPDATE galleries
SET 
    title = @title,
    slug = @slug,
    description = @description,
    is_published = @is_published,
    created_at = @created_at,
    updated_at = @updated_at 
WHERE slug = @id
RETURNING *;
