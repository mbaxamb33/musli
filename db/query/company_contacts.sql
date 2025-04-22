-- name: AddContactToCompany :one
INSERT INTO company_contacts (
  company_id,
  contact_id,
  is_primary
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCompanyContact :one
SELECT * FROM company_contacts
WHERE company_id = $1 AND contact_id = $2
LIMIT 1;

-- name: GetCompanyContactById :one
SELECT * FROM company_contacts
WHERE company_contact_id = $1
LIMIT 1;

-- name: ListCompanyContacts :many
SELECT 
  cc.*,
  c.first_name,
  c.last_name,
  c.title,
  c.email,
  c.phone
FROM company_contacts cc
JOIN contacts c ON cc.contact_id = c.contact_id
WHERE cc.company_id = $1
ORDER BY cc.is_primary DESC, c.last_name, c.first_name
LIMIT $2 OFFSET $3;

-- name: UpdateCompanyContact :one
UPDATE company_contacts
SET
  is_primary = $3,
  updated_at = CURRENT_TIMESTAMP
WHERE company_id = $1 AND contact_id = $2
RETURNING *;

-- name: SetPrimaryContact :exec
-- First, set all contacts for the company to not primary
UPDATE company_contacts
SET 
  is_primary = FALSE,
  updated_at = CURRENT_TIMESTAMP
WHERE company_id = $1;

-- name: SetContactAsPrimary :one
-- Then set the specific contact as primary
UPDATE company_contacts
SET 
  is_primary = TRUE,
  updated_at = CURRENT_TIMESTAMP
WHERE company_id = $1 AND contact_id = $2
RETURNING *;

-- name: RemoveContactFromCompany :exec
DELETE FROM company_contacts
WHERE company_id = $1 AND contact_id = $2;

-- name: RemoveCompanyContactById :exec
DELETE FROM company_contacts
WHERE company_contact_id = $1;

-- name: GetCompaniesForContact :many
SELECT 
  c.*,
  cc.is_primary
FROM companies c
JOIN company_contacts cc ON c.company_id = cc.company_id
WHERE cc.contact_id = $1
ORDER BY c.name
LIMIT $2 OFFSET $3;

-- name: CreateContactAndLinkToCompany :one
WITH new_contact AS (
  INSERT INTO contacts (
    first_name,
    last_name,
    title,
    email,
    phone
  ) VALUES (
    $1, $2, $3, $4, $5
  ) RETURNING *
)
INSERT INTO company_contacts (
  company_id,
  contact_id,
  is_primary
) VALUES (
  $6,
  (SELECT contact_id FROM new_contact),
  $7
) RETURNING company_contact_id;