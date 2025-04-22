-- name: CreateMatchingScoreDetail :one
INSERT INTO matching_scores_detail (
  project_company_id,
  criteria_id,
  score
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetMatchingScoreDetail :one
SELECT * FROM matching_scores_detail
WHERE project_company_id = $1 AND criteria_id = $2
LIMIT 1;

-- name: GetMatchingScoreDetailById :one
SELECT * FROM matching_scores_detail
WHERE score_detail_id = $1
LIMIT 1;

-- name: ListMatchingScoreDetailsForProjectCompany :many
SELECT 
  msd.*,
  mc.name AS criteria_name,
  mc.description AS criteria_description,
  mc.weight AS criteria_weight
FROM matching_scores_detail msd
JOIN matching_criteria mc ON msd.criteria_id = mc.criteria_id
WHERE msd.project_company_id = $1
ORDER BY mc.weight DESC, mc.name
LIMIT $2 OFFSET $3;

-- name: UpdateMatchingScoreDetail :one
UPDATE matching_scores_detail
SET
  score = $3,
  updated_at = CURRENT_TIMESTAMP
WHERE project_company_id = $1 AND criteria_id = $2
RETURNING *;

-- name: DeleteMatchingScoreDetail :exec
DELETE FROM matching_scores_detail
WHERE project_company_id = $1 AND criteria_id = $2;

-- name: DeleteMatchingScoreDetailById :exec
DELETE FROM matching_scores_detail
WHERE score_detail_id = $1;

-- name: DeleteAllMatchingScoreDetailsForProjectCompany :exec
DELETE FROM matching_scores_detail
WHERE project_company_id = $1;

-- name: CalculateWeightedScore :one
SELECT 
  SUM(msd.score * mc.weight) / SUM(mc.weight) AS weighted_score
FROM matching_scores_detail msd
JOIN matching_criteria mc ON msd.criteria_id = mc.criteria_id
WHERE msd.project_company_id = $1;

-- name: GetCriteriaScoresForProject :many
SELECT 
  c.name AS company_name,
  mc.name AS criteria_name,
  msd.score,
  mc.weight,
  (msd.score * mc.weight) AS weighted_score
FROM matching_scores_detail msd
JOIN project_companies pc ON msd.project_company_id = pc.project_company_id
JOIN companies c ON pc.company_id = c.company_id
JOIN matching_criteria mc ON msd.criteria_id = mc.criteria_id
WHERE pc.project_id = $1
ORDER BY c.name, mc.weight DESC;