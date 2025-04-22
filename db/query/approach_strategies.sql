-- name: CreateApproachStrategy :one
INSERT INTO approach_strategies (
  name,
  description,
  recommended_score_min,
  recommended_score_max
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetApproachStrategy :one
SELECT * FROM approach_strategies
WHERE strategy_id = $1 LIMIT 1;

-- name: ListApproachStrategies :many
SELECT * FROM approach_strategies
ORDER BY recommended_score_min
LIMIT $1 OFFSET $2;

-- name: ListApproachStrategiesInScoreRange :many
SELECT * FROM approach_strategies
WHERE recommended_score_min <= $1 AND recommended_score_max >= $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateApproachStrategy :one
UPDATE approach_strategies
SET
  name = COALESCE($2, name),
  description = COALESCE($3, description),
  recommended_score_min = COALESCE($4, recommended_score_min),
  recommended_score_max = COALESCE($5, recommended_score_max),
  updated_at = CURRENT_TIMESTAMP
WHERE strategy_id = $1
RETURNING *;

-- name: DeleteApproachStrategy :exec
DELETE FROM approach_strategies
WHERE strategy_id = $1;

-- name: GetApproachStrategyByName :one
SELECT * FROM approach_strategies
WHERE name = $1 LIMIT 1;

-- name: GetRecommendedStrategyForScore :one
SELECT * FROM approach_strategies
WHERE recommended_score_min <= $1 AND recommended_score_max >= $1
LIMIT 1;

-- name: GetApproachStrategyWithUsageCount :many
SELECT 
  approach_strategies.*,
  COUNT(pc.project_company_id) AS usage_count
FROM approach_strategies
LEFT JOIN project_companies pc ON approach_strategies.strategy_id = pc.approach_strategy_id
GROUP BY approach_strategies.strategy_id
ORDER BY approach_strategies.recommended_score_min
LIMIT $1 OFFSET $2;