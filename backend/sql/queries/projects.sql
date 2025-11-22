-- backend/sql/queries/projects.sql

-- ============================================
-- Projects Queries
-- ============================================

-- name: CreateProject :one
INSERT INTO projects (
    title,
    slug,
    description,
    project_date,
    status,
    client_name,
    project_year,
    project_url,
    technologies,
    project_status,
    featured_image_id,
    published_at
)
VALUES (
    @title,
    @slug,
    @description,
    @project_date,
    @status,
    @client_name,
    @project_year,
    @project_url,
    @technologies,
    @project_status,
    @featured_image_id,
    @published_at
)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = @id AND deleted_at IS NULL;

-- name: GetProjectBySlug :one
SELECT * FROM projects
WHERE slug = @slug AND deleted_at IS NULL;

-- name: ListProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
  AND (@status = '' OR status = @status::page_status)
ORDER BY 
  CASE WHEN @sort_by = 'created_at' AND @sort_order = 'desc' THEN created_at END DESC,
  CASE WHEN @sort_by = 'created_at' AND @sort_order = 'asc' THEN created_at END ASC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'desc' THEN published_at END DESC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'asc' THEN published_at END ASC,
  CASE WHEN @sort_by = 'title' AND @sort_order = 'desc' THEN title END DESC,
  CASE WHEN @sort_by = 'title' AND @sort_order = 'asc' THEN title END ASC,
  CASE WHEN @sort_by = 'project_date' AND @sort_order = 'desc' THEN project_date END DESC,
  CASE WHEN @sort_by = 'project_date' AND @sort_order = 'asc' THEN project_date END ASC,
  CASE WHEN @sort_by = 'project_year' AND @sort_order = 'desc' THEN project_year END DESC,
  CASE WHEN @sort_by = 'project_year' AND @sort_order = 'asc' THEN project_year END ASC,
  created_at DESC
LIMIT @limit_val OFFSET @offset_val;


-- name: ListPublishedProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
  AND status = 'published'
ORDER BY 
  CASE WHEN @sort_by = 'created_at' AND @sort_order = 'desc' THEN created_at END DESC,
  CASE WHEN @sort_by = 'created_at' AND @sort_order = 'asc' THEN created_at END ASC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'desc' THEN published_at END DESC,
  CASE WHEN @sort_by = 'published_at' AND @sort_order = 'asc' THEN published_at END ASC,
  CASE WHEN @sort_by = 'title' AND @sort_order = 'desc' THEN title END DESC,
  CASE WHEN @sort_by = 'title' AND @sort_order = 'asc' THEN title END ASC,
  CASE WHEN @sort_by = 'project_date' AND @sort_order = 'desc' THEN project_date END DESC,
  CASE WHEN @sort_by = 'project_date' AND @sort_order = 'asc' THEN project_date END ASC,
  CASE WHEN @sort_by = 'project_year' AND @sort_order = 'desc' THEN project_year END DESC,
  CASE WHEN @sort_by = 'project_year' AND @sort_order = 'asc' THEN project_year END ASC,
  created_at DESC
LIMIT @limit_val OFFSET @offset_val;



-- name: CountProjects :one
SELECT COUNT(*) FROM projects
WHERE deleted_at IS NULL
  AND (@status = '' OR status = @status::page_status);


-- name: CountPublishedProjects :one
SELECT COUNT(*) FROM projects
WHERE deleted_at IS NULL
  AND status = 'published';



-- name: UpdateProject :one
UPDATE projects
SET 
    title = COALESCE(@title, title),
    slug = COALESCE(@slug, slug),
    description = @description,
    project_date = @project_date,
    status = COALESCE(@status, status),
    client_name = @client_name,
    project_year = @project_year,
    project_url = @project_url,
    technologies = COALESCE(@technologies, technologies),
    project_status = COALESCE(@project_status, project_status),
    featured_image_id = @featured_image_id,
    published_at = @published_at,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UpdateProjectStatus :one
UPDATE projects
SET status = @status, updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: SoftDeleteProject :exec
UPDATE projects
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id;

-- name: CheckProjectSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM projects 
    WHERE slug = @slug 
    AND deleted_at IS NULL
    AND (@exclude_id::uuid IS NULL OR id != @exclude_id)
);
