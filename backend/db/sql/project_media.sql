
-- name: AddMediaToProject :exec
INSERT INTO project_media (project_id, media_id, sort_order)
VALUES (@project_id, @media_id, @sort_order);

-- name: RemoveMediaFromProject :exec
DELETE FROM project_media WHERE project_id = @project_id AND media_id = @media_id;

-- name: ListMediaForProject :many
SELECT m.*
FROM project_media pm
JOIN media m ON m.id = pm.media_id
WHERE pm.project_id = @project_id
ORDER BY pm.sort_order ASC;




-- name: UpdateProjectMediaSortOrder :exec
WITH sorted(media_id, sort_order) AS (
  SELECT unnest(@media_ids::uuid[]), generate_series(0, cardinality(@media_ids) - 1)
)
UPDATE project_media pm
SET sort_order = sorted.sort_order
FROM sorted
WHERE pm.project_id = @project_id
  AND pm.media_id = sorted.media_id;

