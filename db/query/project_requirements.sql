-- =============================================================================
-- PROJECT REQUIREMENTS QUERIES
-- =============================================================================

-- name: CreateProjectRequirements :one
INSERT INTO project_requirements (
    brief_id, project_scope, success_criteria, implementation_timeline, resource_allocation,
    project_team_structure, change_management_approach, training_requirements, rollout_strategy,
    pilot_phase_design, risk_mitigation_plan, communication_plan, stakeholder_engagement,
    performance_metrics, governance_structure, escalation_procedures
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, project_scope, success_criteria, implementation_timeline, resource_allocation,
    project_team_structure, change_management_approach, training_requirements, rollout_strategy,
    pilot_phase_design, risk_mitigation_plan, communication_plan, stakeholder_engagement,
    performance_metrics, governance_structure, escalation_procedures, created_at, updated_at;

-- name: GetProjectRequirementsByBriefID :one
SELECT id, brief_id, project_scope, success_criteria, implementation_timeline, resource_allocation,
    project_team_structure, change_management_approach, training_requirements, rollout_strategy,
    pilot_phase_design, risk_mitigation_plan, communication_plan, stakeholder_engagement,
    performance_metrics, governance_structure, escalation_procedures, created_at, updated_at
FROM project_requirements
WHERE brief_id = $1;

-- name: UpdateProjectRequirements :one
UPDATE project_requirements
SET project_scope = $2,
    success_criteria = $3,
    implementation_timeline = $4,
    resource_allocation = $5,
    project_team_structure = $6,
    change_management_approach = $7,
    training_requirements = $8,
    rollout_strategy = $9,
    pilot_phase_design = $10,
    risk_mitigation_plan = $11,
    communication_plan = $12,
    stakeholder_engagement = $13,
    performance_metrics = $14,
    governance_structure = $15,
    escalation_procedures = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, project_scope, success_criteria, implementation_timeline, resource_allocation,
    project_team_structure, change_management_approach, training_requirements, rollout_strategy,
    pilot_phase_design, risk_mitigation_plan, communication_plan, stakeholder_engagement,
    performance_metrics, governance_structure, escalation_procedures, created_at, updated_at;

-- name: DeleteProjectRequirements :exec
DELETE FROM project_requirements
WHERE brief_id = $1;