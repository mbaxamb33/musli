-- =============================================================================
-- SALES PROCESS TRACKING QUERIES
-- =============================================================================

-- name: CreateSalesProcessTracking :one
INSERT INTO sales_process_tracking (
    brief_id, lead_source, opportunity_stage, probability_percentage, weighted_value,
    next_action_required, key_milestones, sales_velocity, deal_momentum, competitive_position,
    win_probability_factors
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, brief_id, lead_source, opportunity_stage, probability_percentage, weighted_value,
    next_action_required, key_milestones, sales_velocity, deal_momentum, competitive_position,
    win_probability_factors, created_at, updated_at;

-- name: GetSalesProcessTrackingByBriefID :one
SELECT id, brief_id, lead_source, opportunity_stage, probability_percentage, weighted_value,
    next_action_required, key_milestones, sales_velocity, deal_momentum, competitive_position,
    win_probability_factors, created_at, updated_at
FROM sales_process_tracking
WHERE brief_id = $1;

-- name: UpdateSalesProcessTracking :one
UPDATE sales_process_tracking
SET lead_source = $2,
    opportunity_stage = $3,
    probability_percentage = $4,
    weighted_value = $5,
    next_action_required = $6,
    key_milestones = $7,
    sales_velocity = $8,
    deal_momentum = $9,
    competitive_position = $10,
    win_probability_factors = $11,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, lead_source, opportunity_stage, probability_percentage, weighted_value,
    next_action_required, key_milestones, sales_velocity, deal_momentum, competitive_position,
    win_probability_factors, created_at, updated_at;

-- name: DeleteSalesProcessTracking :exec
DELETE FROM sales_process_tracking
WHERE brief_id = $1;