-- name: CreateContactNewsItem :one
INSERT INTO contact_news (
    contact_id, title, content, datasource_id
)
VALUES ($1, $2, $3, $4)
RETURNING contact_news_id, contact_id, title, content, datasource_id, created_at;

-- name: GetContactNewsItemByID :one
SELECT contact_news_id, contact_id, title, content, datasource_id, created_at
FROM contact_news
WHERE contact_news_id = $1;

-- name: ListNewsItemsByContact :many
SELECT contact_news_id, contact_id, title, content, datasource_id, created_at
FROM contact_news
WHERE contact_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateContactNewsItem :one
UPDATE contact_news
SET title = $2,
    content = $3
WHERE contact_news_id = $1
RETURNING contact_news_id, contact_id, title, content, datasource_id, created_at;

-- name: DeleteContactNewsItem :exec
DELETE FROM contact_news
WHERE contact_news_id = $1;