-- name: CreateDataSource :one
INSERT INTO data_sources (
  name,
  url_pattern,
  api_endpoint,
  api_key,
  is_active
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetDataSource :one
SELECT * FROM data_sources
WHERE source_id = $1 LIMIT 1;

-- name: ListDataSources :many
SELECT * FROM data_sources
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListActiveDataSources :many
SELECT * FROM data_sources
WHERE is_active = TRUE
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateDataSource :one
UPDATE data_sources
SET
  name = COALESCE($2, name),
  url_pattern = COALESCE($3, url_pattern),
  api_endpoint = COALESCE($4, api_endpoint),
  api_key = COALESCE($5, api_key),
  is_active = COALESCE($6, is_active),
  updated_at = CURRENT_TIMESTAMP
WHERE source_id = $1
RETURNING *;

-- name: ToggleDataSourceActive :one
UPDATE data_sources
SET 
  is_active = NOT is_active,
  updated_at = CURRENT_TIMESTAMP
WHERE source_id = $1
RETURNING *;

-- name: DeleteDataSource :exec
DELETE FROM data_sources
WHERE source_id = $1;

-- name: GetDataSourceByName :one
SELECT * FROM data_sources
WHERE name = $1 LIMIT 1;

-- name: GetDataSourceWithUsageCount :many
SELECT 
  ds.*,
  COUNT(pcd.source_id) AS usage_count
FROM data_sources ds
LEFT JOIN web_scrape_data wsd ON ds.source_id = wsd.source_id
LEFT JOIN processed_company_data pcd ON ds.source_id = pcd.source_id
GROUP BY ds.source_id
ORDER BY ds.name
LIMIT $1 OFFSET $2;

-- name: GetDataSourcesWithApiEndpoint :many
SELECT * FROM data_sources
WHERE api_endpoint IS NOT NULL AND api_endpoint != ''
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: GetDataSourcesByUrlPattern :many
SELECT * FROM data_sources
WHERE url_pattern LIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;