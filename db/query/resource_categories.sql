-- name: CreateResourceCategory :one
INSERT INTO resource_categories (
  name,
  description
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetResourceCategory :one
SELECT * FROM resource_categories
WHERE category_id = $1 LIMIT 1;

-- name: ListResourceCategories :many
SELECT * FROM resource_categories
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateResourceCategory :one
UPDATE resource_categories
SET
  name = COALESCE($2, name),
  description = COALESCE($3, description),
  updated_at = CURRENT_TIMESTAMP
WHERE category_id = $1
RETURNING *;

-- name: DeleteResourceCategory :exec
DELETE FROM resource_categories
WHERE category_id = $1;

-- name: GetResourceCategoryWithResourceCount :many
SELECT 
  rc.*,
  COUNT(r.resource_id) AS resource_count
FROM resource_categories rc
LEFT JOIN resources r ON rc.category_id = r.category_id
GROUP BY rc.category_id
ORDER BY rc.name
LIMIT $1 OFFSET $2;

-- name: GetResourceCategoryByName :one
SELECT * FROM resource_categories
WHERE name = $1 LIMIT 1;