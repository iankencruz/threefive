-- name: AddMediaToGallery :exec
INSERT INTO gallery_media (gallery_id, media_id, sort_order)
VALUES (@gallery_id, @media_id, @sort_order);

-- name: RemoveMediaFromGallery :exec
DELETE FROM gallery_media WHERE gallery_id = @gallery_id AND media_id = @media_id;

-- name: ListMediaForGallery :many
SELECT m.*
FROM gallery_media gm
JOIN media m ON m.id = gm.media_id
WHERE gm.project_id = @project_id
ORDER BY gm.sort_order ASC;

-- name: UpdateGalleryMediaSortOrder :exec
WITH sorted(media_id, sort_order) AS (
  SELECT unnest(@media_ids::uuid[]), generate_series(0, cardinality(@media_ids) - 1)
)
UPDATE gallery_media gm
SET sort_order = sorted.sort_order
FROM sorted
WHERE gm.gallery_id = @gallery_id
  AND gm.media_id = @sorted.media_id

