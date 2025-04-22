-- name: CreateResource :one
INSERT INTO resources (
  name,
  description,
  category_id,
  unit,
  cost_per_unit
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetResource :one
SELECT * FROM resources
WHERE resource_id = $1 LIMIT 1;

-- name: ListResources :many
SELECT * FROM resources
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListResourcesByCategory :many
SELECT * FROM resources
WHERE category_id = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateResource :one
UPDATE resources
SET
  name = COALESCE($2, name),
  description = COALESCE($3, description),
  category_id = COALESCE($4, category_id),
  unit = COALESCE($5, unit),
  cost_per_unit = COALESCE($6, cost_per_unit),
  updated_at = CURRENT_TIMESTAMP
WHERE resource_id = $1
RETURNING *;

-- name: DeleteResource :exec
DELETE FROM resources
WHERE resource_id = $1;

-- name: SearchResources :many
SELECT * FROM resources
WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: ListResourcesWithCategory :many
SELECT 
  r.*,
  rc.name AS category_name
FROM resources r
LEFT JOIN resource_categories rc ON r.category_id = rc.category_id
ORDER BY r.name
LIMIT $1 OFFSET $2;

-- name: GetResourcesUsedInProject :many
SELECT 
  r.*,
  pr.quantity,
  pr.notes
FROM resources r
JOIN project_resources pr ON r.resource_id = pr.resource_id
WHERE pr.project_id = $1
ORDER BY r.name;