-- name: CreateContact :one
INSERT INTO contacts (
    name,
    email,
    subject,
    message,
    ip_address,
    user_agent
) VALUES (
    @name,
    @email,
    @subject,
    @message,
    @ip_address,
    @user_agent
) RETURNING *;

-- name: GetContactByID :one
SELECT * FROM contacts
WHERE id = @id AND deleted_at IS NULL;

-- name: ListContacts :many
SELECT * FROM contacts
WHERE deleted_at IS NULL
ORDER BY
    CASE WHEN @order_by = 'created_at_desc' THEN created_at END DESC,
    CASE WHEN @order_by = 'created_at_asc' THEN created_at END ASC
LIMIT @limit_count
OFFSET @offset_count;

-- name: ListContactsByStatus :many
SELECT * FROM contacts
WHERE status = @status AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_count
OFFSET @offset_count;

-- name: CountContacts :one
SELECT COUNT(*) FROM contacts
WHERE deleted_at IS NULL;

-- name: CountContactsByStatus :one
SELECT COUNT(*) FROM contacts
WHERE status = @status AND deleted_at IS NULL;

-- name: UpdateContactStatus :one
UPDATE contacts
SET status = @status, updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING *;

-- name: MarkEmailSent :exec
UPDATE contacts
SET email_sent = true,
    email_sent_at = NOW(),
    email_error = NULL,
    updated_at = NOW()
WHERE id = @id;

-- name: MarkEmailFailed :exec
UPDATE contacts
SET email_sent = false,
    email_error = @error,
    updated_at = NOW()
WHERE id = @id;

-- name: GetUnsentEmails :many
SELECT * FROM contacts
WHERE email_sent = false
  AND deleted_at IS NULL
ORDER BY created_at ASC
LIMIT @limit_count;

-- name: SoftDeleteContact :exec
UPDATE contacts
SET deleted_at = NOW()
WHERE id = @id;

-- name: HardDeleteContact :exec
DELETE FROM contacts
WHERE id = @id;

-- name: GetContactsByIPAddress :many
SELECT * FROM contacts
WHERE ip_address = @ip_address
  AND created_at > @after_time
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CleanupDeletedContacts :exec
DELETE FROM contacts
WHERE deleted_at IS NOT NULL
  AND deleted_at < NOW() - INTERVAL '30 days';
