-- name: CreateAnalysisInput :one
INSERT INTO analysis_inputs (
    analysis_id, input_type, datasource_id, content
)
VALUES ($1, $2, $3, $4)
RETURNING input_id, analysis_id, input_type, datasource_id, content, created_at;

-- name: GetAnalysisInputByID :one
SELECT input_id, analysis_id, input_type, datasource_id, content, created_at
FROM analysis_inputs
WHERE input_id = $1;

-- name: ListInputsByAnalysis :many
SELECT input_id, analysis_id, input_type, datasource_id, content, created_at
FROM analysis_inputs
WHERE analysis_id = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: ListInputsByAnalysisAndType :many
SELECT input_id, analysis_id, input_type, datasource_id, content, created_at
FROM analysis_inputs
WHERE analysis_id = $1 AND input_type = $2
ORDER BY created_at ASC
LIMIT $3 OFFSET $4;

-- name: UpdateAnalysisInput :one
UPDATE analysis_inputs
SET content = $2
WHERE input_id = $1
RETURNING input_id, analysis_id, input_type, datasource_id, content, created_at;

-- name: DeleteAnalysisInput :exec
DELETE FROM analysis_inputs
WHERE input_id = $1;