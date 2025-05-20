-- name: CreateProject :one
INSERT INTO projects (
    cognito_sub, project_name, main_idea
)
VALUES ($1, $2, $3)
RETURNING project_id, cognito_sub, project_name, main_idea, created_at, updated_at;

-- name: GetProjectByID :one
SELECT project_id, cognito_sub, project_name, main_idea, created_at, updated_at
FROM projects
WHERE project_id = $1;

-- name: ListProjectsByCognitoSub :many
SELECT project_id, cognito_sub, project_name, main_idea, created_at, updated_at
FROM projects
WHERE cognito_sub = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchProjectsByName :many
SELECT project_id, cognito_sub, project_name, main_idea, created_at, updated_at
FROM projects
WHERE cognito_sub = $1 AND project_name ILIKE '%' || $2 || '%'
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateProject :one
UPDATE projects
SET project_name = $2,
    main_idea = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1
RETURNING project_id, cognito_sub, project_name, main_idea, created_at, updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;