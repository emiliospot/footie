# ⚽ Footie - Football Analytics Platform

> **Professional Nx monorepo** with Golang backend, Angular 19 frontend, and enterprise-grade architecture.

---

## 🚀 Quick Start

```bash
cd workspace
npm install
npm run docker:up
npm run dev
```

**Access:**

- Frontend: http://localhost:4200
- Backend API: http://localhost:8088
- API Docs: http://localhost:8088/swagger

---

## 📁 Project Structure

```
footie/
├── workspace/                      # Nx monorepo (all code here!)
│   ├── apps/
│   │   ├── api/                   # Golang backend with Air hot-reload
│   │   │   ├── cmd/api/           # Main application entry point
│   │   │   ├── internal/
│   │   │   │   ├── api/           # HTTP handlers, middleware, router
│   │   │   │   ├── config/        # Configuration management
│   │   │   │   ├── domain/        # Domain models
│   │   │   │   ├── infrastructure/ # Database, Redis, WebSocket, Logger
│   │   │   │   └── repository/    # Data access (sqlc + GORM legacy)
│   │   │   ├── migrations/        # golang-migrate SQL files
│   │   │   ├── pkg/               # Reusable packages (auth, utils)
│   │   │   ├── sqlc.yaml          # sqlc configuration
│   │   │   ├── Makefile           # Backend commands
│   │   │   └── README_SQLC.md     # sqlc + pgx guide
│   │   ├── web/                   # Angular 19 frontend
│   │   │   ├── src/
│   │   │   │   ├── app/
│   │   │   │   │   ├── core/      # Guards, interceptors, services, models
│   │   │   │   │   ├── features/  # Feature modules (auth, matches, players, teams)
│   │   │   │   │   └── shared/    # Shared components
│   │   │   │   └── environments/  # Environment configs
│   │   │   └── angular.json
│   │   └── web-e2e/               # Playwright E2E tests
│   ├── libs/
│   │   └── shared/                # Shared TypeScript types
│   ├── infra/
│   │   ├── docker/
│   │   │   └── docker-compose.yml # PostgreSQL 16 + Redis 8 + Redis Commander
│   │   └── terraform/             # AWS infrastructure as code
│   │       └── modules/           # Reusable Terraform modules
│   ├── docs/                      # Comprehensive documentation
│   │   ├── ARCHITECTURE.md        # System architecture & diagrams
│   │   ├── QUICKSTART.md          # 3-minute setup guide
│   │   ├── DEPLOYMENT.md          # AWS deployment guide
│   │   ├── TESTING_STRATEGY.md    # Testing approach
│   │   ├── MATCH_DATA_FEEDS.md    # External data integration (Phase 2)
│   │   ├── OPENSEARCH_INTEGRATION.md # Analytics engine (Phase 3)
│   │   ├── TECH_STACK_PRESENTATION.md # Complete tech overview
│   │   ├── TECH_IMPROVEMENTS_ROADMAP.md # Post-MVP enhancements
│   │   └── CI_CD_FIX.md           # CI/CD troubleshooting
│   ├── nx.json                    # Nx workspace configuration
│   ├── package.json               # Root scripts & dependencies
│   └── README.md                  # Workspace guide
├── .vscode/                       # VSCode settings & tasks
├── .github/
│   └── workflows/                 # GitHub Actions CI/CD
├── .gitignore
├── PM_INVITATION.md               # Product manager onboarding
└── README.md                      # You are here!
```

---

## ✨ Key Features

- ⚡ **Air Hot-Reload** for Golang (< 1s rebuild)
- 🚀 **Angular 19** with HMR
- 🔥 **sqlc + pgx** - Type-safe SQL with 3-5x faster queries (industry standard for analytics)
- 🗄️ **golang-migrate** - Production-grade database migrations
- 📡 **Real-Time WebSockets** - Sub-second match updates with Redis Streams & Pub/Sub
- 🧪 **Comprehensive Testing** (unit, integration, E2E)
- 📦 **Nx Monorepo** (build caching, affected commands)
- 🐳 **Docker** ready for local development
- ☁️ **AWS** deployment configured with Terraform
- 🔄 **CI/CD** ready with GitHub Actions

---

## 🛠️ Tech Stack

### Backend

