-- name: CreateNeedsDatasourceMatch :one
INSERT INTO needs_datasource_matches (
    need_id, datasource_id, match_score
)
VALUES ($1, $2, $3)
RETURNING match_id, need_id, datasource_id, match_score, created_at;

-- name: GetNeedsDatasourceMatchByID :one
SELECT match_id, need_id, datasource_id, match_score, created_at
FROM needs_datasource_matches
WHERE match_id = $1;

-- name: ListMatchesByNeed :many
SELECT match_id, need_id, datasource_id, match_score, created_at
FROM needs_datasource_matches
WHERE need_id = $1
ORDER BY match_score DESC
LIMIT $2 OFFSET $3;

-- name: ListMatchesByDatasource :many
SELECT match_id, need_id, datasource_id, match_score, created_at
FROM needs_datasource_matches
WHERE datasource_id = $1
ORDER BY match_score DESC
LIMIT $2 OFFSET $3;

-- name: UpdateNeedsDatasourceMatch :one
UPDATE needs_datasource_matches
SET match_score = $3
WHERE need_id = $1 AND datasource_id = $2
RETURNING match_id, need_id, datasource_id, match_score, created_at;

-- name: DeleteNeedsDatasourceMatch :exec
DELETE FROM needs_datasource_matches
WHERE match_id = $1;