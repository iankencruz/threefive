
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
UPDATE project_media
SET sort_order = @sort_order
WHERE project_id = @project_id AND media_id = @media_id;
