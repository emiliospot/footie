# Architecture Guide - Clean Architecture vs Current Structure

## 🎯 Current vs Clean Architecture

### **Current Structure** (What We Have)

```
apps/api/internal/
├── api/                    # HTTP handlers & routes
├── domain/                 # Models (entities only)
├── infrastructure/         # DB, Redis, Logger implementations
├── repository/             # Data access layer
│   ├── interfaces.go       # Repository contracts
│   └── gorm/              # GORM implementation
├── config/                 # Configuration
└── pkg/                   # Reusable packages (auth, utils)
```

**Pros:**

- ✅ Simple, easy to navigate
- ✅ Repository pattern for DB abstraction
- ✅ Good for small-to-medium teams
- ✅ Familiar to most Go developers

**Cons:**

- ⚠️ Business logic mixed with HTTP handlers
- ⚠️ Not strictly following Clean/Hexagonal Architecture
- ⚠️ Testing requires mocking repositories directly

---

### **Clean Architecture** (Hexagonal/Ports & Adapters)

```
apps/api/internal/
├── domain/                          # CORE: Business entities + interfaces
│   ├── entities/
│   │   ├── user.go                 # Pure business entity
│   │   ├── team.go
│   │   └── match.go
│   └── repositories/                # Repository interfaces (ports)
│       ├── user_repository.go
│       ├── team_repository.go
│       └── match_repository.go
│
├── application/                     # USE CASES: Business logic
│   ├── usecases/
│   │   ├── user/
│   │   │   ├── create_user.go      # Single responsibility
│   │   │   ├── authenticate_user.go
│   │   │   └── get_user_profile.go
│   │   ├── team/
│   │   │   ├── create_team.go
│   │   │   └── get_team_statistics.go
│   │   └── match/
│   │       ├── create_match.go
│   │       └── analyze_match.go
│   └── services/                    # Domain services
│       ├── authentication_service.go
│       └── statistics_service.go
│
├── interfaces/                      # ADAPTERS: External interfaces
│   ├── http/                       # HTTP adapter
│   │   ├── handlers/
│   │   │   ├── user_handler.go
│   │   │   ├── team_handler.go
│   │   │   └── match_handler.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   └── logging.go
│   │   └── router.go
│   ├── grpc/                       # (Optional) gRPC adapter
│   │   └── user_grpc.go
│   └── cli/                        # (Optional) CLI adapter
│       └── commands.go
│
└── infrastructure/                  # IMPLEMENTATIONS: External concerns
    ├── persistence/                # Database implementations
    │   ├── gorm/
    │   │   ├── user_repository.go
    │   │   ├── team_repository.go
    │   │   └── transaction.go
    │   └── redis/
    │       └── cache_repository.go
    ├── external/                   # External services
    │   ├── email_service.go
    │   └── storage_service.go
    └── config/
        └── config.go
```

**Pros:**

- ✅ **True separation of concerns**
- ✅ **Business logic independent** of frameworks
- ✅ **Highly testable** (mock use cases, not repositories)
- ✅ **Easy to swap** HTTP ↔ gRPC ↔ CLI
- ✅ **Screaming architecture** (you can see what the app does)
- ✅ **Industry standard** for large applications

**Cons:**

- ⚠️ More complex for small projects
- ⚠️ More files and indirection
- ⚠️ Steeper learning curve for juniors

---

## 📊 Detailed Comparison

### **Example: Creating a User**

#### Current Structure

```go
// internal/api/handlers/user_handler.go (❌ Mixed concerns)
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.ShouldBindJSON(&req)

    // Business logic in handler!
    passwordHash, _ := auth.HashPassword(req.Password)
    user := &models.User{
        Email: req.Email,
        PasswordHash: passwordHash,
    }

    // Direct repository call
    err := h.repo.User().Create(c.Request.Context(), user)
    c.JSON(200, user)
}
```

#### Clean Architecture

```go
// 1. Domain Entity (internal/domain/entities/user.go)
type User struct {
    ID           uint
    Email        string
    PasswordHash string
    Role         string
}

// 2. Repository Interface (internal/domain/repositories/user_repository.go)
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
}

// 3. Use Case (internal/application/usecases/user/create_user.go)
type CreateUserUseCase struct {
    userRepo UserRepository
    hasher   PasswordHasher
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Pure business logic!

    // Check if user exists
    existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrUserAlreadyExists
    }

    // Hash password
    hash, err := uc.hasher.Hash(req.Password)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &User{
        Email:        req.Email,
        PasswordHash: hash,
        Role:         "user",
    }

    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

// 4. HTTP Handler (internal/interfaces/http/handlers/user_handler.go)
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Delegate to use case
    user, err := h.createUserUseCase.Execute(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, user)
}

// 5. GORM Implementation (internal/infrastructure/persistence/gorm/user_repository.go)
type GormUserRepository struct {
    db *gorm.DB
}

func (r *GormUserRepository) Save(ctx context.Context, user *entities.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

---

## 🎯 Which Should You Use?

### **Use Current Structure** if:

- ✅ Small to medium project (< 50k LOC)
- ✅ Team is < 5 developers
- ✅ Tight deadlines
- ✅ CRUD-heavy application
- ✅ Team less experienced with Clean Architecture

### **Use Clean Architecture** if:

- ✅ Large project (> 50k LOC)
- ✅ Multiple teams working on codebase
- ✅ Complex business logic
- ✅ Need to support multiple interfaces (HTTP, gRPC, CLI)
- ✅ Long-term maintainability critical
- ✅ **For senior/staff engineer interviews** (shows architectural maturity)

---

## 🔄 Migration Path (If You Want Clean Architecture)

### Phase 1: Extract Use Cases

```go
// Create application/usecases/
internal/application/usecases/
├── user/
│   ├── create_user.go
│   ├── authenticate_user.go
│   └── update_user_profile.go
└── team/
    ├── create_team.go
    └── calculate_team_stats.go
