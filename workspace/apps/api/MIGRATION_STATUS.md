# 🔄 sqlc + pgx Migration Status

## ✅ Completed

### Database Layer

- ✅ **golang-migrate** installed and configured
- ✅ **pgx v5** connection pool setup (`internal/infrastructure/database/pgx.go`)
- ✅ **Migration runner** implemented (`internal/infrastructure/database/migrate.go`)
- ✅ **Initial schema migration** created (`migrations/000001_init_schema.up.sql`)
- ✅ **Rollback migration** created (`migrations/000001_init_schema.down.sql`)
- ✅ **Migrations tested** - All 7 tables + 46 indexes created successfully

### sqlc Code Generation

- ✅ **sqlc** installed and configured (`sqlc.yaml`)
- ✅ **SQL queries** written for all entities:
  - `queries/users.sql` - User CRUD + authentication
  - `queries/teams.sql` - Team management + search
  - `queries/players.sql` - Player management + team relations
  - `queries/matches.sql` - Match management + filtering
  - `queries/match_events.sql` - Event tracking + analytics
  - `queries/statistics.sql` - Player/team stats + leaderboards
- ✅ **Type-safe Go code** generated successfully
- ✅ **Models** generated (`internal/repository/sqlc/models.go`)
- ✅ **Querier interface** generated (`internal/repository/sqlc/querier.go`)

### Build & Dependencies

- ✅ **Go dependencies** updated (Dependabot PRs merged)
- ✅ **GORM code** removed (`postgres.go` deleted)
- ✅ **Router** updated to use `*pgxpool.Pool`
- ✅ **API compiles** successfully
- ✅ **Type errors** resolved
- ✅ **go vet** passes

### Documentation

- ✅ **README files** updated with sqlc + pgx references
- ✅ **QUICKSTART** updated with migration steps
- ✅ **README_SQLC.md** created with complete usage guide
- ✅ **Database commands** added to root `package.json`
- ✅ **Makefile** updated with db commands

### Scripts & Tools

- ✅ `npm run db:up` - Run migrations
- ✅ `npm run db:down` - Rollback migrations
- ✅ `npm run db:reset` - Drop all & re-run
- ✅ `npm run db:status` - Check migration version
- ✅ `npm run sqlc:generate` - Generate Go code
- ✅ `npm run sqlc:vet` - Validate SQL queries

---

## 🚧 TODO: Handler Refactoring

The following handlers need to be refactored to use sqlc queries instead of GORM:

### Priority 1: Authentication & Users

- [ ] **AuthHandler** (`internal/api/handlers/auth_handler.go`)
  - Register (create user)
  - Login (get user by email)
  - RefreshToken (get user by ID)
- [ ] **UserHandler** (`internal/api/handlers/user_handler.go`)
  - GetCurrentUser
  - UpdateCurrentUser
  - GetUser
  - ListUsers (admin)
  - UpdateUserRole (admin)
  - DeleteUser (admin)

### Priority 2: Core Entities

- [ ] **TeamHandler** (`internal/api/handlers/team_handler.go`)
  - ListTeams
  - GetTeam
  - CreateTeam
  - UpdateTeam
  - DeleteTeam
  - GetTeamPlayers
  - GetTeamStatistics

- [ ] **PlayerHandler** (`internal/api/handlers/player_handler.go`)
  - ListPlayers
  - GetPlayer
  - CreatePlayer
  - UpdatePlayer
  - DeletePlayer
  - GetPlayerStatistics

- [ ] **MatchHandler** (`internal/api/handlers/match_handler.go`)
  - ListMatches
  - GetMatch
  - CreateMatch
  - UpdateMatch
  - DeleteMatch
  - GetMatchEvents
  - CreateMatchEvent

---

## 📝 How to Refactor a Handler

### Example: UserHandler

**Before (GORM):**

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    var user models.User
    if err := h.db.First(&user, id).Error; err != nil {
        // handle error
    }
    c.JSON(200, user)
}
```

**After (sqlc + pgx):**

```go
import "github.com/emiliospot/footie/api/internal/repository/sqlc"

type UserHandler struct {
    queries *sqlc.Queries
    logger  *logger.Logger
}

func NewUserHandler(pool *pgxpool.Pool, logger *logger.Logger) *UserHandler {
    return &UserHandler{
        queries: sqlc.New(pool),
        logger:  logger,
    }
}

func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.queries.GetUserByID(c.Request.Context(), int32(id))
    if err != nil {
        // handle error
    }
    c.JSON(200, user)
}
```

### Steps:

1. Replace `*gorm.DB` with `*sqlc.Queries` in handler struct
2. Update constructor to accept `*pgxpool.Pool`
3. Replace GORM queries with sqlc methods
4. Use `c.Request.Context()` for all queries
5. Handle `pgx.ErrNoRows` for not found cases
6. Update error handling

---

## 🎯 Benefits After Migration

### Performance

- **3-5x faster** queries (pgx vs database/sql)
- **Zero reflection** overhead
- **Connection pooling** optimized for production

### Type Safety

- **Compile-time** SQL validation
- **Type-safe** parameters and results
- **No runtime** SQL errors

### Developer Experience

- **Raw SQL** - perfect for complex analytics
- **IDE autocomplete** for all queries
- **Easy to optimize** - just write better SQL

### Production Ready

- Used by **betting companies**
- Used by **sports analytics** teams
- **Industry standard** for high-performance apps

---

## 🔧 Maintenance

### Adding New Queries

1. Write SQL in `internal/repository/sqlc/queries/*.sql`
2. Run `npm run sqlc:generate`
3. Use generated functions in handlers

### Creating Migrations

1. Run `npm run db:create name=add_xg_field`
2. Edit `migrations/NNNN_add_xg_field.up.sql`
3. Edit `migrations/NNNN_add_xg_field.down.sql`
4. Run `npm run db:up`

### Testing Queries

```bash
# Check migration status
npm run db:status

# Validate SQL
npm run sqlc:vet

# Generate code
npm run sqlc:generate
```

---

## 📚 Resources

- [sqlc Documentation](https://docs.sqlc.dev/)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [Project README_SQLC.md](./README_SQLC.md)

---

**Status:** Database layer complete ✅ | Handlers pending 🚧 | API compiles ✅
