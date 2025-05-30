-- =============================================================================
-- COMPETITIVE INTELLIGENCE QUERIES
-- =============================================================================

-- name: CreateCompetitiveIntelligence :one
INSERT INTO competitive_intelligence (
    brief_id, competitors_in_evaluation, preferred_vendor_bias, previous_vendor_history,
    competitive_strengths, competitive_weaknesses, pricing_expectations, feature_comparison_matrix,
    vendor_selection_criteria, criteria_weighting, evaluation_process, reference_requirements,
    proof_of_concept_needs, pilot_program_scope, final_presentation_format, decision_timeline
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, competitors_in_evaluation, preferred_vendor_bias, previous_vendor_history,
    competitive_strengths, competitive_weaknesses, pricing_expectations, feature_comparison_matrix,
    vendor_selection_criteria, criteria_weighting, evaluation_process, reference_requirements,
    proof_of_concept_needs, pilot_program_scope, final_presentation_format, decision_timeline,
    created_at, updated_at;

-- name: GetCompetitiveIntelligenceByBriefID :one
SELECT id, brief_id, competitors_in_evaluation, preferred_vendor_bias, previous_vendor_history,
    competitive_strengths, competitive_weaknesses, pricing_expectations, feature_comparison_matrix,
    vendor_selection_criteria, criteria_weighting, evaluation_process, reference_requirements,
    proof_of_concept_needs, pilot_program_scope, final_presentation_format, decision_timeline,
    created_at, updated_at
FROM competitive_intelligence
WHERE brief_id = $1;

-- name: UpdateCompetitiveIntelligence :one
UPDATE competitive_intelligence
SET competitors_in_evaluation = $2,
    preferred_vendor_bias = $3,
    previous_vendor_history = $4,
    competitive_strengths = $5,
    competitive_weaknesses = $6,
    pricing_expectations = $7,
    feature_comparison_matrix = $8,
    vendor_selection_criteria = $9,
    criteria_weighting = $10,
    evaluation_process = $11,
    reference_requirements = $12,
    proof_of_concept_needs = $13,
    pilot_program_scope = $14,
    final_presentation_format = $15,
    decision_timeline = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, competitors_in_evaluation, preferred_vendor_bias, previous_vendor_history,
    competitive_strengths, competitive_weaknesses, pricing_expectations, feature_comparison_matrix,
    vendor_selection_criteria, criteria_weighting, evaluation_process, reference_requirements,
    proof_of_concept_needs, pilot_program_scope, final_presentation_format, decision_timeline,
    created_at, updated_at;

-- name: DeleteCompetitiveIntelligence :exec
DELETE FROM competitive_intelligence
WHERE brief_id = $1;