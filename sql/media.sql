-- sql/queries/media.sql

-- name: CreateMedia :one
INSERT INTO media (
    id,
    filename,
    original_filename,
    mime_type,
    file_size,
    width,
    height,
    duration,
    storage_type,
    s3_bucket,
    s3_region,
    original_key,
    large_key,
    medium_key,
    thumbnail_key,
    alt_text,
    uploaded_by,
    created_at,
    updated_at
) VALUES (
    @id,
    @filename,
    @original_filename,
    @mime_type,
    @file_size,
    @width,
    @height,
    @duration,
    @storage_type,
    @s3_bucket,
    @s3_region,
    @original_key,
    @large_key,
    @medium_key,
    @thumbnail_key,
    @alt_text,
    @uploaded_by,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetMediaByID :one
SELECT * FROM media
WHERE id = @id
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetMediaByFilename :one
SELECT * FROM media
WHERE filename = @filename
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListMedia :many
SELECT * FROM media
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: ListMediaByType :many
SELECT * FROM media
WHERE mime_type LIKE @mime_type_pattern
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: ListMediaByUploader :many
SELECT * FROM media
WHERE uploaded_by = @uploaded_by
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: UpdateMedia :one
UPDATE media
SET
    alt_text = COALESCE(@alt_text, alt_text),
    updated_at = NOW()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: UpdateMediaAltText :one
UPDATE media
SET
    alt_text = @alt_text,
    updated_at = NOW()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteMedia :exec
UPDATE media
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = @id;

-- name: RestoreMedia :exec
UPDATE media
SET
    deleted_at = NULL,
    updated_at = NOW()
WHERE id = @id;

-- name: HardDeleteMedia :exec
DELETE FROM media
WHERE id = @id;

-- name: GetDeletedMedia :many
SELECT * FROM media
WHERE deleted_at IS NOT NULL
ORDER BY deleted_at DESC
LIMIT @limit_val
OFFSET @offset_val;

-- name: PurgeOldDeletedMedia :exec
DELETE FROM media
WHERE deleted_at IS NOT NULL
  AND deleted_at < NOW() - INTERVAL '30 days';

-- name: CountMedia :one
SELECT COUNT(*) FROM media
WHERE deleted_at IS NULL;

-- name: CountMediaByType :one
SELECT COUNT(*) FROM media
WHERE mime_type LIKE @mime_type_pattern
  AND deleted_at IS NULL;

-- Media Relations Queries

-- name: CreateMediaRelation :one
INSERT INTO media_relations (
    id,
    media_id,
    entity_type,
    entity_id,
    relation_type,
    sort_order
) VALUES (
    @id,
    @media_id,
    @entity_type,
    @entity_id,
    @relation_type,
    @sort_order
)
RETURNING *;

-- name: GetMediaForEntity :many
SELECT m.* FROM media m
JOIN media_relations mr ON m.id = mr.media_id
WHERE mr.entity_type = @entity_type
  AND mr.entity_id = @entity_id
  AND m.deleted_at IS NULL
ORDER BY mr.sort_order ASC, m.created_at DESC;

-- name: GetFeaturedMediaForEntity :one
SELECT m.* FROM media m
JOIN media_relations mr ON m.id = mr.media_id
WHERE mr.entity_type = @entity_type
  AND mr.entity_id = @entity_id
  AND mr.relation_type = 'featured'
  AND m.deleted_at IS NULL
LIMIT 1;

-- name: GetGalleryMediaForEntity :many
SELECT m.* FROM media m
JOIN media_relations mr ON m.id = mr.media_id
WHERE mr.entity_type = @entity_type
  AND mr.entity_id = @entity_id
  AND mr.relation_type = 'gallery'
  AND m.deleted_at IS NULL
ORDER BY mr.sort_order ASC, m.created_at DESC;

-- name: DeleteMediaRelation :exec
DELETE FROM media_relations
WHERE media_id = @media_id
  AND entity_type = @entity_type
  AND entity_id = @entity_id;

-- name: DeleteAllMediaRelationsForEntity :exec
DELETE FROM media_relations
WHERE entity_type = @entity_type
  AND entity_id = @entity_id;

-- name: ReorderGalleryMedia :exec
UPDATE media_relations
SET sort_order = @sort_order
WHERE media_id = @media_id
  AND entity_type = @entity_type
  AND entity_id = @entity_id
  AND relation_type = 'gallery';

-- Media Statistics

-- name: GetMediaStats :one
SELECT
    COUNT(*) as total_count,
    COALESCE(SUM(file_size), 0) as total_size,
    COUNT(*) FILTER (WHERE mime_type LIKE 'image/%') as image_count,
    COUNT(*) FILTER (WHERE mime_type LIKE 'video/%') as video_count,
    COUNT(*) FILTER (WHERE mime_type LIKE 'audio/%') as audio_count
FROM media
WHERE deleted_at IS NULL;

-- name: GetOrphanedMedia :many
SELECT m.* FROM media m
LEFT JOIN media_relations mr ON m.id = mr.media_id
WHERE mr.id IS NULL
  AND m.deleted_at IS NULL
  AND m.created_at < NOW() - INTERVAL '1 day'
ORDER BY m.created_at DESC;

-- name: GetMediaUsageByEntity :many
SELECT
    mr.entity_type,
    mr.entity_id,
    COUNT(*) as media_count
FROM media_relations mr
JOIN media m ON mr.media_id = m.id
WHERE m.deleted_at IS NULL
GROUP BY mr.entity_type, mr.entity_id
ORDER BY media_count DESC
LIMIT @limit_val;

-- Batch operations

-- name: BatchCreateMediaRelations :copyfrom
INSERT INTO media_relations (
    id,
    media_id,
    entity_type,
    entity_id,
    relation_type,
    sort_order
) VALUES (
    $1, $2, $3, $4, $5, $6
);
