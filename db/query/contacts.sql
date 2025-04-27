-- name: CreateContact :one
INSERT INTO contacts (
    company_id, first_name, last_name, position, email, phone, notes
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING contact_id, company_id, first_name, last_name, position, email, phone, notes, created_at;

-- name: GetContactByID :one
SELECT contact_id, company_id, first_name, last_name, position, email, phone, notes, created_at
FROM contacts
WHERE contact_id = $1;

-- name: ListContactsByCompanyID :many
SELECT contact_id, company_id, first_name, last_name, position, email, phone, notes, created_at
FROM contacts
WHERE company_id = $1
ORDER BY last_name, first_name ASC
LIMIT $2 OFFSET $3;

-- name: SearchContactsByName :many
SELECT contact_id, company_id, first_name, last_name, position, email, phone, notes, created_at
FROM contacts
WHERE (first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%')
ORDER BY last_name, first_name ASC
LIMIT $2 OFFSET $3;

-- name: SearchContactsByCompanyAndName :many
SELECT contact_id, company_id, first_name, last_name, position, email, phone, notes, created_at
FROM contacts
WHERE company_id = $1 AND (first_name ILIKE '%' || $2 || '%' OR last_name ILIKE '%' || $2 || '%')
ORDER BY last_name, first_name ASC
LIMIT $3 OFFSET $4;

-- name: UpdateContact :one
UPDATE contacts
SET first_name = $2,
    last_name = $3,
    position = $4,
    email = $5,
    phone = $6,
    notes = $7
WHERE contact_id = $1
RETURNING contact_id, company_id, first_name, last_name, position, email, phone, notes, created_at;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE contact_id = $1;