-- name: GetTeamByID :one
SELECT * FROM teams
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetTeamByCode :one
SELECT * FROM teams
WHERE code = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListTeams :many
SELECT * FROM teams
WHERE deleted_at IS NULL
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: SearchTeamsByName :many
SELECT * FROM teams
WHERE deleted_at IS NULL
  AND name ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: GetTeamsByCountry :many
SELECT * FROM teams
WHERE country = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: CreateTeam :one
INSERT INTO teams (
    name, short_name, code, country, city, stadium, stadium_capacity,
    founded, logo, colors, website
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateTeam :one
UPDATE teams
SET
    name = COALESCE(sqlc.narg('name'), name),
    short_name = COALESCE(sqlc.narg('short_name'), short_name),
    country = COALESCE(sqlc.narg('country'), country),
    city = COALESCE(sqlc.narg('city'), city),
    stadium = COALESCE(sqlc.narg('stadium'), stadium),
    stadium_capacity = COALESCE(sqlc.narg('stadium_capacity'), stadium_capacity),
    logo = COALESCE(sqlc.narg('logo'), logo),
    colors = COALESCE(sqlc.narg('colors'), colors),
    website = COALESCE(sqlc.narg('website'), website)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteTeam :exec
UPDATE teams
SET deleted_at = NOW()
WHERE id = $1;

-- name: CountTeams :one
SELECT COUNT(*) FROM teams
WHERE deleted_at IS NULL;
