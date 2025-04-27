-- name: CreateTask :one
INSERT INTO tasks (
    sales_process_id, title, description, status, due_date
)
VALUES ($1, $2, $3, $4, $5)
RETURNING task_id, sales_process_id, title, description, status, due_date, created_at, updated_at;

-- name: GetTaskByID :one
SELECT task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE task_id = $1;

-- name: ListTasksBySalesProcess :many
SELECT task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE sales_process_id = $1
ORDER BY due_date ASC, created_at ASC
LIMIT $2 OFFSET $3;

-- name: ListTasksBySalesProcessAndStatus :many
SELECT task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE sales_process_id = $1 AND status = $2
ORDER BY due_date ASC, created_at ASC
LIMIT $3 OFFSET $4;

-- name: UpdateTask :one
UPDATE tasks
SET title = $2,
    description = $3,
    status = $4,
    due_date = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE task_id = $1
RETURNING task_id, sales_process_id, title, description, status, due_date, created_at, updated_at;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE task_id = $1
RETURNING task_id, sales_process_id, title, description, status, due_date, created_at, updated_at;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE task_id = $1;