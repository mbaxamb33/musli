-- name: CreateCompanyNews :one
INSERT INTO company_news (
    company_id, title, content, datasource_id
)
VALUES ($1, $2, $3, $4)
RETURNING company_news_id, company_id, title, content, datasource_id, created_at;

-- name: GetCompanyNewsByID :one
SELECT company_news_id, company_id, title, content, datasource_id, created_at
FROM company_news
WHERE company_news_id = $1;

-- name: ListNewsByCompany :many
SELECT company_news_id, company_id, title, content, datasource_id, created_at
FROM company_news
WHERE company_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCompanyNews :one
UPDATE company_news
SET title = $2,
    content = $3
WHERE company_news_id = $1
RETURNING company_news_id, company_id, title, content, datasource_id, created_at;

-- name: DeleteCompanyNews :exec
DELETE FROM company_news
WHERE company_news_id = $1;