-- name: CreateSalesProcess :one
INSERT INTO sales_processes (
    user_id, contact_id, overall_matching_score, status
)
VALUES ($1, $2, $3, $4)
RETURNING sales_process_id, user_id, contact_id, overall_matching_score, status, created_at, updated_at;

-- name: GetSalesProcessByID :one
SELECT sales_process_id, user_id, contact_id, overall_matching_score, status, created_at, updated_at
FROM sales_processes
WHERE sales_process_id = $1;

-- name: ListSalesProcessesByUser :many
SELECT sales_process_id, user_id, contact_id, overall_matching_score, status, created_at, updated_at
FROM sales_processes
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListSalesProcessesByContact :many
SELECT sales_process_id, user_id, contact_id, overall_matching_score, status, created_at, updated_at
FROM sales_processes
WHERE contact_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListSalesProcessesByStatus :many
SELECT sales_process_id, user_id, contact_id, overall_matching_score, status, created_at, updated_at
FROM sales_processes
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateSalesProcess :one
UPDATE sales_processes
SET overall_matching_score = $2,
    status = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE sales_process_id = $1
RETURNING sales_process_id, user_id, contact_id, overall_matching_score, status, created_at, updated_at;

-- name: DeleteSalesProcess :exec
DELETE FROM sales_processes
WHERE sales_process_id = $1;