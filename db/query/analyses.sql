-- name: CreateAnalysis :one
INSERT INTO analyses (
    sales_process_id, version
)
VALUES ($1, $2)
RETURNING analysis_id, sales_process_id, version, created_at, updated_at;

-- name: GetAnalysisByID :one
SELECT analysis_id, sales_process_id, version, created_at, updated_at
FROM analyses
WHERE analysis_id = $1;

-- name: GetLatestAnalysisBySalesProcess :one
SELECT analysis_id, sales_process_id, version, created_at, updated_at
FROM analyses
WHERE sales_process_id = $1
ORDER BY version DESC
LIMIT 1;

-- name: ListAnalysesBySalesProcess :many
SELECT analysis_id, sales_process_id, version, created_at, updated_at
FROM analyses
WHERE sales_process_id = $1
ORDER BY version DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAnalysis :one
UPDATE analyses
SET updated_at = CURRENT_TIMESTAMP
WHERE analysis_id = $1
RETURNING analysis_id, sales_process_id, version, created_at, updated_at;

-- name: DeleteAnalysis :exec
DELETE FROM analyses
WHERE analysis_id = $1;