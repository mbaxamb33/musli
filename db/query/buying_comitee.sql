-- =============================================================================
-- BUYING COMMITTEE QUERIES
-- =============================================================================

-- name: CreateBuyingCommittee :one
INSERT INTO buying_committee (
    brief_id, economic_buyer_name, economic_buyer_title, economic_buyer_influence, economic_buyer_motivations,
    technical_buyer_name, technical_buyer_concerns, user_buyer_representatives, coach_champion_name,
    coach_influence_level, blocker_identification, blocker_concerns, committee_dynamics,
    decision_making_process, consensus_requirements, individual_risk_tolerance, career_motivations,
    personal_success_metrics, relationship_mapping, communication_preferences, influence_networks
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
RETURNING id, brief_id, economic_buyer_name, economic_buyer_title, economic_buyer_influence, economic_buyer_motivations,
    technical_buyer_name, technical_buyer_concerns, user_buyer_representatives, coach_champion_name,
    coach_influence_level, blocker_identification, blocker_concerns, committee_dynamics,
    decision_making_process, consensus_requirements, individual_risk_tolerance, career_motivations,
    personal_success_metrics, relationship_mapping, communication_preferences, influence_networks,
    created_at, updated_at;

-- name: GetBuyingCommitteeByBriefID :one
SELECT id, brief_id, economic_buyer_name, economic_buyer_title, economic_buyer_influence, economic_buyer_motivations,
    technical_buyer_name, technical_buyer_concerns, user_buyer_representatives, coach_champion_name,
    coach_influence_level, blocker_identification, blocker_concerns, committee_dynamics,
    decision_making_process, consensus_requirements, individual_risk_tolerance, career_motivations,
    personal_success_metrics, relationship_mapping, communication_preferences, influence_networks,
    created_at, updated_at
FROM buying_committee
WHERE brief_id = $1;

-- name: UpdateBuyingCommittee :one
UPDATE buying_committee
SET economic_buyer_name = $2,
    economic_buyer_title = $3,
    economic_buyer_influence = $4,
    economic_buyer_motivations = $5,
    technical_buyer_name = $6,
    technical_buyer_concerns = $7,
    user_buyer_representatives = $8,
    coach_champion_name = $9,
    coach_influence_level = $10,
    blocker_identification = $11,
    blocker_concerns = $12,
    committee_dynamics = $13,
    decision_making_process = $14,
    consensus_requirements = $15,
    individual_risk_tolerance = $16,
    career_motivations = $17,
    personal_success_metrics = $18,
    relationship_mapping = $19,
    communication_preferences = $20,
    influence_networks = $21,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, economic_buyer_name, economic_buyer_title, economic_buyer_influence, economic_buyer_motivations,
    technical_buyer_name, technical_buyer_concerns, user_buyer_representatives, coach_champion_name,
    coach_influence_level, blocker_identification, blocker_concerns, committee_dynamics,
    decision_making_process, consensus_requirements, individual_risk_tolerance, career_motivations,
    personal_success_metrics, relationship_mapping, communication_preferences, influence_networks,
    created_at, updated_at;

-- name: DeleteBuyingCommittee :exec
DELETE FROM buying_committee
WHERE brief_id = $1;