-- name: CreateContact :one
INSERT INTO contacts (
    first_name, last_name, email, phone, linkedin_profile,
    job_title, company_id, location, bio
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING contact_id, first_name, last_name, email, phone, linkedin_profile,
          job_title, company_id, location, bio, scrape_timestamp;

-- name: GetContactByID :one
SELECT contact_id, first_name, last_name, email, phone, linkedin_profile,
       job_title, company_id, location, bio, scrape_timestamp
FROM contacts
WHERE contact_id = $1;

-- name: ListContactsByCompany :many
SELECT contact_id, first_name, last_name, email, phone, linkedin_profile,
       job_title, company_id, location, bio, scrape_timestamp
FROM contacts
WHERE company_id = $1
ORDER BY scrape_timestamp DESC
LIMIT $2 OFFSET $3;

-- name: SearchContactsByName :many
SELECT contact_id, first_name, last_name, email, phone, linkedin_profile,
       job_title, company_id, location, bio, scrape_timestamp
FROM contacts
WHERE LOWER(first_name) LIKE LOWER($1) OR LOWER(last_name) LIKE LOWER($1)
ORDER BY scrape_timestamp DESC
LIMIT $2 OFFSET $3;

-- name: UpdateContact :one
UPDATE contacts
SET first_name = $2,
    last_name = $3,
    email = $4,
    phone = $5,
    linkedin_profile = $6,
    job_title = $7,
    company_id = $8,
    location = $9,
    bio = $10,
    scrape_timestamp = CURRENT_TIMESTAMP
WHERE contact_id = $1
RETURNING contact_id, first_name, last_name, email, phone, linkedin_profile,
          job_title, company_id, location, bio, scrape_timestamp;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE contact_id = $1;
