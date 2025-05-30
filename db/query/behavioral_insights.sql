-- =============================================================================
-- BEHAVIORAL & PSYCHOLOGICAL INSIGHTS QUERIES
-- =============================================================================

-- name: CreateBehavioralInsights :one
INSERT INTO behavioral_insights (
    brief_id, decision_making_style, risk_aversion_level, change_adoption_patterns, innovation_appetite,
    consensus_building_approach, conflict_resolution_style, communication_patterns, trust_building_factors,
    credibility_requirements, relationship_preferences, meeting_effectiveness, follow_up_responsiveness,
    documentation_preferences, presentation_style_preferences, negotiation_approach
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, decision_making_style, risk_aversion_level, change_adoption_patterns, innovation_appetite,
    consensus_building_approach, conflict_resolution_style, communication_patterns, trust_building_factors,
    credibility_requirements, relationship_preferences, meeting_effectiveness, follow_up_responsiveness,
    documentation_preferences, presentation_style_preferences, negotiation_approach, created_at, updated_at;

-- name: GetBehavioralInsightsByBriefID :one
SELECT id, brief_id, decision_making_style, risk_aversion_level, change_adoption_patterns, innovation_appetite,
    consensus_building_approach, conflict_resolution_style, communication_patterns, trust_building_factors,
    credibility_requirements, relationship_preferences, meeting_effectiveness, follow_up_responsiveness,
    documentation_preferences, presentation_style_preferences, negotiation_approach, created_at, updated_at
FROM behavioral_insights
WHERE brief_id = $1;

-- name: UpdateBehavioralInsights :one
UPDATE behavioral_insights
SET decision_making_style = $2,
    risk_aversion_level = $3,
    change_adoption_patterns = $4,
    innovation_appetite = $5,
    consensus_building_approach = $6,
    conflict_resolution_style = $7,
    communication_patterns = $8,
    trust_building_factors = $9,
    credibility_requirements = $10,
    relationship_preferences = $11,
    meeting_effectiveness = $12,
    follow_up_responsiveness = $13,
    documentation_preferences = $14,
    presentation_style_preferences = $15,
    negotiation_approach = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, decision_making_style, risk_aversion_level, change_adoption_patterns, innovation_appetite,
    consensus_building_approach, conflict_resolution_style, communication_patterns, trust_building_factors,
    credibility_requirements, relationship_preferences, meeting_effectiveness, follow_up_responsiveness,
    documentation_preferences, presentation_style_preferences, negotiation_approach, created_at, updated_at;

-- name: DeleteBehavioralInsights :exec
DELETE FROM behavioral_insights
WHERE brief_id = $1;