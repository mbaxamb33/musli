-- name: CreateCompany :one
INSERT INTO companies (
    user_id, company_name, industry, website, address, description
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING company_id, user_id, company_name, industry, website, address, description, created_at;

-- name: GetCompanyByID :one
SELECT company_id, user_id, company_name, industry, website, address, description, created_at
FROM companies
WHERE company_id = $1;

-- name: GetCompaniesByUserID :many
SELECT company_id, user_id, company_name, industry, website, address, description, created_at
FROM companies
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCompanyByName :one
SELECT company_id, user_id, company_name, industry, website, address, description, created_at
FROM companies
WHERE company_name = $1;

-- name: ListCompanies :many
SELECT company_id, user_id, company_name, industry, website, address, description, created_at
FROM companies
ORDER BY company_name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateCompany :one
UPDATE companies
SET company_name = $2,
    industry = $3,
    website = $4,
    address = $5,
    description = $6
WHERE company_id = $1
RETURNING company_id, user_id, company_name, industry, website, address, description, created_at;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE company_id = $1;