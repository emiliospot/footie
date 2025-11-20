# ⚽ Footie - Football Analytics Platform

**Professional Nx monorepo** with Golang backend, Angular 19 frontend, and enterprise-grade architecture.

---

## 🚀 Quick Start

```bash
# Install dependencies
npm install

# Start infrastructure (PostgreSQL + Redis)
npm run docker:up

# Start development servers
npm run dev
```

**Access:**

- **Frontend:** http://localhost:4200
- **Backend API:** http://localhost:8081
- **API Docs:** http://localhost:8081/swagger

> For detailed setup, see [docs/QUICKSTART.md](docs/QUICKSTART.md)

---

## 📁 Project Structure

```
workspace/
├── apps/
│   ├── api/              # Golang backend with Air hot-reload
│   ├── web/              # Angular 19 frontend
│   └── web-e2e/          # Playwright E2E tests
├── libs/
│   └── shared/           # Shared TypeScript types
├── infra/
│   ├── docker/           # Docker Compose for local dev
│   └── terraform/        # AWS infrastructure as code
├── package.json          # Workspace root
└── nx.json               # Nx configuration
```

---

## 🛠️ Development Commands

```bash
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
npm run docker:logs   # View logs

# Nx Commands
npx nx graph          # Visualize dependencies
npx nx affected:test  # Test only affected code
npx nx affected:build # Build only affected code
```

---

## 🎯 Tech Stack

### Backend

- **Golang 1.21+** with Gin framework
- **Air** for hot-reload development (< 1s rebuild)
- **GORM** with repository pattern (easy to swap ORMs)
- **PostgreSQL 16** for data storage
- **Redis 7** for caching
- **testcontainers-go** for integration tests

### Frontend

- **Angular 19** with standalone components
- **TypeScript** (strict mode)
- **RxJS 7** for reactive programming
- **Angular Material** for UI components
- **Playwright** for E2E testing

### Infrastructure

- **Docker & Docker Compose** for local development
- **Nx 22.1** for monorepo management
- **Terraform** for AWS infrastructure as code
- **GitHub Actions** ready for CI/CD

---

## ✨ Key Features

- ⚡ **Air Hot-Reload** - Golang rebuilds in < 1 second
- 🚀 **Angular 19 HMR** - Instant frontend updates
- 🏗️ **Repository Pattern** - Abstracted database layer (easy ORM swap)
- 🧪 **Comprehensive Testing** - Unit, integration, and E2E tests
- 📦 **Nx Monorepo** - Build caching and affected commands
- 🐳 **Docker Ready** - Local development infrastructure
- ☁️ **AWS Deployment** - Terraform configurations included

---

## 🏗️ Architecture

We use a **hybrid approach** combining:

- **Repository pattern** for data access abstraction
- **Use cases** for complex business logic
- **Clean separation** of concerns
- **Interface-based design** for testability

The repository pattern makes it trivial to swap from GORM to sqlx, ent, or any other ORM without changing business logic.

For detailed architectural decisions, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

---

## 🧪 Testing

Comprehensive testing across all layers:

| Test Type           | Tool                  | Location                          |
| ------------------- | --------------------- | --------------------------------- |
| Backend Unit        | `testing` + `testify` | `apps/api/**/*_test.go`           |
| Backend Integration | `testcontainers-go`   | `apps/api/**/integration_test.go` |
| Backend Benchmarks  | `go test -bench`      | Embedded in unit tests            |
| Frontend Unit       | Jasmine + Karma       | `apps/web/**/*.spec.ts`           |
| E2E Tests           | Playwright            | `apps/web-e2e/src/**/*.spec.ts`   |

See [docs/TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md) for complete testing approach.

---

## 🚢 Deployment

### Local Development

```bash
npm run docker:up
npm run dev
```

### Production (AWS)

```bash
cd infra/terraform
terraform init
terraform apply
```

Infrastructure includes:

- VPC with public/private subnets
- ECS Fargate for containers
- RDS PostgreSQL
- ElastiCache Redis
- S3 + CloudFront
- Application Load Balancer

---

## 🔐 Security

- JWT authentication with refresh tokens
- Password hashing with bcrypt
- Role-based access control (Admin, User)
- CORS properly configured
- Rate limiting enabled
- SQL injection protection via GORM
- Environment-based configuration

---

## 📊 Database

- **PostgreSQL 16** as primary database
- **GORM** for ORM (abstracted via repository pattern)
- **Migrations** handled via SQL scripts
- **Redis 7** for caching and sessions
- Easy to swap ORMs (sqlx, ent, etc.)

---

## 🤝 Contributing

1. Create feature branch from `main`
2. Make changes with tests
3. Ensure linting passes (`npm run lint`)
4. Submit PR with description

---

## 📖 Documentation

- **[docs/QUICKSTART.md](docs/QUICKSTART.md)** - 3-minute setup guide
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Architecture decisions
- **[docs/TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md)** - Testing approach
- **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)** - AWS deployment guide
- **[infra/terraform/README.md](infra/terraform/README.md)** - Terraform setup

---

## 📄 License

MIT License

---

## 🆘 Support

For questions or issues:

1. Check [docs/QUICKSTART.md](docs/QUICKSTART.md) for setup help
2. See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for design decisions
3. Review [docs/TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md) for testing info

---

**Built with ❤️ for football analytics**

_Powered by Nx, Golang, and Angular 19_
