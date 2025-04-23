-- name: CreateWebsitePage :one
INSERT INTO website_pages (
    website_id, url, path, title, parent_page_id, depth, extract_status, datasource_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetWebsitePageByID :one
SELECT * FROM website_pages
WHERE page_id = $1;

-- name: GetWebsitePageByURL :one
SELECT * FROM website_pages
WHERE website_id = $1 AND url = $2;

-- name: ListWebsitePagesByWebsiteID :many
SELECT * FROM website_pages
WHERE website_id = $1
ORDER BY path ASC
LIMIT $2 OFFSET $3;

-- name: ListWebsitePageTree :many
WITH RECURSIVE page_tree AS (
    -- Base case: select root pages (parent_page_id IS NULL)
    SELECT 
        wp.page_id, 
        wp.website_id, 
        wp.url, 
        wp.path, 
        wp.title, 
        wp.parent_page_id, 
        wp.depth, 
        wp.last_extracted_at, 
        wp.extract_status, 
        wp.datasource_id,
        wp.title AS page_path
    FROM website_pages wp
    WHERE wp.website_id = $1 AND wp.parent_page_id IS NULL
    
    UNION ALL
    
    -- Recursive case: join with pages that have a parent in our tree
    SELECT 
        wp.page_id, 
        wp.website_id, 
        wp.url, 
        wp.path, 
        wp.title, 
        wp.parent_page_id, 
        wp.depth, 
        wp.last_extracted_at, 
        wp.extract_status, 
        wp.datasource_id,
        pt.page_path || ' > ' || wp.title AS page_path
    FROM website_pages wp
    JOIN page_tree pt ON wp.parent_page_id = pt.page_id
)
SELECT * FROM page_tree
ORDER BY page_path
LIMIT $2 OFFSET $3;

-- name: UpdateWebsitePage :one
UPDATE website_pages
SET title = $3,
    parent_page_id = $4,
    last_extracted_at = $5,
    extract_status = $6,
    datasource_id = $7
WHERE website_id = $1 AND url = $2
RETURNING *;

-- name: UpdateExtractStatus :one
UPDATE website_pages
SET last_extracted_at = CURRENT_TIMESTAMP,
    extract_status = $3
WHERE website_id = $1 AND page_id = $2
RETURNING *;

-- name: DeleteWebsitePage :exec
DELETE FROM website_pages
WHERE page_id = $1;

-- name: GetPagesForExtraction :many
SELECT * FROM website_pages
WHERE extract_status = 'pending' OR 
      (extract_status = 'completed' AND last_extracted_at < NOW() - INTERVAL '30 days')
ORDER BY last_extracted_at ASC NULLS FIRST
LIMIT $1;