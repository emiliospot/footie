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
- Backend API: http://localhost:8081
- API Docs: http://localhost:8081/swagger

---

## 📁 Project Structure

```
footie/
├── workspace/              # Nx monorepo (all code here!)
│   ├── apps/
│   │   ├── api/           # Golang backend with Air hot-reload
│   │   ├── web/           # Angular 19 frontend
│   │   └── web-e2e/       # Playwright E2E tests
│   ├── libs/
│   │   └── shared/        # Shared TypeScript types
│   ├── infra/
│   │   ├── docker/        # Docker Compose for local dev
│   │   └── terraform/     # AWS infrastructure as code
│   └── [documentation]
├── .git/                  # Git repository
├── .gitignore
└── README.md              # You are here!
```

---

## ✨ Key Features

- ⚡ **Air Hot-Reload** for Golang (< 1s rebuild)
- 🚀 **Angular 19** with HMR
- 🏗️ **Repository Pattern** (abstracted database layer)
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
- **GORM** with repository pattern (easy to swap ORMs)
- **PostgreSQL 16** for data storage
- **Redis 7** for caching
- **testcontainers-go** for integration tests

### Frontend

- **Angular 19** with standalone components
- **TypeScript** (strict mode)
- **RxJS 7** for reactive programming
- **Angular Material** for UI
- **Playwright** for E2E testing

### Infrastructure

- **Docker & Docker Compose**
- **Nx** for monorepo management
- **Terraform** for AWS IaC
- **GitHub Actions** for CI/CD

---

## 📚 Documentation

**All documentation is in the `workspace/docs/` directory:**

- **[workspace/README.md](workspace/README.md)** - Complete monorepo guide
- **[workspace/docs/QUICKSTART.md](workspace/docs/QUICKSTART.md)** - 3-minute setup
- **[workspace/docs/ARCHITECTURE.md](workspace/docs/ARCHITECTURE.md)** - Architecture decisions
- **[workspace/docs/TESTING_STRATEGY.md](workspace/docs/TESTING_STRATEGY.md)** - Testing approach
- **[workspace/docs/DEPLOYMENT.md](workspace/docs/DEPLOYMENT.md)** - AWS deployment guide

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

# Docker
npm run docker:up     # Start PostgreSQL & Redis
npm run docker:down   # Stop infrastructure

# Nx Commands
npx nx graph          # Visualize dependencies
npx nx affected:test  # Test only affected code
```

---

## 🏗️ Architecture

We use a **hybrid approach**:

- **Repository pattern** for data access abstraction
- **Use cases** for complex business logic
- **Clean separation** of concerns
- **Easy to test** and maintain

The repository pattern makes it trivial to swap from GORM to sqlx, ent, or any other ORM.

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
