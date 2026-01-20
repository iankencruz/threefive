-- sql/queries/sessions.sql

-- name: CreateSession :exec
INSERT INTO sessions (token, data, expiry)
VALUES ($1, $2, $3);

-- name: GetSession :one
SELECT data FROM sessions
WHERE token = $1 AND expiry > NOW();

-- name: UpdateSession :exec
UPDATE sessions
SET data = $2, expiry = $3
WHERE token = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expiry < NOW();

-- name: CommitSession :exec
INSERT INTO sessions (token, data, expiry)
VALUES ($1, $2, $3)
ON CONFLICT (token) 
DO UPDATE SET data = EXCLUDED.data, expiry = EXCLUDED.expiry;
