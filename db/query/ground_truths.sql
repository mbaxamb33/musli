
-- =============================================================================
-- GROUND TRUTH QUERIES
-- =============================================================================

-- name: CreateGroundTruth :one
INSERT INTO ground_truth (
    master_brief_id, field_name, field_value, confidence_score, source_brief_ids
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, master_brief_id, field_name, field_value, confidence_score, source_brief_ids, last_updated;

-- name: GetGroundTruthByMasterBriefAndField :one
SELECT id, master_brief_id, field_name, field_value, confidence_score, source_brief_ids, last_updated
FROM ground_truth
WHERE master_brief_id = $1 AND field_name = $2;

-- name: ListGroundTruthByMasterBrief :many
SELECT id, master_brief_id, field_name, field_value, confidence_score, source_brief_ids, last_updated
FROM ground_truth
WHERE master_brief_id = $1
ORDER BY field_name;

-- name: UpdateGroundTruth :one
UPDATE ground_truth
SET field_value = $3,
    confidence_score = $4,
    source_brief_ids = $5,
    last_updated = CURRENT_TIMESTAMP
WHERE master_brief_id = $1 AND field_name = $2
RETURNING id, master_brief_id, field_name, field_value, confidence_score, source_brief_ids, last_updated;

-- name: UpsertGroundTruth :one
INSERT INTO ground_truth (master_brief_id, field_name, field_value, confidence_score, source_brief_ids)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (master_brief_id, field_name)
DO UPDATE SET
    field_value = EXCLUDED.field_value,
    confidence_score = EXCLUDED.confidence_score,
    source_brief_ids = EXCLUDED.source_brief_ids,
    last_updated = CURRENT_TIMESTAMP
RETURNING id, master_brief_id, field_name, field_value, confidence_score, source_brief_ids, last_updated;

-- name: DeleteGroundTruth :exec
DELETE FROM ground_truth
WHERE master_brief_id = $1 AND field_name = $2;

-- name: DeleteAllGroundTruthForMasterBrief :exec
DELETE FROM ground_truth
WHERE master_brief_id = $1;