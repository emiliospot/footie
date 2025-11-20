# 🚀 Footie API - sqlc + pgx + golang-migrate Stack

This is a **production-grade** football analytics API using the industry-standard stack:

- **sqlc** - Type-safe SQL → Go code generation
- **pgx** - Fastest PostgreSQL driver for Go
- **golang-migrate** - Database migration tool

This is the same stack used by betting companies, sports data companies, and real-time analytics systems.

---

## 📁 Project Structure

```
apps/api/
├── migrations/                    # Database migrations
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── internal/
│   ├── repository/
│   │   └── sqlc/
│   │       ├── queries/          # SQL queries (you write these)
│   │       │   ├── users.sql
│   │       │   ├── teams.sql
│   │       │   ├── players.sql
│   │       │   ├── matches.sql
│   │       │   ├── match_events.sql
│   │       │   └── statistics.sql
│   │       ├── *.sql.go          # Generated Go code (don't edit)
│   │       ├── models.go         # Generated models
│   │       └── querier.go        # Generated interface
│   └── infrastructure/
│       └── database/
│           ├── pgx.go            # Connection pool setup
│           └── migrate.go        # Migration runner
└── sqlc.yaml                     # sqlc configuration
```

---

## 🛠️ Available Commands

### Database Migrations

```bash
# Run all pending migrations
make db-up

# Rollback last migration
make db-down

# Drop all tables and re-run migrations (⚠️  destructive)
make db-reset

# Create new migration
make db-create name=add_xg_field

# Check migration status
make db-status

# Force migration version (if stuck)
make db-force version=1
```

### sqlc Code Generation

```bash
# Generate Go code from SQL queries
make sqlc-generate

# Vet SQL queries for errors
make sqlc-vet
```

### Development

```bash
# Run with hot-reload
make dev

# Run tests
make test

# Lint code
make lint
```

---

## 📝 How to Write Queries

### 1. Write SQL in `internal/repository/sqlc/queries/*.sql`

Example: `internal/repository/sqlc/queries/players.sql`

```sql
-- name: GetPlayerShots :many
SELECT 
    me.*,
    me.metadata->>'xg' as expected_goals,
    me.metadata->>'shot_type' as shot_type
FROM match_events me
WHERE me.player_id = $1 
  AND me.event_type = 'shot' 
  AND me.deleted_at IS NULL
ORDER BY me.id DESC
LIMIT $2 OFFSET $3;
```

### 2. Generate Go Code

```bash
make sqlc-generate
```

### 3. Use in Your Code

```go
import "github.com/emiliospot/footie/api/internal/repository/sqlc"

// In your handler or service
queries := sqlc.New(pool)

shots, err := queries.GetPlayerShots(ctx, sqlc.GetPlayerShotsParams{
    PlayerID: 123,
    Limit:    10,
    Offset:   0,
})
```

**That's it!** Type-safe, fast, and clean.

---

## 🎯 Why This Stack?

### ✅ sqlc

- You write **raw SQL** (perfect for analytics)
- Generates **type-safe Go code**
- No ORM magic, no reflection
- Catches SQL errors at **compile time**

### ✅ pgx

- **3-5x faster** than database/sql
- Native PostgreSQL protocol
- Connection pooling built-in
- Used by production systems handling millions of queries

### ✅ golang-migrate

- Simple `.sql` migration files
- Works with Docker, Kubernetes, CI/CD
- Can run migrations on app startup
- Supports up/down migrations

---

## 📊 Example: Football Analytics Query

```sql
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
```

Generated Go code:

```go
type GetPlayerPassAccuracyRow struct {
    CompletedPasses         int64
    TotalPasses             int64
    PassAccuracyPercentage  float64
}

func (q *Queries) GetPlayerPassAccuracy(ctx context.Context, playerID int32) (GetPlayerPassAccuracyRow, error)
```

---

## 🔥 Advanced Analytics Queries

All queries support:

- **JSONB** for flexible event metadata (xG, pass types, shot locations)
- **Window functions** for rankings and running totals
- **CTEs** for complex analytics
- **Aggregations** with FILTER for conditional stats
- **Full-text search** with pg_trgm for player/team search

---

## 🐳 Docker Integration

Migrations run automatically on app startup in `main.go`:

```go
// Run database migrations first
if err := database.RunMigrations(cfg.Database.URL, "migrations"); err != nil {
    log.Fatal(err)
}
```

For Docker Compose, migrations run before the API starts.

---

## 📚 Resources

- [sqlc Documentation](https://docs.sqlc.dev/)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)

---

## 🎓 Learning Path

1. **Start with simple queries** (GetUserByID, ListTeams)
2. **Add filters and pagination** (GetMatchesByTeam)
3. **Write analytics queries** (GetTopScorers, GetPlayerPassAccuracy)
4. **Use JSONB for flexible data** (shot xG, pass completion, tracking data)
5. **Optimize with indexes** (already added in migrations)

---

## 💡 Pro Tips

1. **Always use prepared statements** (sqlc does this automatically)
2. **Use connection pooling** (pgx handles this)
3. **Write migrations in pairs** (up + down)
4. **Test queries with `sqlc vet`** before generating
5. **Use JSONB for analytics metadata** (xG, coordinates, etc.)

---

## 🚨 Important Notes

- **Never edit generated files** (`*.sql.go`, `models.go`, etc.)
- **Always run `make sqlc-generate`** after changing SQL queries
- **Test migrations locally** before deploying
- **Use transactions** for multi-step operations (pgx supports this)

---

Happy coding! ⚽📊

