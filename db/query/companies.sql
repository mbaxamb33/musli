-- name: CreateCompany :one
INSERT INTO companies (
    name, website, industry, description,
    headquarters_location, founded_year, is_public, ticker_symbol
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING company_id, name, website, industry, description,
          headquarters_location, founded_year, is_public, ticker_symbol, scrape_timestamp;

-- name: GetCompanyByID :one
SELECT company_id, name, website, industry, description,
       headquarters_location, founded_year, is_public, ticker_symbol, scrape_timestamp
FROM companies
WHERE company_id = $1;

-- name: GetCompanyByName :one
SELECT company_id, name, website, industry, description,
       headquarters_location, founded_year, is_public, ticker_symbol, scrape_timestamp
FROM companies
WHERE name = $1;

-- name: ListCompanies :many
SELECT company_id, name, website, industry, description,
       headquarters_location, founded_year, is_public, ticker_symbol, scrape_timestamp
FROM companies
ORDER BY scrape_timestamp DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCompany :one
UPDATE companies
SET name = $2,
    website = $3,
    industry = $4,
    description = $5,
    headquarters_location = $6,
    founded_year = $7,
    is_public = $8,
    ticker_symbol = $9,
    scrape_timestamp = CURRENT_TIMESTAMP
WHERE company_id = $1
RETURNING company_id, name, website, industry, description,
          headquarters_location, founded_year, is_public, ticker_symbol, scrape_timestamp;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE company_id = $1;