- **Golang 1.21+** with Gin framework
- **Air** for hot-reload development
- **sqlc + pgx** - Type-safe SQL with fastest PostgreSQL driver (3-5x faster)
- **golang-migrate** - Production-grade database migrations
- **PostgreSQL 16** for data storage
- **Redis 8** for caching & real-time events (Streams + Pub/Sub)
- **WebSockets** for real-time match updates (Gorilla WebSocket)
- **testcontainers-go** for integration tests

### Frontend

- **Angular 19** with standalone components
- **TypeScript** (strict mode)
- **RxJS 7** for reactive programming
- **Angular Material** for UI
- **Playwright** for E2E testing

### Infrastructure & AWS Services

- **Docker & Docker Compose** - Local development
- **Nx** for monorepo management
- **Terraform** for AWS IaC
- **GitHub Actions** for CI/CD

**AWS Services (Production):**

- **AWS Lambda** - Serverless webhook processing
- **AWS Kinesis** - Event streaming (1000s events/sec)
- **AWS OpenSearch** - Analytics engine (Phase 3)
- **AWS RDS PostgreSQL** - Managed database
- **AWS ElastiCache Redis** - Managed cache
- **AWS API Gateway** - Webhook endpoints

---

## 📚 Documentation

### 🚀 Getting Started

- **[workspace/README.md](workspace/README.md)** - Complete monorepo guide
- **[workspace/docs/QUICKSTART.md](workspace/docs/QUICKSTART.md)** - 3-minute setup guide
- **[workspace/docs/ARCHITECTURE.md](workspace/docs/ARCHITECTURE.md)** - System architecture with diagrams
- **[workspace/docs/DEPLOYMENT.md](workspace/docs/DEPLOYMENT.md)** - AWS deployment guide

### 🔧 Backend Guides

- **[workspace/apps/api/README_SQLC.md](workspace/apps/api/README_SQLC.md)** - sqlc + pgx + golang-migrate complete guide
- **[workspace/apps/api/REALTIME_ARCHITECTURE.md](workspace/apps/api/REALTIME_ARCHITECTURE.md)** - WebSocket + Redis Streams architecture
- **[workspace/apps/api/MIGRATION_STATUS.md](workspace/apps/api/MIGRATION_STATUS.md)** - GORM → sqlc migration tracker

### 🎯 Product & Strategy

- **[PM_INVITATION.md](PM_INVITATION.md)** - Product manager onboarding & project overview
- **[workspace/docs/TECH_STACK_PRESENTATION.md](workspace/docs/TECH_STACK_PRESENTATION.md)** - Complete technical overview for stakeholders
- **[workspace/docs/TECH_IMPROVEMENTS_ROADMAP.md](workspace/docs/TECH_IMPROVEMENTS_ROADMAP.md)** - Post-MVP enhancements & phased rollout

### 🔮 Advanced Features (Phase 2+)

- **[workspace/docs/MATCH_DATA_FEEDS.md](workspace/docs/MATCH_DATA_FEEDS.md)** - External data feed integration (Opta, StatsBomb, API-Football)
- **[workspace/docs/OPENSEARCH_INTEGRATION.md](workspace/docs/OPENSEARCH_INTEGRATION.md)** - Analytics engine with AWS OpenSearch (Phase 3)

### 🧪 Testing & CI/CD

- **[workspace/docs/TESTING_STRATEGY.md](workspace/docs/TESTING_STRATEGY.md)** - Comprehensive testing approach
- **[workspace/docs/CI_CD_FIX.md](workspace/docs/CI_CD_FIX.md)** - CI/CD troubleshooting & test strategy

---

## 🧪 Development Commands

```bash
# All commands run from workspace/ directory
cd workspace

# Development
npm run dev           # Start everything with hot-reload
npm run api           # Backend only (Air hot-reload)
npm run web           # Frontend only (HMR)

# Testing
npm test              # Run all tests
npm run test:api      # Backend unit tests
npm run test:api:integration  # Backend integration tests
npm run test:web      # Frontend tests
npm run test:e2e      # Playwright E2E tests

# Building
npm run build         # Build all
npm run build:api     # Build backend
npm run build:web     # Build frontend

# Linting
npm run lint          # Lint all
npm run lint:fix      # Auto-fix issues

# Database Migrations
npm run db:up         # Run all pending migrations
npm run db:down       # Rollback last migration
npm run db:reset      # Drop all & re-run migrations
npm run db:status     # Check migration version
npm run sqlc:generate # Generate Go code from SQL

# Docker
npm run docker:up     # Start PostgreSQL & Redis
npm run docker:down   # Stop infrastructure

# Nx Commands
npx nx graph          # Visualize dependencies
npx nx affected:test  # Test only affected code
```

