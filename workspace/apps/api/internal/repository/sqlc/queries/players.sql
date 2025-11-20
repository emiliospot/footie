-- name: GetPlayerByID :one
SELECT * FROM players
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetPlayerWithTeam :one
SELECT
    p.*,
    t.id as team_id,
    t.name as team_name,
    t.short_name as team_short_name,
    t.code as team_code
FROM players p
LEFT JOIN teams t ON p.team_id = t.id AND t.deleted_at IS NULL
WHERE p.id = $1 AND p.deleted_at IS NULL
LIMIT 1;

-- name: ListPlayers :many
SELECT * FROM players
WHERE deleted_at IS NULL
ORDER BY full_name
LIMIT $1 OFFSET $2;

-- name: GetPlayersByTeam :many
SELECT * FROM players
WHERE team_id = $1 AND deleted_at IS NULL
ORDER BY shirt_number NULLS LAST, full_name;

-- name: GetPlayersByPosition :many
SELECT * FROM players
WHERE position = $1 AND deleted_at IS NULL
ORDER BY full_name
LIMIT $2 OFFSET $3;

-- name: SearchPlayersByName :many
SELECT * FROM players
WHERE deleted_at IS NULL
  AND full_name ILIKE '%' || $1 || '%'
ORDER BY full_name
LIMIT $2 OFFSET $3;

-- name: CreatePlayer :one
INSERT INTO players (
    team_id, first_name, last_name, full_name, date_of_birth, nationality,
    position, shirt_number, height, weight, preferred_foot, photo
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdatePlayer :one
UPDATE players
SET
    team_id = COALESCE(sqlc.narg('team_id'), team_id),
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    date_of_birth = COALESCE(sqlc.narg('date_of_birth'), date_of_birth),
    nationality = COALESCE(sqlc.narg('nationality'), nationality),
    position = COALESCE(sqlc.narg('position'), position),
    shirt_number = COALESCE(sqlc.narg('shirt_number'), shirt_number),
    height = COALESCE(sqlc.narg('height'), height),
    weight = COALESCE(sqlc.narg('weight'), weight),
    preferred_foot = COALESCE(sqlc.narg('preferred_foot'), preferred_foot),
    photo = COALESCE(sqlc.narg('photo'), photo)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeletePlayer :exec
UPDATE players
SET deleted_at = NOW()
WHERE id = $1;

-- name: CountPlayers :one
SELECT COUNT(*) FROM players
WHERE deleted_at IS NULL;

-- name: CountPlayersByTeam :one
SELECT COUNT(*) FROM players
WHERE team_id = $1 AND deleted_at IS NULL;
