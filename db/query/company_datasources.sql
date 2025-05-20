-- name: AssociateDatasourceWithCompany :exec
INSERT INTO company_datasources (
    company_id, datasource_id
)
VALUES ($1, $2);

-- name: GetCompanyDatasourceAssociation :one
SELECT company_id, datasource_id, created_at
FROM company_datasources
WHERE company_id = $1 AND datasource_id = $2;

-- name: ListDatasourcesByCompany :many
SELECT d.datasource_id, d.source_type, d.link, d.file_name, d.created_at
FROM datasources d
JOIN company_datasources cd ON d.datasource_id = cd.datasource_id
WHERE cd.company_id = $1
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCompaniesByDatasource :many
SELECT c.company_id, c.cognito_sub, c.company_name, c.industry, c.website, c.address, c.description, c.created_at
FROM companies c
JOIN company_datasources cd ON c.company_id = cd.company_id
WHERE cd.datasource_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: RemoveDatasourceFromCompany :exec
DELETE FROM company_datasources
WHERE company_id = $1 AND datasource_id = $2;