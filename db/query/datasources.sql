-- name: CreateDatasource :one
INSERT INTO datasources (
    source_type, source_id
)
VALUES ($1, $2)
RETURNING datasource_id, source_type, source_id, extraction_timestamp;

-- name: GetDatasourceByID :one
SELECT datasource_id, source_type, source_id, extraction_timestamp
FROM datasources
WHERE datasource_id = $1;

-- name: ListDatasources :many
SELECT datasource_id, source_type, source_id, extraction_timestamp
FROM datasources
ORDER BY extraction_timestamp DESC
LIMIT $1 OFFSET $2;

-- name: ListDatasourcesByType :many
SELECT datasource_id, source_type, source_id, extraction_timestamp
FROM datasources
WHERE source_type = $1
ORDER BY extraction_timestamp DESC
LIMIT $2 OFFSET $3;

-- name: DeleteDatasource :exec
DELETE FROM datasources
WHERE datasource_id = $1;
