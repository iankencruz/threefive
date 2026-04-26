

-- name: GetConfigByCode :one
SELECT * FROM system_config
WHERE config_code = @config_code
LIMIT 1;

-- name: ListConfig :many
SELECT * FROM system_config
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: CreateSystemConfig :one
INSERT INTO system_config (
  value
) VALUES (
  @value
)
RETURNING *;

-- name: UpdateConfig :one
UPDATE system_config SET 
  value = COALESCE(sqlc.narg('value'), value),
  updated_at = NOW()
WHERE config_code = @config_code
RETURNING *;


-- name: DeleteConfig :exec
DELETE FROM system_config
WHERE config_code = @config_code;

