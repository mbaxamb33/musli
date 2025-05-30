-- =============================================================================
-- BRIEFS QUERIES
-- =============================================================================

-- name: CreateBrief :one
INSERT INTO briefs (
    master_brief_id, brief_type, brief_tag, title, text_content
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, master_brief_id, brief_type, brief_tag, title, text_content, created_at, updated_at;

-- name: GetBriefByID :one
SELECT id, master_brief_id, brief_type, brief_tag, title, text_content, created_at, updated_at
FROM briefs
WHERE id = $1;

-- name: ListBriefsByMasterBrief :many
SELECT id, master_brief_id, brief_type, brief_tag, title, text_content, created_at, updated_at
FROM briefs
WHERE master_brief_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListBriefsByType :many
SELECT id, master_brief_id, brief_type, brief_tag, title, text_content, created_at, updated_at
FROM briefs
WHERE master_brief_id = $1 AND brief_type = $2
ORDER BY created_at DESC;

-- name: ListBriefsByTag :many
SELECT id, master_brief_id, brief_type, brief_tag, title, text_content, created_at, updated_at
FROM briefs
WHERE master_brief_id = $1 AND brief_tag = $2
ORDER BY created_at DESC;

-- name: UpdateBrief :one
UPDATE briefs
SET brief_type = $2,
    brief_tag = $3,
    title = $4,
    text_content = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, master_brief_id, brief_type, brief_tag, title, text_content, created_at, updated_at;

-- name: DeleteBrief :exec
DELETE FROM briefs
WHERE id = $1;