# Footie API (Golang Backend)

High-performance football analytics API built with Golang, featuring hot-reload with Air and abstracted database layer.

## Features

- ✅ **Hot Reload** with Air - changes reload automatically
- ✅ **Abstracted Database Layer** - easily switch ORMs (GORM, sqlx, ent, raw SQL)
- ✅ **Repository Pattern** - clean separation of concerns
- ✅ **Context-aware** - proper request context handling
- ✅ **Transaction Support** - built into repository manager
- ✅ **GORM Implementation** - default ORM (can swap easily)

## Architecture

### Repository Pattern

```
internal/
├── repository/
│   ├── interfaces.go           # Repository interfaces
│   └── gorm/                   # GORM implementation
│       ├── repository_manager.go
│       ├── user_repository.go
│       ├── team_repository.go
│       └── ...
```

### Benefits

1. **Testable**: Mock repositories easily
2. **Flexible**: Switch ORMs without changing business logic
3. **Clean**: Separation of concerns
4. **Maintainable**: Single responsibility principle

## Getting Started

### Prerequisites

- Go 1.21+
- Air for hot-reload: `go install github.com/air-verse/air@latest`
- PostgreSQL
- Redis

### Installation

```bash
# Install dependencies
go mod download

# Install Air (if not already installed)
go install github.com/air-verse/air@latest
```

### Development with Hot Reload

```bash
# Start with Air (recommended for development)
air

# Or run without hot-reload
go run cmd/api/main.go
```

Air will watch for changes and automatically rebuild/restart the server.

## Usage

### Using the Repository Pattern

```go
// Initialize repository manager
db, _ := database.NewPostgresDB(&cfg.Database)
repoManager := gorm.NewRepositoryManager(db)

// Use repositories
user, err := repoManager.User().FindByEmail(ctx, "user@example.com")

// With transactions
txManager, _ := repoManager.BeginTx(ctx)
defer txManager.Rollback() // rollback if not committed

err = txManager.User().Create(ctx, user)
err = txManager.TeamStatistics().Create(ctx, stats)

txManager.Commit() // commit transaction
```

### Switching ORMs

To switch from GORM to another ORM:

1. Create new implementation in `internal/repository/your-orm/`
2. Implement all repository interfaces
3. Update initialization in main.go
4. Business logic remains unchanged!

Example with sqlx:

```go
// internal/repository/sqlx/user_repository.go
type SqlxUserRepository struct {
    db *sqlx.DB
}

func (r *SqlxUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
    return &user, err
}
```

## Configuration

### Air Configuration (.air.toml)

```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api/main.go"
  bin = "./tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor", "../../node_modules"]
```

### Environment Variables

See main workspace `.env` file.

## API Endpoints

- `GET /health` - Health check
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/teams` - List teams
- `GET /api/v1/players` - List players
- `GET /api/v1/matches` - List matches

Full API documentation: http://localhost:8080/swagger

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/repository/gorm/...
```

### Testing with Mocks

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.User), args.Error(1)
}

// In tests
mockRepo := new(MockUserRepository)
mockRepo.On("FindByID", mock.Anything, uint(1)).Return(&user, nil)
```

## Database Migrations

```bash
# Run migrations
go run cmd/migrate/main.go up

# Rollback
go run cmd/migrate/main.go down
```

## Building

```bash
# Development build
go build -o bin/api cmd/api/main.go

# Production build (optimized)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go
```

## Performance

- Hot-reload rebuilds typically complete in < 1 second
- Context-aware queries prevent resource leaks
- Repository pattern allows for caching layers
- Transaction support ensures data consistency

## Why Repository Pattern?

### Without Repository Pattern

```go
// Business logic tightly coupled to GORM
func (s *UserService) GetUser(id uint) (*User, error) {
    var user User
    err := s.db.First(&user, id).Error
    return &user, err
}
```

### With Repository Pattern

```go
// Business logic independent of ORM
func (s *UserService) GetUser(ctx context.Context, id uint) (*User, error) {
    return s.repo.User().FindByID(ctx, id)
}
```

Benefits:

- ✅ Easy to test (mock repositories)
- ✅ Easy to switch ORMs
- ✅ Easy to add caching
- ✅ Clean architecture
- ✅ Single responsibility

## Learn More

- [Air Documentation](https://github.com/air-verse/air)
- [GORM Documentation](https://gorm.io)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
