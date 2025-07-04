// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: project_media.sql

package generated

import (
	"context"

	"github.com/google/uuid"
)

const addMediaToProject = `-- name: AddMediaToProject :exec
INSERT INTO project_media (project_id, media_id, sort_order)
VALUES ($1, $2, $3)
`

type AddMediaToProjectParams struct {
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	MediaID   uuid.UUID `db:"media_id" json:"media_id"`
	SortOrder int32     `db:"sort_order" json:"sort_order"`
}

func (q *Queries) AddMediaToProject(ctx context.Context, arg AddMediaToProjectParams) error {
	_, err := q.db.Exec(ctx, addMediaToProject, arg.ProjectID, arg.MediaID, arg.SortOrder)
	return err
}

const listMediaForProject = `-- name: ListMediaForProject :many
SELECT m.id, m.url, m.thumbnail_url, m.type, m.is_public, m.title, m.alt_text, m.mime_type, m.file_size, m.sort_order, m.created_at, m.updated_at, m.medium_url
FROM project_media pm
JOIN media m ON m.id = pm.media_id
WHERE pm.project_id = $1
ORDER BY pm.sort_order ASC
`

func (q *Queries) ListMediaForProject(ctx context.Context, projectID uuid.UUID) ([]Media, error) {
	rows, err := q.db.Query(ctx, listMediaForProject, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Media
	for rows.Next() {
		var i Media
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.ThumbnailUrl,
			&i.Type,
			&i.IsPublic,
			&i.Title,
			&i.AltText,
			&i.MimeType,
			&i.FileSize,
			&i.SortOrder,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MediumUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeMediaFromProject = `-- name: RemoveMediaFromProject :exec
DELETE FROM project_media WHERE project_id = $1 AND media_id = $2
`

type RemoveMediaFromProjectParams struct {
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	MediaID   uuid.UUID `db:"media_id" json:"media_id"`
}

func (q *Queries) RemoveMediaFromProject(ctx context.Context, arg RemoveMediaFromProjectParams) error {
	_, err := q.db.Exec(ctx, removeMediaFromProject, arg.ProjectID, arg.MediaID)
	return err
}

const updateProjectMediaSortOrder = `-- name: UpdateProjectMediaSortOrder :exec
WITH sorted(media_id, sort_order) AS (
  SELECT unnest($2::uuid[]), generate_series(0, cardinality($2) - 1)
)
UPDATE project_media pm
SET sort_order = sorted.sort_order
FROM sorted
WHERE pm.project_id = $1
  AND pm.media_id = sorted.media_id
`

type UpdateProjectMediaSortOrderParams struct {
	ProjectID uuid.UUID   `db:"project_id" json:"project_id"`
	MediaIds  []uuid.UUID `db:"media_ids" json:"media_ids"`
}

func (q *Queries) UpdateProjectMediaSortOrder(ctx context.Context, arg UpdateProjectMediaSortOrderParams) error {
	_, err := q.db.Exec(ctx, updateProjectMediaSortOrder, arg.ProjectID, arg.MediaIds)
	return err
}
