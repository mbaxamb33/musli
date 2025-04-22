-- name: CreateMatchingCriteria :one
INSERT INTO matching_criteria (
  name,
  description,
  weight,
  is_active
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetMatchingCriteria :one
SELECT * FROM matching_criteria
WHERE criteria_id = $1 LIMIT 1;

-- name: ListMatchingCriteria :many
SELECT * FROM matching_criteria
ORDER BY weight DESC, name
LIMIT $1 OFFSET $2;

-- name: ListActiveMatchingCriteria :many
SELECT * FROM matching_criteria
WHERE is_active = TRUE
ORDER BY weight DESC, name
LIMIT $1 OFFSET $2;

-- name: UpdateMatchingCriteria :one
UPDATE matching_criteria
SET
  name = COALESCE($2, name),
  description = COALESCE($3, description),
  weight = COALESCE($4, weight),
  is_active = COALESCE($5, is_active),
  updated_at = CURRENT_TIMESTAMP
WHERE criteria_id = $1
RETURNING *;

-- name: ToggleCriteriaActive :one
UPDATE matching_criteria
SET 
  is_active = NOT is_active,
  updated_at = CURRENT_TIMESTAMP
WHERE criteria_id = $1
RETURNING *;

-- name: DeleteMatchingCriteria :exec
DELETE FROM matching_criteria
WHERE criteria_id = $1;

-- name: GetMatchingCriteriaByName :one
SELECT * FROM matching_criteria
WHERE name = $1 LIMIT 1;

-- name: GetTotalCriteriaWeight :one
SELECT SUM(weight) AS total_weight
FROM matching_criteria
WHERE is_active = TRUE;

-- name: GetMatchingCriteriaWithUsageCount :many
SELECT 
  mc.*,
  COUNT(msd.score_detail_id) AS usage_count
FROM matching_criteria mc
LEFT JOIN matching_scores_detail msd ON mc.criteria_id = msd.criteria_id
GROUP BY mc.criteria_id
ORDER BY mc.weight DESC, mc.name
LIMIT $1 OFFSET $2;