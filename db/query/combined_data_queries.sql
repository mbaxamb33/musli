-- name: GetCompanyWithDatasources :many
SELECT c.company_id, c.company_name, c.industry, c.website, 
       d.datasource_id, d.source_type, d.link, d.file_name, d.created_at as datasource_created_at
FROM companies c
JOIN company_datasources cd ON c.company_id = cd.company_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
WHERE c.company_id = $1
ORDER BY d.created_at DESC;

-- name: GetContactWithDatasources :many
SELECT ct.contact_id, ct.first_name, ct.last_name, ct.position, ct.email,
       d.datasource_id, d.source_type, d.link, d.file_name, d.created_at as datasource_created_at
FROM contacts ct
JOIN contact_datasources cd ON ct.contact_id = cd.contact_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
WHERE ct.contact_id = $1
ORDER BY d.created_at DESC;

-- name: GetCompanyParagraphs :many
SELECT c.company_id, c.company_name, d.datasource_id, d.source_type, 
       p.paragraph_id, p.title, p.main_idea, p.content, p.created_at
FROM companies c
JOIN company_datasources cd ON c.company_id = cd.company_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
JOIN paragraphs p ON d.datasource_id = p.datasource_id
WHERE c.company_id = $1
ORDER BY d.created_at DESC, p.paragraph_id ASC
LIMIT $2 OFFSET $3;

-- name: GetContactParagraphs :many
SELECT ct.contact_id, ct.first_name, ct.last_name, d.datasource_id, d.source_type,
       p.paragraph_id, p.title, p.main_idea, p.content, p.created_at
FROM contacts ct
JOIN contact_datasources cd ON ct.contact_id = cd.contact_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
JOIN paragraphs p ON d.datasource_id = p.datasource_id
WHERE ct.contact_id = $1
ORDER BY d.created_at DESC, p.paragraph_id ASC
LIMIT $2 OFFSET $3;

-- name: SearchCompanyParagraphs :many
SELECT c.company_id, c.company_name, d.datasource_id, d.source_type, 
       p.paragraph_id, p.title, p.main_idea, p.content
FROM companies c
JOIN company_datasources cd ON c.company_id = cd.company_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
JOIN paragraphs p ON d.datasource_id = p.datasource_id
WHERE c.company_id = $1 AND (p.content ILIKE '%' || $2 || '%' OR p.main_idea ILIKE '%' || $2 || '%')
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: SearchContactParagraphs :many
SELECT ct.contact_id, ct.first_name, ct.last_name, d.datasource_id, d.source_type,
       p.paragraph_id, p.title, p.main_idea, p.content
FROM contacts ct
JOIN contact_datasources cd ON ct.contact_id = cd.contact_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
JOIN paragraphs p ON d.datasource_id = p.datasource_id
WHERE ct.contact_id = $1 AND (p.content ILIKE '%' || $2 || '%' OR p.main_idea ILIKE '%' || $2 || '%')
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetCompanyAllData :many
SELECT c.company_id, c.company_name, c.industry, c.website, c.description,
       d.datasource_id, d.source_type, d.link, d.file_name,
       p.paragraph_id, p.title, p.main_idea, p.content
FROM companies c
JOIN company_datasources cd ON c.company_id = cd.company_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
LEFT JOIN paragraphs p ON d.datasource_id = p.datasource_id
WHERE c.company_id = $1
ORDER BY d.created_at DESC, p.paragraph_id ASC;

-- name: GetContactAllData :many
SELECT ct.contact_id, ct.first_name, ct.last_name, ct.position, ct.email, ct.notes,
       d.datasource_id, d.source_type, d.link, d.file_name, 
       p.paragraph_id, p.title, p.main_idea, p.content
FROM contacts ct
JOIN contact_datasources cd ON ct.contact_id = cd.contact_id
JOIN datasources d ON cd.datasource_id = d.datasource_id
LEFT JOIN paragraphs p ON d.datasource_id = p.datasource_id
WHERE ct.contact_id = $1
ORDER BY d.created_at DESC, p.paragraph_id ASC;