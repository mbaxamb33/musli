-- name: LinkProjectToSalesProcess :exec
INSERT INTO sales_process_projects (
    sales_process_id, project_id
)
VALUES ($1, $2);

-- name: UnlinkProjectFromSalesProcess :exec
DELETE FROM sales_process_projects
WHERE sales_process_id = $1 AND project_id = $2;

-- name: GetProjectsForSalesProcess :many
SELECT p.project_id, p.cognito_sub, p.project_name, p.main_idea, p.created_at, p.updated_at
FROM projects p
JOIN sales_process_projects spp ON p.project_id = spp.project_id
WHERE spp.sales_process_id = $1
ORDER BY p.project_name ASC
LIMIT $2 OFFSET $3;

-- name: GetSalesProcessesForProject :many
SELECT sp.sales_process_id, sp.cognito_sub, sp.contact_id, sp.overall_matching_score, sp.status, sp.created_at, sp.updated_at
FROM sales_processes sp
JOIN sales_process_projects spp ON sp.sales_process_id = spp.sales_process_id
WHERE spp.project_id = $1
ORDER BY sp.created_at DESC
LIMIT $2 OFFSET $3;