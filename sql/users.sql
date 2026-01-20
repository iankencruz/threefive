-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @id AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = @email
LIMIT 1;



-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  password_hash
) VALUES (
  @first_name,
  @last_name,
  @email,
  @password_hash
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
  first_name = COALESCE(sqlc.narg('first_name'), first_name),
  last_name = COALESCE(sqlc.narg('last_name'), last_name),
  email = COALESCE(sqlc.narg('email'), email),
  password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
  updated_at = NOW()
WHERE id = @id
RETURNING *;


-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = @id;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CheckEmailExists :one
SELECT EXISTS(
  SELECT 1 FROM users
  WHERE email = @email
);

