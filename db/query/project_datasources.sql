-- name: AssociateDatasourceWithProject :exec
INSERT INTO project_datasources (
    project_id, datasource_id
)
VALUES ($1, $2);

-- name: GetProjectDatasourceAssociation :one
SELECT project_id, datasource_id
FROM project_datasources
WHERE project_id = $1 AND datasource_id = $2;

-- name: ListDatasourcesByProject :many
SELECT d.datasource_id, d.source_type, d.link, d.file_name, d.created_at
FROM datasources d
JOIN project_datasources pd ON d.datasource_id = pd.datasource_id
WHERE pd.project_id = $1
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListProjectsByDatasource :many
SELECT p.project_id, p.cognito_sub, p.project_name, p.main_idea, p.created_at, p.updated_at
FROM projects p
JOIN project_datasources pd ON p.project_id = pd.project_id
WHERE pd.datasource_id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: RemoveDatasourceFromProject :exec
DELETE FROM project_datasources
WHERE project_id = $1 AND datasource_id = $2;