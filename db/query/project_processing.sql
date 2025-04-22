-- name: AddProcessingToProject :one
INSERT INTO project_processing (
  project_id,
  processing_type,
  description,
  status
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetProjectProcessing :one
SELECT * FROM project_processing
WHERE processing_id = $1 LIMIT 1;

-- name: ListProjectProcessing :many
SELECT * FROM project_processing
WHERE project_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListProjectProcessingByType :many
SELECT * FROM project_processing
WHERE project_id = $1 AND processing_type = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListProjectProcessingByStatus :many
SELECT * FROM project_processing
WHERE project_id = $1 AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateProjectProcessing :one
UPDATE project_processing
SET
  processing_type = COALESCE($2, processing_type),
  description = COALESCE($3, description),
  status = COALESCE($4, status),
  updated_at = CURRENT_TIMESTAMP
WHERE processing_id = $1
RETURNING *;

-- name: UpdateProjectProcessingStatus :one
UPDATE project_processing
SET 
  status = $2,
  updated_at = CURRENT_TIMESTAMP
WHERE processing_id = $1
RETURNING *;

-- name: DeleteProjectProcessing :exec
DELETE FROM project_processing
WHERE processing_id = $1;

-- name: CountProjectProcessingByStatus :one
SELECT COUNT(*) 
FROM project_processing
WHERE project_id = $1 AND status = $2;

-- name: GetProcessingTypeCounts :many
SELECT 
  processing_type, 
  COUNT(*) as count
FROM project_processing
WHERE project_id = $1
GROUP BY processing_type
ORDER BY count DESC;