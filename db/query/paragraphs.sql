-- name: CreateParagraph :one
INSERT INTO paragraphs (
    datasource_id, content, main_idea, classification, confidence_score
)
VALUES ($1, $2, $3, $4, $5)
RETURNING paragraph_id, datasource_id, content, main_idea, classification, confidence_score;

-- name: GetParagraphByID :one
SELECT paragraph_id, datasource_id, content, main_idea, classification, confidence_score
FROM paragraphs
WHERE paragraph_id = $1;

-- name: ListParagraphsByDatasource :many
SELECT paragraph_id, datasource_id, content, main_idea, classification, confidence_score
FROM paragraphs
WHERE datasource_id = $1
ORDER BY paragraph_id ASC
LIMIT $2 OFFSET $3;

-- name: ListParagraphsByClassification :many
SELECT paragraph_id, datasource_id, content, main_idea, classification, confidence_score
FROM paragraphs
WHERE classification = $1
ORDER BY confidence_score DESC
LIMIT $2 OFFSET $3;

-- name: UpdateParagraph :one
UPDATE paragraphs
SET content = $2,
    main_idea = $3,
    classification = $4,
    confidence_score = $5
WHERE paragraph_id = $1
RETURNING paragraph_id, datasource_id, content, main_idea, classification, confidence_score;

-- name: DeleteParagraph :exec
DELETE FROM paragraphs
WHERE paragraph_id = $1;
