-- name: CreateMeeting :one
INSERT INTO meetings (
    sales_process_id, contact_id, task_id, meeting_time, meeting_place, notes
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING meeting_id, sales_process_id, contact_id, task_id, meeting_time, meeting_place, notes, created_at;

-- name: GetMeetingByID :one
SELECT meeting_id, sales_process_id, contact_id, task_id, meeting_time, meeting_place, notes, created_at
FROM meetings
WHERE meeting_id = $1;

-- name: ListMeetingsBySalesProcess :many
SELECT meeting_id, sales_process_id, contact_id, task_id, meeting_time, meeting_place, notes, created_at
FROM meetings
WHERE sales_process_id = $1
ORDER BY meeting_time DESC
LIMIT $2 OFFSET $3;

-- name: ListMeetingsByContact :many
SELECT meeting_id, sales_process_id, contact_id, task_id, meeting_time, meeting_place, notes, created_at
FROM meetings
WHERE contact_id = $1
ORDER BY meeting_time DESC
LIMIT $2 OFFSET $3;

-- name: ListMeetingsByTask :many
SELECT meeting_id, sales_process_id, contact_id, task_id, meeting_time, meeting_place, notes, created_at
FROM meetings
WHERE task_id = $1
ORDER BY meeting_time DESC;

-- name: UpdateMeeting :one
UPDATE meetings
SET meeting_time = $2,
    meeting_place = $3,
    notes = $4
WHERE meeting_id = $1
RETURNING meeting_id, sales_process_id, contact_id, task_id, meeting_time, meeting_place, notes, created_at;

-- name: DeleteMeeting :exec
DELETE FROM meetings
WHERE meeting_id = $1;