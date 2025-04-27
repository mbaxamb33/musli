-- name: CreateProject :one
INSERT INTO projects (
    user_id, project_name, main_idea
)
VALUES ($1, $2, $3)
RETURNING project_id, user_id, project_name, main_idea, created_at, updated_at;

-- name: GetProjectByID :one
SELECT project_id, user_id, project_name, main_idea, created_at, updated_at
FROM projects
WHERE project_id = $1;

-- name: ListProjectsByUserID :many
SELECT project_id, user_id, project_name, main_idea, created_at, updated_at
FROM projects
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchProjectsByName :many
SELECT project_id, user_id, project_name, main_idea, created_at, updated_at
FROM projects
WHERE user_id = $1 AND project_name ILIKE '%' || $2 || '%'
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateProject :one
UPDATE projects
SET project_name = $2,
    main_idea = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1
RETURNING project_id, user_id, project_name, main_idea, created_at, updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;