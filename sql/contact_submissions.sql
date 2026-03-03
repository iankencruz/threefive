-- name: CreateContactSubmission :one
INSERT INTO contact_submissions (first_name, last_name, email, subject, message)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: MarkEmailSent :exec
UPDATE contact_submissions
SET email_sent = TRUE,
    email_attempts = email_attempts + 1,
    email_last_attempted_at = NOW(),
    email_error = NULL
WHERE id = $1;

-- name: MarkEmailFailed :exec
UPDATE contact_submissions
SET email_attempts = email_attempts + 1,
    email_last_attempted_at = NOW(),
    email_error = $2
WHERE id = $1;

-- name: GetUnsentSubmissions :many
SELECT * FROM contact_submissions
WHERE email_sent = FALSE
  AND email_attempts < $1
  AND deleted_at IS NULL
ORDER BY created_at ASC;

-- name: ListContactSubmissions :many
SELECT * FROM contact_submissions
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: MarkContactSubmissionRead :exec
UPDATE contact_submissions
SET read_at = NOW()
WHERE id = $1;

-- name: DeleteContactSubmission :exec
UPDATE contact_submissions
SET deleted_at = NOW()
WHERE id = $1;
