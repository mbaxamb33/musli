-- name: CreateContact :one
INSERT INTO contacts (
  first_name,
  last_name,
  title,
  email,
  phone
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts
WHERE contact_id = $1 LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts
ORDER BY last_name, first_name
LIMIT $1 OFFSET $2;

-- name: UpdateContact :one
UPDATE contacts
SET
  first_name = COALESCE($2, first_name),
  last_name = COALESCE($3, last_name),
  title = COALESCE($4, title),
  email = COALESCE($5, email),
  phone = COALESCE($6, phone),
  updated_at = CURRENT_TIMESTAMP
WHERE contact_id = $1
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE contact_id = $1;

-- name: SearchContacts :many
SELECT * FROM contacts
WHERE 
  first_name ILIKE '%' || $1 || '%' OR 
  last_name ILIKE '%' || $1 || '%' OR
  email ILIKE '%' || $1 || '%'
ORDER BY last_name, first_name
LIMIT $2 OFFSET $3;

-- name: GetContactsByCompany :many
SELECT 
  c.*,
  cc.is_primary
FROM contacts c
JOIN company_contacts cc ON c.contact_id = cc.contact_id
WHERE cc.company_id = $1
ORDER BY cc.is_primary DESC, c.last_name, c.first_name
LIMIT $2 OFFSET $3;

-- name: GetPrimaryContactForCompany :one
SELECT c.* 
FROM contacts c
JOIN company_contacts cc ON c.contact_id = cc.contact_id
WHERE cc.company_id = $1 AND cc.is_primary = TRUE
LIMIT 1;