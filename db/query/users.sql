-- name: CreateUser :one
INSERT INTO users (
    cognito_sub, username, password
)
VALUES ($1, $2, $3)
RETURNING cognito_sub, username, password, created_at;

-- name: GetUserByID :one
SELECT cognito_sub, username, password, created_at
FROM users
WHERE cognito_sub = $1;

-- name: GetUserByUsername :one
SELECT cognito_sub, username, password, created_at
FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT cognito_sub, username, password, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $2
WHERE cognito_sub = $1
RETURNING cognito_sub, username, password, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE cognito_sub = $1;