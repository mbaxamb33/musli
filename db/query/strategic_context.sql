-- =============================================================================
-- STRATEGIC CONTEXT QUERIES
-- =============================================================================

-- name: CreateStrategicContext :one
INSERT INTO strategic_context (
    brief_id, business_strategy, strategic_initiatives, quarterly_priorities, annual_goals,
    transformation_agenda, digital_maturity, innovation_focus, operational_challenges,
    cost_reduction_pressures, revenue_growth_targets, efficiency_mandates, compliance_drivers,
    risk_management_priorities, sustainability_goals, technology_roadmap
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, business_strategy, strategic_initiatives, quarterly_priorities, annual_goals,
    transformation_agenda, digital_maturity, innovation_focus, operational_challenges,
    cost_reduction_pressures, revenue_growth_targets, efficiency_mandates, compliance_drivers,
    risk_management_priorities, sustainability_goals, technology_roadmap, created_at, updated_at;

-- name: GetStrategicContextByBriefID :one
SELECT id, brief_id, business_strategy, strategic_initiatives, quarterly_priorities, annual_goals,
    transformation_agenda, digital_maturity, innovation_focus, operational_challenges,
    cost_reduction_pressures, revenue_growth_targets, efficiency_mandates, compliance_drivers,
    risk_management_priorities, sustainability_goals, technology_roadmap, created_at, updated_at
FROM strategic_context
WHERE brief_id = $1;

-- name: UpdateStrategicContext :one
UPDATE strategic_context
SET business_strategy = $2,
    strategic_initiatives = $3,
    quarterly_priorities = $4,
    annual_goals = $5,
    transformation_agenda = $6,
    digital_maturity = $7,
    innovation_focus = $8,
    operational_challenges = $9,
    cost_reduction_pressures = $10,
    revenue_growth_targets = $11,
    efficiency_mandates = $12,
    compliance_drivers = $13,
    risk_management_priorities = $14,
    sustainability_goals = $15,
    technology_roadmap = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, business_strategy, strategic_initiatives, quarterly_priorities, annual_goals,
    transformation_agenda, digital_maturity, innovation_focus, operational_challenges,
    cost_reduction_pressures, revenue_growth_targets, efficiency_mandates, compliance_drivers,
    risk_management_priorities, sustainability_goals, technology_roadmap, created_at, updated_at;

-- name: DeleteStrategicContext :exec
DELETE FROM strategic_context
WHERE brief_id = $1;