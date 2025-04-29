-- name: AssociateDatasourceWithContact :exec
INSERT INTO contact_datasources (
    contact_id, datasource_id
)
VALUES ($1, $2);

-- name: GetContactDatasourceAssociation :one
SELECT contact_id, datasource_id, created_at
FROM contact_datasources
WHERE contact_id = $1 AND datasource_id = $2;

-- name: ListDatasourcesByContact :many
SELECT d.datasource_id, d.source_type, d.link, d.file_name, d.created_at
FROM datasources d
JOIN contact_datasources cd ON d.datasource_id = cd.datasource_id
WHERE cd.contact_id = $1
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListContactsByDatasource :many
SELECT c.contact_id, c.company_id, c.first_name, c.last_name, c.position, c.email, c.phone, c.notes, c.created_at
FROM contacts c
JOIN contact_datasources cd ON c.contact_id = cd.contact_id
WHERE cd.datasource_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: RemoveDatasourceFromContact :exec
DELETE FROM contact_datasources
WHERE contact_id = $1 AND datasource_id = $2;