```

Move business logic from handlers into use cases.

### Phase 2: Move Entities

```go
// Rename domain/models → domain/entities
internal/domain/entities/
├── user.go
├── team.go
└── match.go
```

### Phase 3: Reorganize Repositories

```go
// Move repository interfaces to domain
internal/domain/repositories/
├── user_repository.go
└── team_repository.go

// Move implementations to infrastructure
internal/infrastructure/persistence/gorm/
├── user_repository.go
└── team_repository.go
```

### Phase 4: Refactor Handlers

```go
// Handlers become thin adapters
internal/interfaces/http/handlers/
└── user_handler.go  // Just HTTP → Use Case → HTTP
```

---

## 📁 Recommended Structure for Footie

Given this is a **football analytics platform**:

### **Hybrid Approach** (Best of Both Worlds)

```
apps/api/internal/
├── domain/                          # Business core
│   ├── entities/                   # Pure entities
│   │   ├── user.go
│   │   ├── team.go
│   │   ├── player.go
│   │   ├── match.go
│   │   └── statistics.go
│   ├── repositories/               # Repository interfaces
│   │   └── interfaces.go
│   └── services/                   # Domain services
│       ├── statistics_calculator.go
│       └── match_analyzer.go
│
├── application/                     # Use cases (for complex flows)
│   ├── auth/
│   │   ├── register_user.go
│   │   └── authenticate_user.go
│   └── analytics/
│       ├── generate_team_report.go
│       └── compare_players.go
│
├── interfaces/                      # HTTP layer
│   └── http/
│       ├── handlers/
│       ├── middleware/
│       └── router.go
│
└── infrastructure/                  # External implementations
    ├── persistence/
    │   └── gorm/
    ├── cache/
    │   └── redis/
    └── config/
```

**Why Hybrid?**

- Simple CRUD → Direct handler → repository
- Complex analytics → Handler → Use Case → Service → Repository
- Best of both worlds!

---

## 🎤 Interview Response

When asked about architecture:

> "I use a **hybrid approach** between repository pattern and clean architecture. For simple CRUD operations, I keep it straightforward with handlers calling repositories directly. But for **complex business logic**—like generating football analytics, comparing team statistics, or calculating player performance metrics—I extract that into **dedicated use cases** in the application layer.
>
> This gives us the **flexibility** of clean architecture where it matters, without the overhead on simple operations. The **repository pattern** ensures we can easily swap ORMs, and the **use case layer** keeps complex business logic testable and independent of HTTP concerns.
>
> For a football analytics platform, this is crucial because the **analytics calculations** are complex and evolving—we don't want that coupled to our HTTP handlers or database implementation."

---

## 🚀 Quick Decision Matrix

| Factor               | Current Structure | Clean Architecture | Hybrid   |
| -------------------- | ----------------- | ------------------ | -------- |
| **Simplicity**       | ⭐⭐⭐⭐⭐        | ⭐⭐               | ⭐⭐⭐⭐ |
| **Testability**      | ⭐⭐⭐            | ⭐⭐⭐⭐⭐         | ⭐⭐⭐⭐ |
| **Scalability**      | ⭐⭐⭐            | ⭐⭐⭐⭐⭐         | ⭐⭐⭐⭐ |
| **Learning Curve**   | ⭐⭐⭐⭐⭐        | ⭐⭐               | ⭐⭐⭐⭐ |
| **Interview Impact** | ⭐⭐⭐            | ⭐⭐⭐⭐⭐         | ⭐⭐⭐⭐ |

---

## 💡 My Recommendation

**For this project (Footie):**

1. **Keep current structure** for now ✅
2. **Add application/usecases/** for complex analytics ✅
3. **Migrate gradually** as business logic grows ✅
4. **Document the pattern** (this file!) ✅

**Why?**

- You have repository pattern (can swap ORMs) ✅
- You can demonstrate understanding of Clean Architecture 🎓
- You're pragmatic (not over-engineering) 💡
- You can evolve as needed 🔄

This shows **senior-level thinking**: knowing when to apply patterns vs when they're overkill! 🎯
