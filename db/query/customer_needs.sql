-- name: CreateCustomerNeed :one
INSERT INTO customer_needs (
    sales_process_id, need_description
)
VALUES ($1, $2)
RETURNING need_id, sales_process_id, need_description, created_at;

-- name: GetCustomerNeedByID :one
SELECT need_id, sales_process_id, need_description, created_at
FROM customer_needs
WHERE need_id = $1;

-- name: ListNeedsBySalesProcess :many
SELECT need_id, sales_process_id, need_description, created_at
FROM customer_needs
WHERE sales_process_id = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: UpdateCustomerNeed :one
UPDATE customer_needs
SET need_description = $2
WHERE need_id = $1
RETURNING need_id, sales_process_id, need_description, created_at;

-- name: DeleteCustomerNeed :exec
DELETE FROM customer_needs
WHERE need_id = $1;