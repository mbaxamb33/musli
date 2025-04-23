-- name: UploadProjectFile :one
INSERT INTO project_files (
    project_id, file_name, file_type, file_path, file_size
)
VALUES ($1, $2, $3, $4, $5)
RETURNING file_id, project_id, file_name, file_type, file_path, uploaded_at, file_size;

-- name: GetFileByID :one
SELECT file_id, project_id, file_name, file_type, file_path, uploaded_at, file_size
FROM project_files
WHERE file_id = $1;

-- name: ListFilesForProject :many
SELECT file_id, project_id, file_name, file_type, file_path, uploaded_at, file_size
FROM project_files
WHERE project_id = $1
ORDER BY uploaded_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteProjectFile :exec
DELETE FROM project_files
WHERE file_id = $1;

-- name: UpdateProjectFileInfo :one
UPDATE project_files
SET file_name = $2,
    file_type = $3,
    file_path = $4,
    file_size = $5
WHERE file_id = $1
RETURNING file_id, project_id, file_name, file_type, file_path, uploaded_at, file_size;
