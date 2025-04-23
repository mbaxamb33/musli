-- name: CreateProject :one
INSERT INTO projects (user_id, project_name, description)
VALUES ($1, $2, $3)
RETURNING project_id, user_id, project_name, description, created_at, last_updated_at;

-- name: GetProjectByID :one
SELECT project_id, user_id, project_name, description, created_at, last_updated_at
FROM projects
WHERE project_id = $1;

-- name: ListProjectsByUser :many
SELECT project_id, user_id, project_name, description, created_at, last_updated_at
FROM projects
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProject :one
UPDATE projects
SET project_name = $2,
    description = $3,
    last_updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1
RETURNING project_id, user_id, project_name, description, created_at, last_updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;
