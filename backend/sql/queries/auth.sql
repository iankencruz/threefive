-- name: CreateSession :one
INSERT INTO sessions (user_id, token, expires_at, ip_address, user_agent)
VALUES (@user_id, @token, @expires_at, @ip_address, @user_agent)
RETURNING *;

-- name: GetSessionByToken :one
SELECT s.*, u.id as user_id, u.email, u.first_name, u.last_name, u.created_at as user_created_at
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token = @token AND s.is_active = true AND s.expires_at > NOW();

-- name: GetActiveSessionsByUserID :many
SELECT * FROM sessions
WHERE user_id = @user_id AND is_active = true AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: UpdateSessionExpiry :one
UPDATE sessions
SET expires_at = @expires_at, updated_at = NOW()
WHERE token = @token AND is_active = true
RETURNING *;

-- name: DeactivateSession :exec
UPDATE sessions
SET is_active = false, updated_at = NOW()
WHERE token = @token;

-- name: DeactivateAllUserSessions :exec
UPDATE sessions
SET is_active = false, updated_at = NOW()
WHERE user_id = @user_id AND is_active = true;

-- name: CleanupExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < NOW() OR (created_at < NOW() - INTERVAL '30 days');

-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES (@user_id, @token, @expires_at)
RETURNING *;

-- name: GetPasswordResetToken :one
SELECT prt.*, u.email, u.first_name, u.last_name
FROM password_reset_tokens prt
JOIN users u ON prt.user_id = u.id
WHERE prt.token = @token AND prt.is_used = false AND prt.expires_at > NOW();

-- name: UsePasswordResetToken :exec
UPDATE password_reset_tokens
SET is_used = true, used_at = NOW()
WHERE token = @token;

-- name: CleanupExpiredPasswordResetTokens :exec
DELETE FROM password_reset_tokens
WHERE expires_at < NOW() OR is_used = true;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash, updated_at = NOW()
WHERE id = @user_id;
