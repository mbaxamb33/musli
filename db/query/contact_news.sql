-- name: CreateContactNews :one
INSERT INTO contact_news (
    contact_id, title, publication_date, source, url, summary, datasource_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING mention_id, contact_id, title, publication_date, source, url, summary, datasource_id;

-- name: GetContactNewsByID :one
SELECT mention_id, contact_id, title, publication_date, source, url, summary, datasource_id
FROM contact_news
WHERE mention_id = $1;

-- name: ListContactNewsByContact :many
SELECT mention_id, contact_id, title, publication_date, source, url, summary, datasource_id
FROM contact_news
WHERE contact_id = $1
ORDER BY publication_date DESC
LIMIT $2 OFFSET $3;

-- name: ListContactNewsByDatasource :many
SELECT mention_id, contact_id, title, publication_date, source, url, summary, datasource_id
FROM contact_news
WHERE datasource_id = $1
ORDER BY publication_date DESC
LIMIT $2 OFFSET $3;

-- name: ListContactNewsBySource :many
SELECT mention_id, contact_id, title, publication_date, source, url, summary, datasource_id
FROM contact_news
WHERE source = $1
ORDER BY publication_date DESC
LIMIT $2 OFFSET $3;

-- name: UpdateContactNews :one
UPDATE contact_news
SET title = $2,
    publication_date = $3,
    source = $4,
    url = $5,
    summary = $6
WHERE mention_id = $1
RETURNING mention_id, contact_id, title, publication_date, source, url, summary, datasource_id;

-- name: DeleteContactNews :exec
DELETE FROM contact_news
WHERE mention_id = $1;
