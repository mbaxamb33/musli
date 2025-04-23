-- name: CreateCompanyNews :one
INSERT INTO company_news (
    company_id, title, publication_date, source, url, summary, sentiment, datasource_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING news_id, company_id, title, publication_date, source, url, summary, sentiment, datasource_id;

-- name: GetCompanyNewsByID :one
SELECT news_id, company_id, title, publication_date, source, url, summary, sentiment, datasource_id
FROM company_news
WHERE news_id = $1;

-- name: ListCompanyNewsByCompany :many
SELECT news_id, company_id, title, publication_date, source, url, summary, sentiment, datasource_id
FROM company_news
WHERE company_id = $1
ORDER BY publication_date DESC
LIMIT $2 OFFSET $3;

-- name: ListCompanyNewsBySentiment :many
SELECT news_id, company_id, title, publication_date, source, url, summary, sentiment, datasource_id
FROM company_news
WHERE sentiment = $1
ORDER BY publication_date DESC
LIMIT $2 OFFSET $3;

-- name: ListCompanyNewsByDatasource :many
SELECT news_id, company_id, title, publication_date, source, url, summary, sentiment, datasource_id
FROM company_news
WHERE datasource_id = $1
ORDER BY publication_date DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCompanyNews :one
UPDATE company_news
SET title = $2,
    publication_date = $3,
    source = $4,
    url = $5,
    summary = $6,
    sentiment = $7
WHERE news_id = $1
RETURNING news_id, company_id, title, publication_date, source, url, summary, sentiment, datasource_id;

-- name: DeleteCompanyNews :exec
DELETE FROM company_news
WHERE news_id = $1;
