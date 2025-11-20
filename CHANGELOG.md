# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Infrastructure & DevOps

- ✅ **Pre-commit hooks** using Husky for automated quality checks
- ✅ **lint-staged** integration for fast, staged-file-only linting
- ✅ **Type checking** in CI/CD pipeline for both TypeScript and Go
- ✅ **Comprehensive VSCode configuration**
  - Settings for Go, Angular, Prettier, ESLint, Markdownlint
  - Recommended extensions list
  - Custom tasks for development workflow
  - Air schema validation for `.air.toml`
- ✅ **GitHub Actions CI/CD workflows**
  - Main CI pipeline with linting, type checking, and testing
  - PR validation workflow
  - AWS ECS deployment workflow
  - Dependabot for automated dependency updates

#### Monorepo Architecture

- ✅ **Nx workspace** for efficient monorepo management
- ✅ **Project structure consolidation** into single workspace
- ✅ **Unified build and test orchestration**
- ✅ **Dependency graph visualization**

#### Backend (Go API)

- ✅ **RESTful API** using Gin framework
- ✅ **PostgreSQL integration** with GORM
- ✅ **Redis caching** support
- ✅ **JWT authentication** with role-based access control
- ✅ **Repository pattern** for database abstraction
- ✅ **Air hot-reload** for development
- ✅ **Comprehensive testing suite**
  - Unit tests with testify
  - Integration tests with testcontainers-go
  - Benchmark tests for performance
- ✅ **Strict linting** with golangci-lint
- ✅ **Type safety** with go vet

#### Frontend (Angular 19)

- ✅ **Angular 19** with standalone components
- ✅ **Strict TypeScript** configuration
- ✅ **Angular Material** UI components
- ✅ **HTTP interceptors** for auth and error handling
- ✅ **Route guards** for authentication
- ✅ **Lazy-loaded routes** for performance
- ✅ **Comprehensive ESLint** configuration
- ✅ **Type checking** with TypeScript compiler

#### Testing

- ✅ **Backend unit tests** using Go testing stdlib + testify
- ✅ **Backend integration tests** with real PostgreSQL containers
- ✅ **Backend benchmarks** for performance-critical code
- ✅ **Frontend unit tests** with Jasmine + Karma
- ✅ **E2E tests** with Playwright
- ✅ **Test coverage reporting**

#### Documentation

- ✅ **Comprehensive README** with setup instructions
- ✅ **QUICKSTART guide** for new developers
- ✅ **TESTING_STRATEGY** documentation
- ✅ **ARCHITECTURE** comparison guide
- ✅ **Technology stack** documentation

#### AWS Infrastructure (Terraform)

- ✅ **VPC with public/private subnets**
- ✅ **ECS Fargate** for container orchestration
- ✅ **RDS PostgreSQL** for production database
- ✅ **ElastiCache Redis** for caching
- ✅ **S3 + CloudFront** for frontend hosting
- ✅ **Application Load Balancer** for traffic distribution
- ✅ **Security groups** and IAM roles

### Fixed

- ✅ **Go variable shadowing** issues in error handling
- ✅ **Unchecked errors** in strconv.Atoi calls
- ✅ **Integer overflow** protection in type conversions
- ✅ **Type assertion safety** checks
- ✅ **Angular member ordering** ESLint warnings
- ✅ **Naming convention** rules for snake_case in models
- ✅ **Markdown linting** configuration (MD034, MD036, MD040)
- ✅ **Port conflicts** for PostgreSQL (5436), Redis (6386), API (8081)
- ✅ **Docker Compose V2** syntax compatibility
- ✅ **Zone.js version** for Angular 19 compatibility
- ✅ **Missing imports** in Go packages
- ✅ **godotenv.Load** error handling
- ✅ **Build tag format** for Go integration tests
- ✅ **Huge parameter** warning with pointer receiver

### Changed

- ✅ **Backend port** from 8080 to 8081 to avoid conflicts
- ✅ **Database port** from 5432 to 5436 (local development)
- ✅ **Redis port** from 6379 to 6386 (local development)
- ✅ **Project structure** consolidated into workspace monorepo
- ✅ **Angular version** upgraded to 19
- ✅ **golangci-lint config** to allow parallel runners

### Removed

- ✅ **Deprecated documentation** files
- ✅ **Old backend/frontend** folders (consolidated into workspace)
- ✅ **Temporary migration** scripts
- ✅ **Unused Makefile** at root level
- ✅ **cSpell** extension (too noisy for football terminology)

## [1.0.0] - 2025-11-20

### Added

- 🎉 Initial project scaffold
- 🎉 Nx monorepo setup
- 🎉 Go backend with Gin + GORM
- 🎉 Angular 19 frontend
- 🎉 PostgreSQL + Redis with Docker
- 🎉 AWS infrastructure with Terraform
- 🎉 GitHub Actions CI/CD
- 🎉 Comprehensive testing strategy
- 🎉 Pre-commit hooks for quality assurance

---

## 📋 **Pre-Commit Hooks**

Every commit now runs:

- ✅ **lint-staged**: Auto-formats and lints only changed files
- ✅ **TypeScript typecheck**: Ensures no type errors
- ✅ **Go vet**: Validates Go code correctness

## 🚀 **CI/CD Pipeline**

On every push/PR:

1. ✅ Lint all affected projects
2. ✅ Type check entire codebase
3. ✅ Run unit tests with coverage
4. ✅ Run integration tests (Go + PostgreSQL)
5. ✅ Run E2E tests (Playwright)
6. ✅ Build all artifacts
7. ✅ Deploy to AWS (on main branch)

---

[Unreleased]: https://github.com/emiliospot/footie/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/emiliospot/footie/releases/tag/v1.0.0
