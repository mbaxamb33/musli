// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: contact_datasources.sql

package db

import (
	"context"
	"database/sql"
)

const associateDatasourceWithContact = `-- name: AssociateDatasourceWithContact :exec
INSERT INTO contact_datasources (
    contact_id, datasource_id
)
VALUES ($1, $2)
`

type AssociateDatasourceWithContactParams struct {
	ContactID    int32 `json:"contact_id"`
	DatasourceID int32 `json:"datasource_id"`
}

func (q *Queries) AssociateDatasourceWithContact(ctx context.Context, arg AssociateDatasourceWithContactParams) error {
	_, err := q.db.ExecContext(ctx, associateDatasourceWithContact, arg.ContactID, arg.DatasourceID)
	return err
}

const getContactDatasourceAssociation = `-- name: GetContactDatasourceAssociation :one
SELECT contact_id, datasource_id, created_at
FROM contact_datasources
WHERE contact_id = $1 AND datasource_id = $2
`

type GetContactDatasourceAssociationParams struct {
	ContactID    int32 `json:"contact_id"`
	DatasourceID int32 `json:"datasource_id"`
}

func (q *Queries) GetContactDatasourceAssociation(ctx context.Context, arg GetContactDatasourceAssociationParams) (ContactDatasource, error) {
	row := q.db.QueryRowContext(ctx, getContactDatasourceAssociation, arg.ContactID, arg.DatasourceID)
	var i ContactDatasource
	err := row.Scan(&i.ContactID, &i.DatasourceID, &i.CreatedAt)
	return i, err
}

const listContactsByDatasource = `-- name: ListContactsByDatasource :many
SELECT c.contact_id, c.company_id, c.first_name, c.last_name, c.position, c.email, c.phone, c.notes, c.created_at
FROM contacts c
JOIN contact_datasources cd ON c.contact_id = cd.contact_id
WHERE cd.datasource_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3
`

type ListContactsByDatasourceParams struct {
	DatasourceID int32 `json:"datasource_id"`
	Limit        int32 `json:"limit"`
	Offset       int32 `json:"offset"`
}

func (q *Queries) ListContactsByDatasource(ctx context.Context, arg ListContactsByDatasourceParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, listContactsByDatasource, arg.DatasourceID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ContactID,
			&i.CompanyID,
			&i.FirstName,
			&i.LastName,
			&i.Position,
			&i.Email,
			&i.Phone,
			&i.Notes,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDatasourcesByContact = `-- name: ListDatasourcesByContact :many
SELECT d.datasource_id, d.source_type, d.link, d.file_name, d.created_at
FROM datasources d
JOIN contact_datasources cd ON d.datasource_id = cd.datasource_id
WHERE cd.contact_id = $1
ORDER BY d.created_at DESC
LIMIT $2 OFFSET $3
`

type ListDatasourcesByContactParams struct {
	ContactID int32 `json:"contact_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

type ListDatasourcesByContactRow struct {
	DatasourceID int32          `json:"datasource_id"`
	SourceType   DatasourceType `json:"source_type"`
	Link         sql.NullString `json:"link"`
	FileName     sql.NullString `json:"file_name"`
	CreatedAt    sql.NullTime   `json:"created_at"`
}

func (q *Queries) ListDatasourcesByContact(ctx context.Context, arg ListDatasourcesByContactParams) ([]ListDatasourcesByContactRow, error) {
	rows, err := q.db.QueryContext(ctx, listDatasourcesByContact, arg.ContactID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDatasourcesByContactRow
	for rows.Next() {
		var i ListDatasourcesByContactRow
		if err := rows.Scan(
			&i.DatasourceID,
			&i.SourceType,
			&i.Link,
			&i.FileName,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeDatasourceFromContact = `-- name: RemoveDatasourceFromContact :exec
DELETE FROM contact_datasources
WHERE contact_id = $1 AND datasource_id = $2
`

type RemoveDatasourceFromContactParams struct {
	ContactID    int32 `json:"contact_id"`
	DatasourceID int32 `json:"datasource_id"`
}

func (q *Queries) RemoveDatasourceFromContact(ctx context.Context, arg RemoveDatasourceFromContactParams) error {
	_, err := q.db.ExecContext(ctx, removeDatasourceFromContact, arg.ContactID, arg.DatasourceID)
	return err
}
