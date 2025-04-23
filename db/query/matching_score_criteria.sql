-- name: CreateMatchingScoreCriteria :one
INSERT INTO matching_score_criteria (
    name, description, weight
)
VALUES ($1, $2, $3)
RETURNING criteria_id, name, description, weight;

-- name: GetMatchingScoreCriteriaByID :one
SELECT criteria_id, name, description, weight
FROM matching_score_criteria
WHERE criteria_id = $1;

-- name: ListMatchingScoreCriteria :many
SELECT criteria_id, name, description, weight
FROM matching_score_criteria
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateMatchingScoreCriteria :one
UPDATE matching_score_criteria
SET name = $2,
    description = $3,
    weight = $4
WHERE criteria_id = $1
RETURNING criteria_id, name, description, weight;

-- name: DeleteMatchingScoreCriteria :exec
DELETE FROM matching_score_criteria
WHERE criteria_id = $1;
