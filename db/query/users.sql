-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING user_id, username, email, password_hash, created_at;

-- name: GetUserByID :one
SELECT user_id, username, email, password_hash, created_at
FROM users
WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT user_id, username, email, password_hash, created_at
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT user_id, username, email, password_hash, created_at
FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT user_id, username, email, password_hash, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2
WHERE user_id = $1
RETURNING user_id, username, email, password_hash, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;
