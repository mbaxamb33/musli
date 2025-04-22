-- name: CreateProject :one
INSERT INTO projects (
  user_id,
  name,
  description,
  start_date,
  end_date,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetProject :one
SELECT * FROM projects
WHERE project_id = $1 LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProjectsByUser :many
SELECT * FROM projects
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProject :one
UPDATE projects
SET
  name = COALESCE($2, name),
  description = COALESCE($3, description),
  start_date = COALESCE($4, start_date),
  end_date = COALESCE($5, end_date),
  status = COALESCE($6, status),
  updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1
RETURNING *;

-- name: UpdateProjectStatus :one
UPDATE projects
SET 
  status = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;

-- name: CountProjects :one
SELECT COUNT(*) FROM projects;

-- name: CountProjectsByUser :one
SELECT COUNT(*) FROM projects
WHERE user_id = $1;

-- name: GetProjectsWithResourceCount :many
SELECT 
  p.*,
  COUNT(pr.project_resource_id) AS resource_count
FROM projects p
LEFT JOIN project_resources pr ON p.project_id = pr.project_id
WHERE p.user_id = $1
GROUP BY p.project_id
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;