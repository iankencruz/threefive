-- Projects CRUD Operations

-- name: CreateProject :one
INSERT INTO projects (
    id,
    title,
    slug,
    description,
    project_date,
    status,
    client_name,
    project_year,
    project_url,
    project_status,
    featured_image_id,
    author_id,
    created_at,
    updated_at
) VALUES (
    @id,
    @title,
    @slug,
    @description,
    @project_date,
    @status,
    @client_name,
    @project_year,
    @project_url,
    @project_status,
    @featured_image_id,
    @author_id,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = @id AND deleted_at IS NULL
LIMIT 1;

-- name: GetProjectBySlug :one
SELECT * FROM projects
WHERE slug = @slug AND deleted_at IS NULL
LIMIT 1;

-- name: GetProjectIDBySlug :one
SELECT id FROM projects
WHERE slug = @slug AND deleted_at IS NULL
LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: ListPublishedProjects :many
SELECT * FROM projects
WHERE status = 'published' 
  AND deleted_at IS NULL
  AND published_at IS NOT NULL
ORDER BY published_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: UpdateProject :one
UPDATE projects SET
    title = COALESCE(sqlc.narg('title'), title),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    project_date = COALESCE(sqlc.narg('project_date'), project_date),
    status = COALESCE(sqlc.narg('status'), status),
    client_name = COALESCE(sqlc.narg('client_name'), client_name),
    project_year = COALESCE(sqlc.narg('project_year'), project_year),
    project_url = COALESCE(sqlc.narg('project_url'), project_url),
    project_status = COALESCE(sqlc.narg('project_status'), project_status),
    featured_image_id = COALESCE(sqlc.narg('featured_image_id'), featured_image_id),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: PublishProject :one
UPDATE projects SET
    status = 'published',
    published_at = NOW(),
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UnpublishProject :one
UPDATE projects SET
    status = 'draft',
    published_at = NULL,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: SoftDeleteProject :exec
UPDATE projects
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = @id;

-- name: RestoreProject :exec
UPDATE projects
SET
    deleted_at = NULL,
    updated_at = NOW()
WHERE id = @id;

-- name: HardDeleteProject :exec
DELETE FROM projects
WHERE id = @id;

-- name: CountProjects :one
SELECT COUNT(*) FROM projects
WHERE deleted_at IS NULL;

-- name: CountPublishedProjects :one
SELECT COUNT(*) FROM projects
WHERE status = 'published' 
  AND deleted_at IS NULL
  AND published_at IS NOT NULL;

-- name: CheckProjectSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM projects 
    WHERE slug = @slug 
      AND id != @project_id
      AND deleted_at IS NULL
);

-- Project Tags Operations

-- name: AddProjectTag :one
INSERT INTO project_tags (project_id, tag_id, created_at)
VALUES (@project_id, @tag_id, NOW())
ON CONFLICT (project_id, tag_id) DO NOTHING
RETURNING *;

-- name: RemoveProjectTag :exec
DELETE FROM project_tags
WHERE project_id = @project_id AND tag_id = @tag_id;

-- name: GetProjectTags :many
SELECT t.* FROM tags t
JOIN project_tags pt ON t.id = pt.tag_id
WHERE pt.project_id = @project_id
ORDER BY t.name ASC;

-- name: ClearProjectTags :exec
DELETE FROM project_tags
WHERE project_id = @project_id;

-- name: CountProjectTags :one
SELECT COUNT(*) FROM project_tags
WHERE project_id = @project_id;

-- Project Gallery Operations (via media_relations)

-- name: GetProjectGallery :many
SELECT m.* FROM media m
JOIN media_relations mr ON m.id = mr.media_id
WHERE mr.entity_type = 'project'
  AND mr.entity_id = @project_id
  AND mr.relation_type = 'gallery'
  AND m.deleted_at IS NULL
ORDER BY mr.sort_order ASC, m.created_at DESC;

-- name: GetProjectFeaturedImage :one
SELECT m.* FROM media m
WHERE m.id = @featured_image_id
  AND m.deleted_at IS NULL
LIMIT 1;

-- name: CountProjectGalleryImages :one
SELECT COUNT(*) FROM media_relations
WHERE entity_type = 'project'
  AND entity_id = @project_id
  AND relation_type = 'gallery';

-- Search and Filter Operations

-- name: SearchProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
  AND (
    title ILIKE @search_term
    OR description ILIKE @search_term
    OR client_name ILIKE @search_term
  )
ORDER BY created_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: GetProjectsByStatus :many
SELECT * FROM projects
WHERE status = @status
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: GetProjectsByTag :many
SELECT p.* FROM projects p
JOIN project_tags pt ON p.id = pt.project_id
WHERE pt.tag_id = @tag_id
  AND p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: GetProjectsByYear :many
SELECT * FROM projects
WHERE project_year = @project_year
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetRecentProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val;
