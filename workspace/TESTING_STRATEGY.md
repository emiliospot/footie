# Testing Strategy - Footie Monorepo

## 🎯 Overview

Comprehensive testing strategy covering **unit**, **integration**, **performance**, and **end-to-end** tests across the full stack.

## 📋 Testing Layers

### 1️⃣ **Backend (Golang) Testing**

#### **Tools & Frameworks**

- **`testing`** (standard library) - Core unit tests
- **`stretchr/testify`** - Assertions (`assert`/`require`) and mocking
- **`httptest`** - HTTP handler testing
- **`testcontainers-go`** - Real Postgres/Redis in Docker for integration tests
- **Benchmarks** - `go test -bench=.` for performance-critical code

#### **Test Types**

**Unit Tests** (`*_test.go`)

```go
// Test repository logic with in-memory SQLite
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t) // In-memory SQLite
    repo := gorm.NewUserRepository(db)

    // Table-driven tests
    tests := []struct {
        name    string
        user    *models.User
        wantErr bool
    }{
        {"valid user", &models.User{...}, false},
        {"duplicate email", &models.User{...}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := repo.Create(ctx, tt.user)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**Integration Tests** (`integration_test.go`)

```go
//go:build integration
// +build integration

func TestIntegration_WithRealPostgres(t *testing.T) {
    // Spin up real Postgres container
    db, cleanup := setupPostgresContainer(t)
    defer cleanup()

    // Test against real database
    repo := gorm.NewRepositoryManager(db)
    // ... full CRUD cycle tests
}
```

**Performance Tests** (benchmarks)

```go
func BenchmarkUserRepository_Create(b *testing.B) {
    db := setupTestDB(&testing.T{})
    repo := gorm.NewUserRepository(db)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        user := &models.User{...}
        _ = repo.Create(ctx, user)
    }
}
```

#### **Running Backend Tests**

```bash
# Unit tests (fast, in-memory)
nx run api:test
# Or: cd apps/api && go test ./...

# Integration tests (with Docker containers)
nx run api:test:integration
# Or: go test ./... -tags=integration

# Benchmarks
nx run api:bench
# Or: go test ./... -bench=. -benchmem

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

### 2️⃣ **Frontend (Angular) Testing**

#### **Tools & Frameworks**

- **Jasmine + Karma** (Angular default)
- **Angular Testing Library** (optional, better DX)
- **Jest** (modern alternative, faster)

#### **Test Types**

**Unit Tests** (Services, Pipes, Utils)

```typescript
describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [AuthService],
    });
    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should login successfully', () => {
    const mockResponse: AuthResponse = {...};

    service.login({email: 'test@test.com', password: 'pass'})
      .subscribe(response => {
        expect(response.token).toBeTruthy();
      });

    const req = httpMock.expectOne('/api/v1/auth/login');
    expect(req.request.method).toBe('POST');
    req.flush(mockResponse);
  });
});
```

**Component Tests**

```typescript
describe('TeamListComponent', () => {
  let component: TeamListComponent;
  let fixture: ComponentFixture<TeamListComponent>;
  let teamService: jasmine.SpyObj<TeamService>;

  beforeEach(() => {
    const spy = jasmine.createSpyObj('TeamService', ['getTeams']);

    TestBed.configureTestingModule({
      imports: [TeamListComponent],
      providers: [{ provide: TeamService, useValue: spy }],
    });

    fixture = TestBed.createComponent(TeamListComponent);
    component = fixture.componentInstance;
    teamService = TestBed.inject(TeamService) as jasmine.SpyObj<TeamService>;
  });

  it('should display teams', () => {
    const mockTeams: Team[] = [...];
    teamService.getTeams.and.returnValue(of({ data: mockTeams, pagination: {...} }));

    fixture.detectChanges();

    const compiled = fixture.nativeElement;
    expect(compiled.querySelectorAll('.team-card').length).toBe(mockTeams.length);
  });
});
```

#### **Running Frontend Tests**

```bash
# Unit & component tests
nx run web:test
# Or: cd apps/web && ng test

# Watch mode (development)
nx run web:test:watch

# CI mode (headless, single run)
nx run web:test --browsers=ChromeHeadless --watch=false
```

---

### 3️⃣ **End-to-End (E2E) Testing**

#### **Tool: Playwright** (recommended)

#### **Test Types**

**User Journeys** (Critical paths)

- Authentication (login, register, logout)
- Team management (list, filter, view details)
- Match analytics (view match, compare teams, export data)
- Player statistics (view stats, filter, sort)

**Example E2E Test**

```typescript
test("should view match analytics with filtering", async ({ page }) => {
  // Login
  await page.goto("/auth/login");
  await page.fill('[name="email"]', "analyst@example.com");
  await page.fill('[name="password"]', "password");
  await page.click('button[type="submit"]');

  // Navigate to matches
  await page.goto("/matches");

  // Filter by competition
  await page.selectOption('select[name="competition"]', "Premier League");
  await expect(page.locator(".match-card")).toHaveCount(10);

  // Open match detail
  await page.click(".match-card:first-child");
  await expect(page).toHaveURL(/\/matches\/\d+/);

  // Verify analytics are displayed
  await expect(page.locator('[data-testid="possession-chart"]')).toBeVisible();
  await expect(page.locator('[data-testid="shots-chart"]')).toBeVisible();
});
```

