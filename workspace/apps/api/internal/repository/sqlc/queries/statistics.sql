-- Player Statistics Queries

-- name: GetPlayerStatsByID :one
SELECT * FROM player_statistics
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetPlayerStatsByPlayerAndSeason :one
SELECT * FROM player_statistics
WHERE player_id = $1 AND season = $2 AND competition = $3 AND deleted_at IS NULL
LIMIT 1;

-- name: GetPlayerStatsByPlayer :many
SELECT * FROM player_statistics
WHERE player_id = $1 AND deleted_at IS NULL
ORDER BY season DESC, competition;

-- name: GetTopScorers :many
SELECT
    ps.*,
    p.full_name,
    p.position,
    t.name as team_name,
    t.short_name as team_short_name
FROM player_statistics ps
JOIN players p ON ps.player_id = p.id AND p.deleted_at IS NULL
JOIN teams t ON p.team_id = t.id AND t.deleted_at IS NULL
WHERE ps.season = $1
  AND ps.competition = $2
  AND ps.deleted_at IS NULL
ORDER BY ps.goals DESC, ps.assists DESC
LIMIT $3;

-- name: GetTopAssisters :many
SELECT
    ps.*,
    p.full_name,
    p.position,
    t.name as team_name
FROM player_statistics ps
JOIN players p ON ps.player_id = p.id AND p.deleted_at IS NULL
JOIN teams t ON p.team_id = t.id AND t.deleted_at IS NULL
WHERE ps.season = $1
  AND ps.competition = $2
  AND ps.deleted_at IS NULL
ORDER BY ps.assists DESC, ps.goals DESC
LIMIT $3;

-- name: CreatePlayerStats :one
INSERT INTO player_statistics (
    player_id, season, competition, matches_played, matches_started, minutes_played,
    sub_on, sub_off, goals, assists, shots_total, shots_on_target, shot_accuracy,
    goal_conversion, passes_total, passes_completed, pass_accuracy, key_passes, crosses,
    tackles, tackles_won, interceptions, clearances, blocked_shots, duels, duels_won,
    aerial_duels, aerial_duels_won, yellow_cards, red_cards, fouls, fouls_drawn,
    clean_sheets, goals_conceded, saves_total, save_percentage, penalties_saved
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
    $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37
)
RETURNING *;

-- name: UpdatePlayerStats :one
UPDATE player_statistics
SET
    matches_played = $2,
    matches_started = $3,
    minutes_played = $4,
    sub_on = $5,
    sub_off = $6,
    goals = $7,
    assists = $8,
    shots_total = $9,
    shots_on_target = $10,
    shot_accuracy = $11,
    goal_conversion = $12,
    passes_total = $13,
    passes_completed = $14,
    pass_accuracy = $15,
    key_passes = $16,
    crosses = $17,
    tackles = $18,
    tackles_won = $19,
    interceptions = $20,
    clearances = $21,
    blocked_shots = $22,
    duels = $23,
    duels_won = $24,
    aerial_duels = $25,
    aerial_duels_won = $26,
    yellow_cards = $27,
    red_cards = $28,
    fouls = $29,
    fouls_drawn = $30
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePlayerStats :exec
UPDATE player_statistics
SET deleted_at = NOW()
WHERE id = $1;

-- Team Statistics Queries

-- name: GetTeamStatsByID :one
SELECT * FROM team_statistics
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetTeamStatsByTeamAndSeason :one
SELECT * FROM team_statistics
WHERE team_id = $1 AND season = $2 AND competition = $3 AND deleted_at IS NULL
LIMIT 1;

-- name: GetTeamStatsByTeam :many
SELECT * FROM team_statistics
WHERE team_id = $1 AND deleted_at IS NULL
ORDER BY season DESC, competition;

-- name: GetLeagueTable :many
SELECT
    ts.*,
    t.name,
    t.short_name,
    t.code,
    t.logo
FROM team_statistics ts
JOIN teams t ON ts.team_id = t.id AND t.deleted_at IS NULL
WHERE ts.season = $1
  AND ts.competition = $2
  AND ts.deleted_at IS NULL
ORDER BY ts.points DESC, ts.goal_difference DESC, ts.goals_scored DESC;

-- name: CreateTeamStats :one
INSERT INTO team_statistics (
    team_id, season, competition, matches_played, wins, draws, losses, points, position,
    goals_scored, goals_conceded, goal_difference, clean_sheets, goals_per_match,
    home_wins, home_draws, home_losses, away_wins, away_draws, away_losses,
    possession, pass_accuracy, shots_per_match, shots_on_target_percentage,
    yellow_cards, red_cards, current_form
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
    $20, $21, $22, $23, $24, $25, $26, $27
)
RETURNING *;

-- name: UpdateTeamStats :one
UPDATE team_statistics
SET
    matches_played = $2,
    wins = $3,
    draws = $4,
    losses = $5,
    points = $6,
    position = $7,
    goals_scored = $8,
    goals_conceded = $9,
    goal_difference = $10,
    clean_sheets = $11,
    goals_per_match = $12,
    home_wins = $13,
    home_draws = $14,
    home_losses = $15,
    away_wins = $16,
    away_draws = $17,
    away_losses = $18,
    possession = $19,
    pass_accuracy = $20,
    shots_per_match = $21,
    shots_on_target_percentage = $22,
    yellow_cards = $23,
    red_cards = $24,
    current_form = $25
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteTeamStats :exec
UPDATE team_statistics
SET deleted_at = NOW()
WHERE id = $1;
