-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  password_hash,
  first_name,
  last_name
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  username = COALESCE($2, username),
  email = COALESCE($3, email),
  password_hash = COALESCE($4, password_hash),
  first_name = COALESCE($5, first_name),
  last_name = COALESCE($6, last_name),
  updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;