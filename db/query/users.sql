-- name: CreateUser :one
INSERT INTO users (
    username, password
)
VALUES ($1, $2)
RETURNING user_id, username, password, created_at;

-- name: GetUserByID :one
SELECT user_id, username, password, created_at
FROM users
WHERE user_id = $1;

-- name: GetUserByUsername :one
SELECT user_id, username, password, created_at
FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT user_id, username, password, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $2
WHERE user_id = $1
RETURNING user_id, username, password, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;