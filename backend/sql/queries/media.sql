-- backend/sql/queries/media.sql

-- name: CreateMedia :one
INSERT INTO media (
    filename, 
    original_filename, 
    mime_type, 
    size_bytes, 
    width, 
    height,
    storage_type, 
    storage_path, 
    s3_bucket, 
    s3_key, 
    s3_region,
    url,
    original_url,
    large_url,
    medium_url,
    thumbnail_url,
    original_path,
    large_path,
    medium_path,
    thumbnail_path,
    uploaded_by
)
VALUES (
    @filename, 
    @original_filename, 
    @mime_type, 
    @size_bytes, 
    @width, 
    @height,
    @storage_type, 
    @storage_path,
    @s3_bucket, 
    @s3_key, 
    @s3_region,
    @url,
    @original_url,
    @large_url,
    @medium_url,
    @thumbnail_url,
    @original_path,
    @large_path,
    @medium_path,
    @thumbnail_path,
    @uploaded_by
)
RETURNING *;

-- name: GetMediaByID :one
SELECT * FROM media
WHERE id = @id AND deleted_at IS NULL;

-- name: ListMedia :many
SELECT * FROM media
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: ListMediaByUser :many
SELECT * FROM media
WHERE uploaded_by = @uploaded_by AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: SoftDeleteMedia :exec
UPDATE media
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id;

-- name: HardDeleteMedia :exec
DELETE FROM media WHERE id = @id;

-- name: GetMediaStats :one
SELECT 
    COUNT(*) as total_files,
    SUM(size_bytes) as total_size_bytes,
    COUNT(DISTINCT uploaded_by) as unique_uploaders
FROM media
WHERE deleted_at IS NULL;

-- Media Relations Queries

-- name: LinkMediaToEntity :one
INSERT INTO media_relations (media_id, entity_type, entity_id, sort_order)
VALUES (@media_id, @entity_type, @entity_id, @sort_order)
ON CONFLICT (media_id, entity_type, entity_id) 
DO UPDATE SET sort_order = @sort_order
RETURNING *;

-- name: UnlinkMediaFromEntity :exec
DELETE FROM media_relations
WHERE media_id = @media_id 
  AND entity_type = @entity_type 
  AND entity_id = @entity_id;

-- name: GetMediaForEntity :many
SELECT m.*
FROM media m
INNER JOIN media_relations mr ON m.id = mr.media_id
WHERE mr.entity_type = @entity_type 
  AND mr.entity_id = @entity_id
  AND m.deleted_at IS NULL
ORDER BY mr.sort_order, mr.created_at;

-- name: GetEntitiesForMedia :many
SELECT entity_type, entity_id, sort_order, created_at
FROM media_relations
WHERE media_id = @media_id
ORDER BY created_at DESC;