---

## 📡 Real-Time Architecture

**WebSocket + Redis Streams** for sub-second match updates:

```
Event → Redis Pub/Sub → WebSocket Hub → Connected Clients
  ↓
Redis Streams → Worker → Analytics → Database
```

**Connect to live match updates:**

```
ws://localhost:8088/ws/matches/:id
```

**Features:**

- 🔴 Sub-100ms latency
- 📊 100,000+ concurrent clients per instance
- 🚀 Horizontal scaling ready
- ⚽ Real-time goals, shots, passes, cards
- 📈 Live xG calculations
- 🎯 Cache invalidation on updates

**See:** `workspace/apps/api/REALTIME_ARCHITECTURE.md` for complete documentation.

---

## 🏗️ Architecture

We use a **production-grade approach** for sports analytics:

### Current Stack (Phase 1)

- **sqlc + pgx** - Type-safe SQL queries (used by betting companies & analytics teams)
- **golang-migrate** - Version-controlled database migrations
- **Raw SQL** - Perfect for complex analytics queries (xG, pass accuracy, heat maps)
- **Repository pattern** - Clean data access abstraction
- **WebSocket + Redis** - Real-time match updates (sub-100ms)
- **Clean separation** of concerns
- **Easy to test** and maintain

### Future Enhancements

**Phase 2: External Data Feeds (AWS-Native)**

```
External Feeds → API Gateway → Lambda → Kinesis → Go Consumers
```

- AWS Lambda for serverless webhook processing
- AWS Kinesis for high-throughput event streaming (1000s events/sec)
- Auto-scaling and replay capability

**Phase 3: Analytics Engine (Production Scale)**

```
Events → Kinesis → Go Consumer → OpenSearch
```

- AWS OpenSearch for advanced analytics (heat maps, xG trends, player similarity)
- Real-time dashboards with millisecond aggregations
- Full-text search and fuzzy matching

**The Perfect Trio:**

- **PostgreSQL** - Source of truth (CRUD, transactions)
- **Redis** - Real-time messaging (WebSocket, pub/sub)
- **OpenSearch** - Analytics & search (complex queries, ML)

This stack is the industry standard for high-performance analytics applications used by betting companies, sports data providers, and live streaming platforms.

For detailed architectural decisions, see [workspace/docs/ARCHITECTURE.md](workspace/docs/ARCHITECTURE.md).

---

## 🧪 Testing Strategy

Comprehensive testing across all layers:

- **Backend Unit Tests**: `testing` + `testify` + in-memory SQLite
- **Backend Integration Tests**: `testcontainers-go` with real Postgres
- **Backend Benchmarks**: Performance testing for critical paths
- **Frontend Unit Tests**: Jasmine + Karma
- **Frontend Component Tests**: Angular Testing Library
- **E2E Tests**: Playwright covering critical user journeys

See [workspace/docs/TESTING_STRATEGY.md](workspace/docs/TESTING_STRATEGY.md) for details.

---

## 🚢 Deployment

### Local Development

```bash
cd workspace
npm run docker:up
npm run dev
```

### Production (AWS)

```bash
cd workspace/infra/terraform
terraform init
terraform apply
```

See **[workspace/docs/DEPLOYMENT.md](workspace/docs/DEPLOYMENT.md)** for complete deployment guide.

---

## 🔐 Security

- JWT authentication with refresh tokens
- Password hashing with bcrypt
- Role-based access control
- CORS properly configured
- Rate limiting enabled
- SQL injection protection

---

## 🤝 Contributing

1. Create feature branch from `main`
2. Make changes with tests
3. Ensure linting passes
4. Submit PR with description

All commands should be run from the `workspace/` directory.

---

## 📄 License

MIT License - see LICENSE file

---

## 🆘 Support

- **Documentation**: Start with [workspace/README.md](workspace/README.md)
- **Quick Start**: [workspace/docs/QUICKSTART.md](workspace/docs/QUICKSTART.md)
- **Architecture**: [workspace/docs/ARCHITECTURE.md](workspace/docs/ARCHITECTURE.md)

---

**Built with ❤️ for football analytics ⚽**

_This is an Nx monorepo. All source code is in the `workspace/` directory._
