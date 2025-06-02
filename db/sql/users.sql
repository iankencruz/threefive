
-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password_hash)
VALUES (@first_name, @last_name, @email, @password_hash)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = @email;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @id;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;


