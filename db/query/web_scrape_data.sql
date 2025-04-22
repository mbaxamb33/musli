-- name: CreateWebScrapeData :one
INSERT INTO web_scrape_data (
  company_id,
  source_url,
  data_type,
  content,
  is_processed
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetWebScrapeData :one
SELECT * FROM web_scrape_data
WHERE scrape_id = $1 LIMIT 1;

-- name: ListWebScrapeDataForCompany :many
SELECT * FROM web_scrape_data
WHERE company_id = $1
ORDER BY scrape_date DESC
LIMIT $2 OFFSET $3;

-- name: ListWebScrapeDataByType :many
SELECT * FROM web_scrape_data
WHERE company_id = $1 AND data_type = $2
ORDER BY scrape_date DESC
LIMIT $3 OFFSET $4;

-- name: ListUnprocessedWebScrapeData :many
SELECT * FROM web_scrape_data
WHERE is_processed = FALSE
ORDER BY scrape_date
LIMIT $1 OFFSET $2;

-- name: UpdateWebScrapeData :one
UPDATE web_scrape_data
SET
  source_url = COALESCE($2, source_url),
  data_type = COALESCE($3, data_type),
  content = COALESCE($4, content),
  is_processed = COALESCE($5, is_processed)
WHERE scrape_id = $1
RETURNING *;

-- name: MarkWebScrapeDataAsProcessed :one
UPDATE web_scrape_data
SET 
  is_processed = TRUE
WHERE scrape_id = $1
RETURNING *;

-- name: DeleteWebScrapeData :exec
DELETE FROM web_scrape_data
WHERE scrape_id = $1;

-- name: DeleteAllWebScrapeDataForCompany :exec
DELETE FROM web_scrape_data
WHERE company_id = $1;

-- name: GetWebScrapeDataTypes :many
SELECT DISTINCT data_type
FROM web_scrape_data
ORDER BY data_type;

-- name: GetWebScrapeDataCountByType :many
SELECT 
  data_type, 
  COUNT(*) as count
FROM web_scrape_data
WHERE company_id = $1
GROUP BY data_type
ORDER BY count DESC;

-- name: GetLatestWebScrapeByType :one
SELECT * FROM web_scrape_data
WHERE company_id = $1 AND data_type = $2
ORDER BY scrape_date DESC
LIMIT 1;