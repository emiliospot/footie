-- name: GetMatchEventByID :one
SELECT * FROM match_events
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetMatchEvents :many
SELECT * FROM match_events
WHERE match_id = $1 AND deleted_at IS NULL
ORDER BY minute ASC, extra_minute ASC, id ASC;

-- name: GetMatchEventsByType :many
SELECT * FROM match_events
WHERE match_id = $1 AND event_type = $2 AND deleted_at IS NULL
ORDER BY minute ASC, extra_minute ASC;

-- name: GetPlayerEvents :many
SELECT * FROM match_events
WHERE player_id = $1 AND deleted_at IS NULL
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: GetTeamEventsInMatch :many
SELECT * FROM match_events
WHERE match_id = $1 AND team_id = $2 AND deleted_at IS NULL
ORDER BY minute ASC, extra_minute ASC;

-- name: GetGoalsByMatch :many
SELECT * FROM match_events
WHERE match_id = $1 AND event_type = 'goal' AND deleted_at IS NULL
ORDER BY minute ASC, extra_minute ASC;

-- name: GetCardsByMatch :many
SELECT * FROM match_events
WHERE match_id = $1
  AND event_type IN ('yellow_card', 'red_card')
  AND deleted_at IS NULL
ORDER BY minute ASC, extra_minute ASC;

-- name: GetShotsByMatch :many
SELECT * FROM match_events
WHERE match_id = $1 AND event_type = 'shot' AND deleted_at IS NULL
ORDER BY minute ASC, extra_minute ASC;

-- name: GetPassesByMatch :many
SELECT * FROM match_events
WHERE match_id = $1 AND event_type = 'pass' AND deleted_at IS NULL
ORDER BY minute ASC, extra_minute ASC;

-- name: CreateMatchEvent :one
INSERT INTO match_events (
    match_id, team_id, player_id, secondary_player_id, event_type,
    minute, extra_minute, position_x, position_y, description, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateMatchEvent :one
UPDATE match_events
SET
    event_type = COALESCE(sqlc.narg('event_type'), event_type),
    minute = COALESCE(sqlc.narg('minute'), minute),
    extra_minute = COALESCE(sqlc.narg('extra_minute'), extra_minute),
    position_x = COALESCE(sqlc.narg('position_x'), position_x),
    position_y = COALESCE(sqlc.narg('position_y'), position_y),
    description = COALESCE(sqlc.narg('description'), description),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteMatchEvent :exec
UPDATE match_events
SET deleted_at = NOW()
WHERE id = $1;

-- name: CountMatchEvents :one
SELECT COUNT(*) FROM match_events
WHERE match_id = $1 AND deleted_at IS NULL;

-- name: CountEventsByType :one
SELECT COUNT(*) FROM match_events
WHERE match_id = $1 AND event_type = $2 AND deleted_at IS NULL;

-- Analytics queries for match events
-- name: GetPlayerShotsWithXG :many
SELECT
    me.*,
    me.metadata->>'xg' as expected_goals,
    me.metadata->>'shot_type' as shot_type,
    me.metadata->>'body_part' as body_part
FROM match_events me
WHERE me.player_id = $1
  AND me.event_type = 'shot'
  AND me.deleted_at IS NULL
ORDER BY me.id DESC
LIMIT $2 OFFSET $3;

-- name: GetPlayerPassAccuracy :one
SELECT
    COUNT(*) FILTER (WHERE metadata->>'completed' = 'true') as completed_passes,
    COUNT(*) as total_passes,
    CASE
        WHEN COUNT(*) > 0 THEN
            ROUND((COUNT(*) FILTER (WHERE metadata->>'completed' = 'true')::numeric / COUNT(*) * 100), 2)
        ELSE 0
    END as pass_accuracy_percentage
FROM match_events
WHERE player_id = $1
  AND event_type = 'pass'
  AND deleted_at IS NULL;

-- name: GetTeamPossessionEvents :many
SELECT
    team_id,
    COUNT(*) as total_events,
    COUNT(*) FILTER (WHERE event_type = 'pass') as passes,
    COUNT(*) FILTER (WHERE event_type = 'shot') as shots
FROM match_events
WHERE match_id = $1 AND deleted_at IS NULL
GROUP BY team_id;
