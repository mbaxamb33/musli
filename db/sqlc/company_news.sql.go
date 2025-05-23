// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: company_news.sql

package db

import (
	"context"
	"database/sql"
)

const createCompanyNews = `-- name: CreateCompanyNews :one
INSERT INTO company_news (
    company_id, title, content, datasource_id
)
VALUES ($1, $2, $3, $4)
RETURNING company_news_id, company_id, title, content, datasource_id, created_at
`

type CreateCompanyNewsParams struct {
	CompanyID    int32          `json:"company_id"`
	Title        string         `json:"title"`
	Content      sql.NullString `json:"content"`
	DatasourceID sql.NullInt32  `json:"datasource_id"`
}

func (q *Queries) CreateCompanyNews(ctx context.Context, arg CreateCompanyNewsParams) (CompanyNews, error) {
	row := q.db.QueryRowContext(ctx, createCompanyNews,
		arg.CompanyID,
		arg.Title,
		arg.Content,
		arg.DatasourceID,
	)
	var i CompanyNews
	err := row.Scan(
		&i.CompanyNewsID,
		&i.CompanyID,
		&i.Title,
		&i.Content,
		&i.DatasourceID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteCompanyNews = `-- name: DeleteCompanyNews :exec
DELETE FROM company_news
WHERE company_news_id = $1
`

func (q *Queries) DeleteCompanyNews(ctx context.Context, companyNewsID int32) error {
	_, err := q.db.ExecContext(ctx, deleteCompanyNews, companyNewsID)
	return err
}

const getCompanyNewsByID = `-- name: GetCompanyNewsByID :one
SELECT company_news_id, company_id, title, content, datasource_id, created_at
FROM company_news
WHERE company_news_id = $1
`

func (q *Queries) GetCompanyNewsByID(ctx context.Context, companyNewsID int32) (CompanyNews, error) {
	row := q.db.QueryRowContext(ctx, getCompanyNewsByID, companyNewsID)
	var i CompanyNews
	err := row.Scan(
		&i.CompanyNewsID,
		&i.CompanyID,
		&i.Title,
		&i.Content,
		&i.DatasourceID,
		&i.CreatedAt,
	)
	return i, err
}

const listNewsByCompany = `-- name: ListNewsByCompany :many
SELECT company_news_id, company_id, title, content, datasource_id, created_at
FROM company_news
WHERE company_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type ListNewsByCompanyParams struct {
	CompanyID int32 `json:"company_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListNewsByCompany(ctx context.Context, arg ListNewsByCompanyParams) ([]CompanyNews, error) {
	rows, err := q.db.QueryContext(ctx, listNewsByCompany, arg.CompanyID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CompanyNews
	for rows.Next() {
		var i CompanyNews
		if err := rows.Scan(
			&i.CompanyNewsID,
			&i.CompanyID,
			&i.Title,
			&i.Content,
			&i.DatasourceID,
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

const updateCompanyNews = `-- name: UpdateCompanyNews :one
UPDATE company_news
SET title = $2,
    content = $3
WHERE company_news_id = $1
RETURNING company_news_id, company_id, title, content, datasource_id, created_at
`

type UpdateCompanyNewsParams struct {
	CompanyNewsID int32          `json:"company_news_id"`
	Title         string         `json:"title"`
	Content       sql.NullString `json:"content"`
}

func (q *Queries) UpdateCompanyNews(ctx context.Context, arg UpdateCompanyNewsParams) (CompanyNews, error) {
	row := q.db.QueryRowContext(ctx, updateCompanyNews, arg.CompanyNewsID, arg.Title, arg.Content)
	var i CompanyNews
	err := row.Scan(
		&i.CompanyNewsID,
		&i.CompanyID,
		&i.Title,
		&i.Content,
		&i.DatasourceID,
		&i.CreatedAt,
	)
	return i, err
}
