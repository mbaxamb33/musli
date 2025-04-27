-- name: CreateParagraph :one
INSERT INTO paragraphs (
    datasource_id, title, main_idea, content
)
VALUES ($1, $2, $3, $4)
RETURNING paragraph_id, datasource_id, title, main_idea, content, created_at;

-- name: GetParagraphByID :one
SELECT paragraph_id, datasource_id, title, main_idea, content, created_at
FROM paragraphs
WHERE paragraph_id = $1;

-- name: ListParagraphsByDatasource :many
SELECT paragraph_id, datasource_id, title, main_idea, content, created_at
FROM paragraphs
WHERE datasource_id = $1
ORDER BY paragraph_id ASC
LIMIT $2 OFFSET $3;

-- name: SearchParagraphsByContent :many
SELECT paragraph_id, datasource_id, title, main_idea, content, created_at
FROM paragraphs
WHERE content ILIKE '%' || $1 || '%' OR main_idea ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateParagraph :one
UPDATE paragraphs
SET title = $2,
    main_idea = $3,
    content = $4
WHERE paragraph_id = $1
RETURNING paragraph_id, datasource_id, title, main_idea, content, created_at;

-- name: DeleteParagraph :exec
DELETE FROM paragraphs
WHERE paragraph_id = $1;