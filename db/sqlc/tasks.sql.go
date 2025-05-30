// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: tasks.sql

package db

import (
	"context"
	"database/sql"
)

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (
    sales_process_id, title, description, status, due_date
)
VALUES ($1, $2, $3, $4, $5)
RETURNING task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
`

type CreateTaskParams struct {
	SalesProcessID int32          `json:"sales_process_id"`
	Title          string         `json:"title"`
	Description    sql.NullString `json:"description"`
	Status         TaskStatus     `json:"status"`
	DueDate        sql.NullTime   `json:"due_date"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, createTask,
		arg.SalesProcessID,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.DueDate,
	)
	var i Task
	err := row.Scan(
		&i.TaskID,
		&i.SalesProcessID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.DueDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE task_id = $1
`

func (q *Queries) DeleteTask(ctx context.Context, taskID int32) error {
	_, err := q.db.ExecContext(ctx, deleteTask, taskID)
	return err
}

const getTaskByID = `-- name: GetTaskByID :one
SELECT task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE task_id = $1
`

func (q *Queries) GetTaskByID(ctx context.Context, taskID int32) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskByID, taskID)
	var i Task
	err := row.Scan(
		&i.TaskID,
		&i.SalesProcessID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.DueDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTasksBySalesProcess = `-- name: ListTasksBySalesProcess :many
SELECT task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE sales_process_id = $1
ORDER BY due_date ASC, created_at ASC
LIMIT $2 OFFSET $3
`

type ListTasksBySalesProcessParams struct {
	SalesProcessID int32 `json:"sales_process_id"`
	Limit          int32 `json:"limit"`
	Offset         int32 `json:"offset"`
}

func (q *Queries) ListTasksBySalesProcess(ctx context.Context, arg ListTasksBySalesProcessParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasksBySalesProcess, arg.SalesProcessID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.TaskID,
			&i.SalesProcessID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTasksBySalesProcessAndStatus = `-- name: ListTasksBySalesProcessAndStatus :many
SELECT task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
FROM tasks
WHERE sales_process_id = $1 AND status = $2
ORDER BY due_date ASC, created_at ASC
LIMIT $3 OFFSET $4
`

type ListTasksBySalesProcessAndStatusParams struct {
	SalesProcessID int32      `json:"sales_process_id"`
	Status         TaskStatus `json:"status"`
	Limit          int32      `json:"limit"`
	Offset         int32      `json:"offset"`
}

func (q *Queries) ListTasksBySalesProcessAndStatus(ctx context.Context, arg ListTasksBySalesProcessAndStatusParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasksBySalesProcessAndStatus,
		arg.SalesProcessID,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.TaskID,
			&i.SalesProcessID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :one
UPDATE tasks
SET title = $2,
    description = $3,
    status = $4,
    due_date = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE task_id = $1
RETURNING task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
`

type UpdateTaskParams struct {
	TaskID      int32          `json:"task_id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Status      TaskStatus     `json:"status"`
	DueDate     sql.NullTime   `json:"due_date"`
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, updateTask,
		arg.TaskID,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.DueDate,
	)
	var i Task
	err := row.Scan(
		&i.TaskID,
		&i.SalesProcessID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.DueDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateTaskStatus = `-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE task_id = $1
RETURNING task_id, sales_process_id, title, description, status, due_date, created_at, updated_at
`

type UpdateTaskStatusParams struct {
	TaskID int32      `json:"task_id"`
	Status TaskStatus `json:"status"`
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, updateTaskStatus, arg.TaskID, arg.Status)
	var i Task
	err := row.Scan(
		&i.TaskID,
		&i.SalesProcessID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.DueDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
