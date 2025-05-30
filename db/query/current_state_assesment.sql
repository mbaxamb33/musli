-- =============================================================================
-- CURRENT STATE ASSESSMENT QUERIES
-- =============================================================================

-- name: CreateCurrentStateAssessment :one
INSERT INTO current_state_assessment (
    brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues,
    created_at, updated_at;

-- name: GetCurrentStateAssessmentByBriefID :one
SELECT id, brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues,
    created_at, updated_at
FROM current_state_assessment
WHERE brief_id = $1;

-- name: UpdateCurrentStateAssessment :one
UPDATE current_state_assessment
SET current_solution_provider = $2,
    current_solution_satisfaction = $3,
    specific_pain_points = $4,
    workaround_solutions = $5,
    cost_of_status_quo = $6,
    switching_barriers = $7,
    contract_end_dates = $8,
    renewal_timing = $9,
    vendor_relationship_health = $10,
    support_satisfaction = $11,
    functionality_gaps = $12,
    performance_issues = $13,
    scalability_constraints = $14,
    integration_challenges = $15,
    user_adoption_issues = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues,
    created_at, updated_at;

-- name: DeleteCurrentStateAssessment :exec
DELETE FROM current_state_assessment
WHERE brief_id = $1;