-- name: CreateProcessedCompanyData :one
INSERT INTO processed_company_data (
  company_id,
  data_type,
  data_key,
  data_value,
  confidence_score,
  source_id
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetProcessedCompanyData :one
SELECT * FROM processed_company_data
WHERE data_id = $1 LIMIT 1;

-- name: GetProcessedCompanyDataByKey :one
SELECT * FROM processed_company_data
WHERE company_id = $1 AND data_type = $2 AND data_key = $3
ORDER BY confidence_score DESC, updated_at DESC
LIMIT 1;

-- name: ListProcessedCompanyData :many
SELECT * FROM processed_company_data
WHERE company_id = $1
ORDER BY data_type, data_key, confidence_score DESC
LIMIT $2 OFFSET $3;

-- name: ListProcessedCompanyDataByType :many
SELECT * FROM processed_company_data
WHERE company_id = $1 AND data_type = $2
ORDER BY data_key, confidence_score DESC
LIMIT $3 OFFSET $4;

-- name: UpdateProcessedCompanyData :one
UPDATE processed_company_data
SET
  data_value = COALESCE($5, data_value),
  confidence_score = COALESCE($6, confidence_score),
  source_id = COALESCE($7, source_id),
  updated_at = CURRENT_TIMESTAMP
WHERE company_id = $1 AND data_type = $2 AND data_key = $3 AND data_id = $4
RETURNING *;

-- name: UpsertProcessedCompanyData :one
INSERT INTO processed_company_data (
  company_id,
  data_type,
  data_key,
  data_value,
  confidence_score,
  source_id
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (company_id, data_type, data_key) 
WHERE data_id = (
  SELECT data_id FROM processed_company_data 
  WHERE company_id = $1 AND data_type = $2 AND data_key = $3
  ORDER BY confidence_score DESC, updated_at DESC 
  LIMIT 1
)
DO UPDATE SET
  data_value = EXCLUDED.data_value,
  confidence_score = EXCLUDED.confidence_score,
  source_id = EXCLUDED.source_id,
  updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: DeleteProcessedCompanyData :exec
DELETE FROM processed_company_data
WHERE data_id = $1;

-- name: DeleteProcessedCompanyDataByKey :exec
DELETE FROM processed_company_data
WHERE company_id = $1 AND data_type = $2 AND data_key = $3;

-- name: DeleteAllProcessedCompanyData :exec
DELETE FROM processed_company_data
WHERE company_id = $1;

-- name: GetCompanyDataTypes :many
SELECT DISTINCT data_type
FROM processed_company_data
WHERE company_id = $1
ORDER BY data_type;

-- name: GetCompanyDataKeysByType :many
SELECT DISTINCT data_key
FROM processed_company_data
WHERE company_id = $1 AND data_type = $2
ORDER BY data_key;

-- name: GetCompanyDataSummary :many
SELECT 
  data_type, 
  COUNT(DISTINCT data_key) as key_count
FROM processed_company_data
WHERE company_id = $1
GROUP BY data_type
ORDER BY key_count DESC;

-- name: GetHighConfidenceCompanyData :many
SELECT *
FROM processed_company_data
WHERE company_id = $1 AND confidence_score >= $2
ORDER BY data_type, data_key
LIMIT $3 OFFSET $4;