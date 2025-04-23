-- name: AssociateContactWithProject :one
INSERT INTO project_contacts (
    project_id, contact_id, association_notes
)
VALUES ($1, $2, $3)
RETURNING project_id, contact_id, association_notes;

-- name: GetProjectContactAssociation :one
SELECT project_id, contact_id, association_notes
FROM project_contacts
WHERE project_id = $1 AND contact_id = $2;

-- name: ListContactsForProject :many
SELECT pc.project_id, pc.contact_id, pc.association_notes,
       c.first_name, c.last_name, c.email, c.job_title
FROM project_contacts pc
JOIN contacts c ON pc.contact_id = c.contact_id
WHERE pc.project_id = $1
ORDER BY c.last_name ASC, c.first_name ASC
LIMIT $2 OFFSET $3;

-- name: ListProjectsForContact :many
SELECT pc.project_id, pc.contact_id, pc.association_notes,
       p.project_name, p.description
FROM project_contacts pc
JOIN projects p ON pc.project_id = p.project_id
WHERE pc.contact_id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProjectContactAssociation :one
UPDATE project_contacts
SET association_notes = $3
WHERE project_id = $1 AND contact_id = $2
RETURNING project_id, contact_id, association_notes;

-- name: RemoveProjectContactAssociation :exec
DELETE FROM project_contacts
WHERE project_id = $1 AND contact_id = $2;
