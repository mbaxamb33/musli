-- name: CreateCompanyWebsite :one
INSERT INTO company_websites (
    company_id, base_url, site_title, scrape_frequency_days, is_active, datasource_id
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCompanyWebsiteByID :one
SELECT * FROM company_websites
WHERE website_id = $1;

-- name: GetCompanyWebsitesByCompanyID :many
SELECT * FROM company_websites
WHERE company_id = $1
ORDER BY last_scraped_at DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: UpdateCompanyWebsite :one
UPDATE company_websites
SET base_url = $2,
    site_title = $3,
    last_scraped_at = $4,
    scrape_frequency_days = $5,
    is_active = $6,
    datasource_id = $7
WHERE website_id = $1
RETURNING *;

-- name: UpdateLastScrapedAt :one
UPDATE company_websites
SET last_scraped_at = CURRENT_TIMESTAMP
WHERE website_id = $1
RETURNING *;

-- name: DeleteCompanyWebsite :exec
DELETE FROM company_websites
WHERE website_id = $1;

-- name: ListCompanyWebsitesForScraping :many
SELECT * FROM company_websites
WHERE is_active = true AND 
      (last_scraped_at IS NULL OR 
       last_scraped_at < NOW() - (scrape_frequency_days * INTERVAL '1 day'))
LIMIT $1;