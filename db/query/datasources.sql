-- name: CreateDatasource :one
INSERT INTO datasources (
    source_type, link, file_data, file_name
)
VALUES ($1, $2, $3, $4)
RETURNING datasource_id, source_type, link, file_data, file_name, created_at;

-- name: GetDatasourceByID :one
SELECT datasource_id, source_type, link, file_name, created_at
FROM datasources
WHERE datasource_id = $1;

-- name: ListDatasources :many
SELECT datasource_id, source_type, link, file_name, created_at
FROM datasources
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListDatasourcesByType :many
SELECT datasource_id, source_type, link, file_name, created_at
FROM datasources
WHERE source_type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteDatasource :exec
DELETE FROM datasources
WHERE datasource_id = $1;

-- name: GetFullDatasourceByID :one
SELECT datasource_id, source_type, link, file_data, file_name, created_at
FROM datasources
WHERE datasource_id = $1;