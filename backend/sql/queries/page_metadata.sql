-- backend/sql/queries/page_metadata.sql

-- ============================================
-- Page SEO Queries
-- ============================================

-- name: CreatePageSEO :one
INSERT INTO page_seo (
    page_id,
    meta_title,
    meta_description,
    og_title,
    og_description,
    og_image_id,
    canonical_url,
    robots_index,
    robots_follow
)
VALUES (
    @page_id,
    @meta_title,
    @meta_description,
    @og_title,
    @og_description,
    @og_image_id,
    @canonical_url,
    @robots_index,
    @robots_follow
)
RETURNING *;

-- name: GetPageSEO :one
SELECT * FROM page_seo
WHERE page_id = @page_id;

-- name: UpdatePageSEO :one
UPDATE page_seo
SET 
    meta_title = @meta_title,
    meta_description = @meta_description,
    og_title = @og_title,
    og_description = @og_description,
    og_image_id = @og_image_id,
    canonical_url = @canonical_url,
    robots_index = @robots_index,
    robots_follow = @robots_follow,
    updated_at = NOW()
WHERE page_id = @page_id
RETURNING *;

-- name: UpsertPageSEO :one
INSERT INTO page_seo (
    page_id,
    meta_title,
    meta_description,
    og_title,
    og_description,
    og_image_id,
    canonical_url,
    robots_index,
    robots_follow
)
VALUES (
    @page_id,
    @meta_title,
    @meta_description,
    @og_title,
    @og_description,
    @og_image_id,
    @canonical_url,
    @robots_index,
    @robots_follow
)
ON CONFLICT (page_id) 
DO UPDATE SET
    meta_title = EXCLUDED.meta_title,
    meta_description = EXCLUDED.meta_description,
    og_title = EXCLUDED.og_title,
    og_description = EXCLUDED.og_description,
    og_image_id = EXCLUDED.og_image_id,
    canonical_url = EXCLUDED.canonical_url,
    robots_index = EXCLUDED.robots_index,
    robots_follow = EXCLUDED.robots_follow,
    updated_at = NOW()
RETURNING *;

-- name: DeletePageSEO :exec
DELETE FROM page_seo WHERE page_id = @page_id;

-- ============================================
-- Project Data Queries
-- ============================================

-- name: CreateProjectData :one
INSERT INTO page_project_data (
    page_id,
    client_name,
    project_year,
    project_url,
    technologies,
    project_status
)
VALUES (
    @page_id,
    @client_name,
    @project_year,
    @project_url,
    @technologies,
    @project_status
)
RETURNING *;

-- name: GetProjectData :one
SELECT * FROM page_project_data
WHERE page_id = @page_id;

-- name: UpdateProjectData :one
UPDATE page_project_data
SET 
    client_name = @client_name,
    project_year = @project_year,
    project_url = @project_url,
    technologies = @technologies,
    project_status = @project_status,
    updated_at = NOW()
WHERE page_id = @page_id
RETURNING *;

-- name: UpsertProjectData :one
INSERT INTO page_project_data (
    page_id,
    client_name,
    project_year,
    project_url,
    technologies,
    project_status
)
VALUES (
    @page_id,
    @client_name,
    @project_year,
    @project_url,
    @technologies,
    @project_status
)
ON CONFLICT (page_id) 
DO UPDATE SET
    client_name = EXCLUDED.client_name,
    project_year = EXCLUDED.project_year,
    project_url = EXCLUDED.project_url,
    technologies = EXCLUDED.technologies,
    project_status = EXCLUDED.project_status,
    updated_at = NOW()
RETURNING *;

-- name: DeleteProjectData :exec
DELETE FROM page_project_data WHERE page_id = @page_id;

-- name: ListProjectPages :many
SELECT p.*, ppd.*
FROM pages p
INNER JOIN page_project_data ppd ON ppd.page_id = p.id
WHERE p.status = 'published' AND p.deleted_at IS NULL
ORDER BY ppd.project_year DESC, p.published_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- ============================================
-- Blog Data Queries
-- ============================================

-- name: CreateBlogData :one
INSERT INTO page_blog_data (
    page_id,
    excerpt,
    reading_time,
    is_featured
)
VALUES (
    @page_id,
    @excerpt,
    @reading_time,
    @is_featured
)
RETURNING *;

-- name: GetBlogData :one
SELECT * FROM page_blog_data
WHERE page_id = @page_id;

-- name: UpdateBlogData :one
UPDATE page_blog_data
SET 
    excerpt = @excerpt,
    reading_time = @reading_time,
    is_featured = @is_featured,
    updated_at = NOW()
WHERE page_id = @page_id
RETURNING *;

-- name: UpsertBlogData :one
INSERT INTO page_blog_data (
    page_id,
    excerpt,
    reading_time,
    is_featured
)
VALUES (
    @page_id,
    @excerpt,
    @reading_time,
    @is_featured
)
ON CONFLICT (page_id) 
DO UPDATE SET
    excerpt = EXCLUDED.excerpt,
    reading_time = EXCLUDED.reading_time,
    is_featured = EXCLUDED.is_featured,
    updated_at = NOW()
RETURNING *;

-- name: DeleteBlogData :exec
DELETE FROM page_blog_data WHERE page_id = @page_id;

-- name: ListBlogPages :many
SELECT p.*, pbd.*
FROM pages p
INNER JOIN page_blog_data pbd ON pbd.page_id = p.id
WHERE p.status = 'published' AND p.deleted_at IS NULL
ORDER BY p.published_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: ListFeaturedBlogPages :many
SELECT p.*, pbd.*
FROM pages p
INNER JOIN page_blog_data pbd ON pbd.page_id = p.id
WHERE p.status = 'published' 
  AND p.deleted_at IS NULL
  AND pbd.is_featured = true
ORDER BY p.published_at DESC
LIMIT @limit_val;
