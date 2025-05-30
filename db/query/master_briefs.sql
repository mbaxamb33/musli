-- =============================================================================
-- MASTER BRIEFS QUERIES
-- =============================================================================

-- name: CreateMasterBrief :one
INSERT INTO master_briefs (
    cognito_sub, company_id, contact_id, company_reference, contact_reference
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, cognito_sub, company_id, contact_id, company_reference, contact_reference, created_at, updated_at;

-- name: GetMasterBriefByID :one
SELECT id, cognito_sub, company_id, contact_id, company_reference, contact_reference, created_at, updated_at
FROM master_briefs
WHERE id = $1;

-- name: ListMasterBriefsByUser :many
SELECT id, cognito_sub, company_id, contact_id, company_reference, contact_reference, created_at, updated_at
FROM master_briefs
WHERE cognito_sub = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListMasterBriefsByCompany :many
SELECT id, cognito_sub, company_id, contact_id, company_reference, contact_reference, created_at, updated_at
FROM master_briefs
WHERE company_id = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListMasterBriefsByContact :many
SELECT id, cognito_sub, company_id, contact_id, company_reference, contact_reference, created_at, updated_at
FROM master_briefs
WHERE contact_id = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMasterBrief :one
UPDATE master_briefs
SET company_id = $2,
    contact_id = $3,
    company_reference = $4,
    contact_reference = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, cognito_sub, company_id, contact_id, company_reference, contact_reference, created_at, updated_at;

-- name: DeleteMasterBrief :exec
DELETE FROM master_briefs
WHERE id = $1;
