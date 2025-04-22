-- name: AddResourceToProject :one
INSERT INTO project_resources (
  project_id,
  resource_id,
  quantity,
  notes
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetProjectResource :one
SELECT * FROM project_resources
WHERE project_id = $1 AND resource_id = $2
LIMIT 1;

-- name: GetProjectResourceById :one
SELECT * FROM project_resources
WHERE project_resource_id = $1
LIMIT 1;

-- name: ListProjectResources :many
SELECT 
  pr.*,
  r.name AS resource_name,
  r.unit,
  r.cost_per_unit,
  (pr.quantity * r.cost_per_unit) AS total_cost
FROM project_resources pr
JOIN resources r ON pr.resource_id = r.resource_id
WHERE pr.project_id = $1
ORDER BY r.name
LIMIT $2 OFFSET $3;

-- name: UpdateProjectResource :one
UPDATE project_resources
SET
  quantity = COALESCE($3, quantity),
  notes = COALESCE($4, notes),
  updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1 AND resource_id = $2
RETURNING *;

-- name: UpdateProjectResourceById :one
UPDATE project_resources
SET
  quantity = COALESCE($2, quantity),
  notes = COALESCE($3, notes),
  updated_at = CURRENT_TIMESTAMP
WHERE project_resource_id = $1
RETURNING *;

-- name: RemoveResourceFromProject :exec
DELETE FROM project_resources
WHERE project_id = $1 AND resource_id = $2;

-- name: RemoveProjectResourceById :exec
DELETE FROM project_resources
WHERE project_resource_id = $1;

-- name: GetTotalProjectCost :one
SELECT 
  SUM(pr.quantity * r.cost_per_unit) AS total_cost
FROM project_resources pr
JOIN resources r ON pr.resource_id = r.resource_id
WHERE pr.project_id = $1;

-- name: GetResourcesNotInProject :many
SELECT 
  r.*
FROM resources r
WHERE r.resource_id NOT IN (
  SELECT resource_id FROM project_resources WHERE project_id = $1
)
ORDER BY r.name
LIMIT $2 OFFSET $3;