#### **Running E2E Tests**

```bash
# Run all E2E tests
nx run web-e2e:e2e

# UI mode (interactive)
nx run web-e2e:e2e:ui

# Debug mode
nx run web-e2e:e2e:debug

# Specific browser
nx run web-e2e:e2e --project=chromium
```

---

## 🔄 CI/CD Integration

### **GitHub Actions Workflow**

```yaml
name: Test Suite

on: [push, pull_request]

jobs:
  backend-unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.21"
      - name: Run unit tests
        run: nx run api:test
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage/api/coverage.out

  backend-integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - name: Run integration tests
        run: nx run api:test:integration

  backend-bench:
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - name: Run benchmarks
        run: nx run api:bench

  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - name: Install dependencies
        run: npm ci
      - name: Run tests
        run: nx run web:test
      - name: Upload coverage
        uses: codecov/codecov-action@v4

  e2e:
    runs-on: ubuntu-latest
    needs: [backend-unit, frontend-test]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - name: Install dependencies
        run: npm ci
      - name: Install Playwright
        run: npx playwright install --with-deps
      - name: Run E2E tests
        run: nx run web-e2e:e2e
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: playwright-report
          path: apps/web-e2e/playwright-report/
```

### **Nx Commands in CI**

```bash
# Test everything
nx run-many --target=test --all

# Test only affected (based on git diff)
nx affected --target=test

# Parallel execution
nx affected --target=test --parallel=3

# With caching (cached results reused)
nx affected --target=test --skip-nx-cache=false
```

---

## 📊 Coverage Goals

| Layer               | Target    | Current |
| ------------------- | --------- | ------- |
| Backend Unit        | 80%+      | -       |
| Backend Integration | Key paths | -       |
| Frontend Unit       | 70%+      | -       |
| Frontend Component  | 70%+      | -       |
| E2E Critical Paths  | 100%      | -       |

---

## 🎤 Interview-Ready Summary

> "In the monorepo, I use **table-driven tests** with Go's standard `testing` package plus `testify` for assertions. For anything touching the database, I use **`testcontainers-go`** so tests run against a real ephemeral Postgres instance, ensuring integration tests mirror production.
>
> On the frontend, I prefer **Jest with Angular Testing Library** for better DX, though I'm comfortable with Jasmine/Karma.
>
> For E2E, I use **Playwright** to cover critical user journeys—login, match selection, filtering analytics, and data export—running against a test environment with seeded match data.
>
> Everything is wired into **Nx targets**, so in CI we run:
>
> - `nx run api:test` → unit tests
> - `nx run api:test:integration` → integration with Docker
> - `nx run web:test` → Angular tests
> - `nx run web-e2e:e2e` → Playwright E2E
>
> Plus, we run **benchmarks** on performance-critical paths like statistics aggregation and match event processing to catch regressions early."

---

## 📁 Test File Organization

```
workspace/
├── apps/
│   ├── api/
│   │   ├── internal/
│   │   │   └── repository/
│   │   │       └── gorm/
│   │   │           ├── user_repository.go
│   │   │           ├── user_repository_test.go      # Unit tests
│   │   │           └── integration_test.go          # Integration tests
│   │   └── cmd/
│   │       └── api/
│   │           └── main_test.go                     # HTTP handler tests
│   ├── web/
│   │   └── src/
│   │       ├── app/
│   │       │   ├── core/
│   │       │   │   └── services/
│   │       │   │       ├── auth.service.ts
│   │       │   │       └── auth.service.spec.ts     # Unit tests
│   │       │   └── features/
│   │       │       └── teams/
│   │       │           ├── team-list.component.ts
│   │       │           └── team-list.component.spec.ts  # Component tests
│   │       └── test.ts                               # Test setup
│   └── web-e2e/
│       ├── src/
│       │   ├── auth.spec.ts                          # E2E: Auth
│       │   ├── teams.spec.ts                         # E2E: Teams
│       │   └── matches.spec.ts                       # E2E: Analytics
│       └── playwright.config.ts                      # Playwright config
├── coverage/                                          # Coverage reports
└── test-results/                                      # Test artifacts
```

---

## 🚀 Quick Commands

```bash
# Run all tests
npm test

# Test specific app
nx run api:test
nx run web:test
nx run web-e2e:e2e

# Test affected (only changed code)
nx affected:test

# Watch mode
nx run web:test:watch

# Coverage
nx run api:test --coverage
nx run web:test --code-coverage

# Benchmarks
nx run api:bench

# E2E with UI
nx run web-e2e:e2e:ui
```

---

This testing strategy ensures **high confidence** in code quality while maintaining **fast feedback loops** for developers. 🎯
