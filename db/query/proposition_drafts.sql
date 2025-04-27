-- name: CreatePropositionDraft :one
INSERT INTO proposition_drafts (
    sales_process_id, title, content, version
)
VALUES ($1, $2, $3, $4)
RETURNING draft_id, sales_process_id, title, content, version, created_at, updated_at;

-- name: GetPropositionDraftByID :one
SELECT draft_id, sales_process_id, title, content, version, created_at, updated_at
FROM proposition_drafts
WHERE draft_id = $1;

-- name: GetLatestPropositionDraft :one
SELECT draft_id, sales_process_id, title, content, version, created_at, updated_at
FROM proposition_drafts
WHERE sales_process_id = $1
ORDER BY version DESC
LIMIT 1;

-- name: ListPropositionDraftsBySalesProcess :many
SELECT draft_id, sales_process_id, title, content, version, created_at, updated_at
FROM proposition_drafts
WHERE sales_process_id = $1
ORDER BY version DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePropositionDraft :one
UPDATE proposition_drafts
SET title = $2,
    content = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE draft_id = $1
RETURNING draft_id, sales_process_id, title, content, version, created_at, updated_at;

-- name: DeletePropositionDraft :exec
DELETE FROM proposition_drafts
WHERE draft_id = $1;