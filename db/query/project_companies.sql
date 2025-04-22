-- name: LinkCompanyToProject :one
INSERT INTO project_companies (
  project_id,
  company_id,
  matching_score,
  approach_strategy_id,
  status
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetProjectCompany :one
SELECT * FROM project_companies
WHERE project_id = $1 AND company_id = $2
LIMIT 1;

-- name: GetProjectCompanyById :one
SELECT * FROM project_companies
WHERE project_company_id = $1
LIMIT 1;

-- name: ListProjectCompanies :many
SELECT
  pc.*,
  c.name AS company_name,
  c.industry,
  c.size,
  strategy.name AS approach_strategy
FROM project_companies pc
JOIN companies c ON pc.company_id = c.company_id
LEFT JOIN approach_strategies strategy ON pc.approach_strategy_id = strategy.strategy_id
WHERE pc.project_id = $1
ORDER BY pc.matching_score DESC
LIMIT $2 OFFSET $3;

-- name: ListProjectCompaniesByStatus :many
SELECT
  pc.*,
  c.name AS company_name,
  c.industry,
  c.size,
  strategy.name AS approach_strategy
FROM project_companies pc
JOIN companies c ON pc.company_id = c.company_id
LEFT JOIN approach_strategies strategy ON pc.approach_strategy_id = strategy.strategy_id
WHERE pc.project_id = $1 AND pc.status = $2
ORDER BY pc.matching_score DESC
LIMIT $3 OFFSET $4;

-- name: UpdateProjectCompany :one
UPDATE project_companies
SET
  matching_score = COALESCE($3, matching_score),
  approach_strategy_id = COALESCE($4, approach_strategy_id),
  status = COALESCE($5, status),
  updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1 AND company_id = $2
RETURNING *;

-- name: UpdateProjectCompanyStatus :one
UPDATE project_companies
SET 
  status = $3,
  updated_at = CURRENT_TIMESTAMP
WHERE project_id = $1 AND company_id = $2
RETURNING *;

-- name: RemoveCompanyFromProject :exec
DELETE FROM project_companies
WHERE project_id = $1 AND company_id = $2;

-- name: GetProjectsForCompany :many
SELECT 
  p.*,
  pc.matching_score,
  pc.status AS relationship_status,
  u.username AS owner_username
FROM projects p
JOIN project_companies pc ON p.project_id = pc.project_id
JOIN users u ON p.user_id = u.user_id
WHERE pc.company_id = $1
ORDER BY pc.matching_score DESC
LIMIT $2 OFFSET $3;

-- name: GetTopMatchingCompaniesForProject :many
SELECT
  c.*,
  pc.matching_score,
  pc.status
FROM companies c
JOIN project_companies pc ON c.company_id = pc.company_id
WHERE pc.project_id = $1
ORDER BY pc.matching_score DESC
LIMIT $2;

-- name: GetProjectCompanyMatches :many
SELECT 
  pc.*,
  c.name AS company_name,
  strategy.name AS strategy_name,
  (
    SELECT COUNT(*) 
    FROM matching_scores_detail msd 
    WHERE msd.project_company_id = pc.project_company_id
  ) AS criteria_count
FROM project_companies pc
JOIN companies c ON pc.company_id = c.company_id
LEFT JOIN approach_strategies strategy ON pc.approach_strategy_id = strategy.strategy_id
WHERE pc.project_id = $1
ORDER BY pc.matching_score DESC
LIMIT $2 OFFSET $3;