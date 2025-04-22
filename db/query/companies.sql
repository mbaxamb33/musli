-- name: CreateCompany :one
INSERT INTO companies (
  name,
  industry,
  size,
  location,
  website,
  description
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetCompany :one
SELECT * FROM companies
WHERE company_id = $1 LIMIT 1;

-- name: ListCompanies :many
SELECT * FROM companies
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListCompaniesByIndustry :many
SELECT * FROM companies
WHERE industry = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateCompany :one
UPDATE companies
SET
  name = COALESCE($2, name),
  industry = COALESCE($3, industry),
  size = COALESCE($4, size),
  location = COALESCE($5, location),
  website = COALESCE($6, website),
  description = COALESCE($7, description),
  updated_at = CURRENT_TIMESTAMP
WHERE company_id = $1
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE company_id = $1;

-- name: SearchCompanies :many
SELECT * FROM companies
WHERE 
  name ILIKE '%' || $1 || '%' OR 
  industry ILIKE '%' || $1 || '%' OR 
  location ILIKE '%' || $1 || '%' OR
  description ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: GetCompanyWithContactCount :many
SELECT 
  c.*,
  COUNT(cc.contact_id) AS contact_count
FROM companies c
LEFT JOIN company_contacts cc ON c.company_id = cc.company_id
GROUP BY c.company_id
ORDER BY c.name
LIMIT $1 OFFSET $2;

-- name: GetCompaniesForProject :many
SELECT 
  c.*,
  pc.matching_score,
  pc.status AS relationship_status,
  strategy.name AS approach_strategy
FROM companies c
JOIN project_companies pc ON c.company_id = pc.company_id
LEFT JOIN approach_strategies strategy ON pc.approach_strategy_id = strategy.strategy_id
WHERE pc.project_id = $1
ORDER BY pc.matching_score DESC
LIMIT $2 OFFSET $3;