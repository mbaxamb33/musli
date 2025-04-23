-- name: AssociateCompanyWithProject :one
INSERT INTO project_companies (
    project_id, company_id, association_notes, matching_score, approach_strategy
)
VALUES ($1, $2, $3, $4, $5)
RETURNING project_id, company_id, association_notes, matching_score, approach_strategy;

-- name: GetProjectCompanyAssociation :one
SELECT project_id, company_id, association_notes, matching_score, approach_strategy
FROM project_companies
WHERE project_id = $1 AND company_id = $2;

-- name: ListCompaniesForProject :many
SELECT pc.project_id, pc.company_id, pc.association_notes, pc.matching_score, pc.approach_strategy,
       c.name, c.industry, c.website
FROM project_companies pc
JOIN companies c ON pc.company_id = c.company_id
WHERE pc.project_id = $1
ORDER BY pc.matching_score DESC
LIMIT $2 OFFSET $3;

-- name: ListProjectsForCompany :many
SELECT pc.project_id, pc.company_id, pc.association_notes, pc.matching_score, pc.approach_strategy,
       p.project_name, p.description
FROM project_companies pc
JOIN projects p ON pc.project_id = p.project_id
WHERE pc.company_id = $1
ORDER BY pc.matching_score DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProjectCompanyAssociation :one
UPDATE project_companies
SET association_notes = $3,
    matching_score = $4,
    approach_strategy = $5
WHERE project_id = $1 AND company_id = $2
RETURNING project_id, company_id, association_notes, matching_score, approach_strategy;

-- name: RemoveProjectCompanyAssociation :exec
DELETE FROM project_companies
WHERE project_id = $1 AND company_id = $2;
