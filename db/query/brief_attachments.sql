-- =============================================================================
-- BRIEF ATTACHMENTS QUERIES
-- =============================================================================

-- name: CreateBriefAttachment :one
INSERT INTO brief_attachments (
    brief_id, datasource_id, attachment_type
)
VALUES ($1, $2, $3)
RETURNING id, brief_id, datasource_id, attachment_type, created_at;

-- name: GetBriefAttachmentByID :one
SELECT id, brief_id, datasource_id, attachment_type, created_at
FROM brief_attachments
WHERE id = $1;

-- name: ListAttachmentsByBrief :many
SELECT id, brief_id, datasource_id, attachment_type, created_at
FROM brief_attachments
WHERE brief_id = $1
ORDER BY created_at DESC;

-- name: ListAttachmentsByType :many
SELECT id, brief_id, datasource_id, attachment_type, created_at
FROM brief_attachments
WHERE brief_id = $1 AND attachment_type = $2
ORDER BY created_at DESC;

-- name: DeleteBriefAttachment :exec
DELETE FROM brief_attachments
WHERE id = $1;