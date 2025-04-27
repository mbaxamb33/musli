-- name: CreateAnalysisContactInfo :one
INSERT INTO analysis_contact_info (
    analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information, updated_at;

-- name: GetAnalysisContactInfo :one
SELECT analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information, updated_at
FROM analysis_contact_info
WHERE analysis_id = $1;

-- name: UpdateAnalysisContactInfo :one
UPDATE analysis_contact_info
SET problems = $2,
    needs = $3,
    urgency = $4,
    priorities = $5,
    decision_process = $6,
    budget = $7,
    resources = $8,
    relevant_information = $9,
    updated_at = CURRENT_TIMESTAMP
WHERE analysis_id = $1
RETURNING analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information, updated_at;

-- name: DeleteAnalysisContactInfo :exec
DELETE FROM analysis_contact_info
WHERE analysis_id = $1;