-- name: GetMatchByID :one
SELECT * FROM matches
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetMatchWithTeams :one
SELECT
    m.*,
    ht.id as home_team_id,
    ht.name as home_team_name,
    ht.short_name as home_team_short_name,
    ht.code as home_team_code,
    ht.logo as home_team_logo,
    at.id as away_team_id,
    at.name as away_team_name,
    at.short_name as away_team_short_name,
    at.code as away_team_code,
    at.logo as away_team_logo
FROM matches m
LEFT JOIN teams ht ON m.home_team_id = ht.id AND ht.deleted_at IS NULL
LEFT JOIN teams at ON m.away_team_id = at.id AND at.deleted_at IS NULL
WHERE m.id = $1 AND m.deleted_at IS NULL
LIMIT 1;

-- name: ListMatches :many
SELECT * FROM matches
WHERE deleted_at IS NULL
ORDER BY match_date DESC
LIMIT $1 OFFSET $2;

-- name: GetMatchesByTeam :many
SELECT * FROM matches
WHERE (home_team_id = $1 OR away_team_id = $1)
  AND deleted_at IS NULL
ORDER BY match_date DESC
LIMIT $2 OFFSET $3;

-- name: GetMatchesByCompetition :many
SELECT * FROM matches
WHERE competition = $1 AND deleted_at IS NULL
ORDER BY match_date DESC
LIMIT $2 OFFSET $3;

-- name: GetMatchesBySeason :many
SELECT * FROM matches
WHERE season = $1 AND deleted_at IS NULL
ORDER BY match_date DESC
LIMIT $2 OFFSET $3;

-- name: GetMatchesByCompetitionAndSeason :many
SELECT * FROM matches
WHERE competition = $1 AND season = $2 AND deleted_at IS NULL
ORDER BY match_date DESC;

-- name: GetMatchesByStatus :many
SELECT * FROM matches
WHERE status = $1 AND deleted_at IS NULL
ORDER BY match_date DESC
LIMIT $2 OFFSET $3;

-- name: GetUpcomingMatches :many
SELECT * FROM matches
WHERE match_date > NOW() AND status = 'scheduled' AND deleted_at IS NULL
ORDER BY match_date ASC
LIMIT $1;

-- name: GetLiveMatches :many
SELECT * FROM matches
WHERE status = 'live' AND deleted_at IS NULL
ORDER BY match_date DESC;

-- name: CreateMatch :one
INSERT INTO matches (
    home_team_id, away_team_id, match_date, competition, season, round,
    stadium, attendance, status, referee, home_team_score, away_team_score
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateMatch :one
UPDATE matches
SET
    match_date = COALESCE(sqlc.narg('match_date'), match_date),
    competition = COALESCE(sqlc.narg('competition'), competition),
    season = COALESCE(sqlc.narg('season'), season),
    round = COALESCE(sqlc.narg('round'), round),
    stadium = COALESCE(sqlc.narg('stadium'), stadium),
    attendance = COALESCE(sqlc.narg('attendance'), attendance),
    status = COALESCE(sqlc.narg('status'), status),
    referee = COALESCE(sqlc.narg('referee'), referee),
    home_team_score = COALESCE(sqlc.narg('home_team_score'), home_team_score),
    away_team_score = COALESCE(sqlc.narg('away_team_score'), away_team_score)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: UpdateMatchScore :one
UPDATE matches
SET
    home_team_score = $2,
    away_team_score = $3
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateMatchStatus :one
UPDATE matches
SET status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteMatch :exec
UPDATE matches
SET deleted_at = NOW()
WHERE id = $1;

-- name: CountMatches :one
SELECT COUNT(*) FROM matches
WHERE deleted_at IS NULL;

-- name: CountMatchesByTeam :one
SELECT COUNT(*) FROM matches
WHERE (home_team_id = $1 OR away_team_id = $1) AND deleted_at IS NULL;
