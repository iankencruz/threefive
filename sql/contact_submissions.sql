-- name: CreateContactSubmission :one
INSERT INTO contact_submissions (first_name, last_name, email, subject, message)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListContactSubmissions :many
SELECT * FROM contact_submissions
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetContactSubmission :one
SELECT * FROM contact_submissions
WHERE id = $1 AND deleted_at IS NULL;

-- name: MarkContactSubmissionRead :one
UPDATE contact_submissions
SET read_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteContactSubmission :exec
UPDATE contact_submissions
SET deleted_at = NOW()
WHERE id = $1;
