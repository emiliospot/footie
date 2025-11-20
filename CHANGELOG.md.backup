# Changelog

All notable changes to this project will be documented in this file. See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.

## [Unreleased]

### ✨ Features

- **monorepo**: Nx workspace for efficient monorepo management
- **backend**: RESTful API using Gin framework with PostgreSQL + Redis
- **backend**: JWT authentication with role-based access control
- **backend**: Repository pattern for database abstraction
- **backend**: Air hot-reload for development
- **backend**: Comprehensive testing suite (unit, integration, benchmarks)
- **frontend**: Angular 19 with standalone components
- **frontend**: Angular Material UI components
- **frontend**: HTTP interceptors and route guards
- **testing**: E2E tests with Playwright
- **infra**: AWS infrastructure with Terraform (VPC, ECS, RDS, ElastiCache, S3, CloudFront)
- **ci**: GitHub Actions workflows for CI/CD
- **ci**: Pre-commit hooks using Husky + lint-staged
- **ci**: Type checking in CI/CD pipeline

### 🐛 Bug Fixes

- **backend**: Go variable shadowing issues in error handling
- **backend**: Unchecked errors in strconv.Atoi calls
- **backend**: Integer overflow protection in type conversions
- **backend**: Type assertion safety checks
- **backend**: godotenv.Load error handling
- **backend**: Build tag format for Go integration tests
- **backend**: Huge parameter warning with pointer receiver
- **frontend**: Angular member ordering ESLint warnings
- **frontend**: Naming convention rules for snake_case in models
- **frontend**: Zone.js version for Angular 19 compatibility
- **docs**: Markdown linting configuration
- **infra**: Port conflicts for PostgreSQL (5436), Redis (6386), API (8081)
- **infra**: Docker Compose V2 syntax compatibility

### 📚 Documentation

- **docs**: Comprehensive README with setup instructions
- **docs**: QUICKSTART guide for new developers
- **docs**: TESTING_STRATEGY documentation
- **docs**: ARCHITECTURE comparison guide
- **docs**: DEPLOYMENT guide for AWS
- **docs**: Organized documentation into workspace/docs/ folder

### 🔧 Chores

- **backend**: Backend port changed from 8080 to 8081
- **infra**: Database port changed from 5432 to 5436 (local dev)
- **infra**: Redis port changed from 6379 to 6386 (local dev)
- **monorepo**: Project structure consolidated into workspace
- **frontend**: Angular upgraded to version 19
- **backend**: golangci-lint config to allow parallel runners
- **cleanup**: Removed deprecated documentation files
- **cleanup**: Removed old backend/frontend folders
- **cleanup**: Removed temporary migration scripts

---

## 🚀 How to Use This Changelog

This changelog is now **automatically generated** from your commit messages using [Conventional Commits](https://www.conventionalcommits.org/).

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature (appears in changelog)
- `fix`: Bug fix (appears in changelog)
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

### Examples

```bash
feat(backend): add match statistics endpoint
fix(frontend): resolve authentication token refresh issue
docs(readme): update installation instructions
chore(deps): upgrade Angular to v19
```

### Generate Changelog

```bash
# Add new entries since last release
npm run changelog

# Generate entire changelog from scratch
npm run changelog:first
```

The changelog will be automatically updated when you run `npm version` to bump the version.

---

## 📋 Pre-Commit Hooks

Every commit now runs:

- ✅ **lint-staged**: Auto-formats and lints only changed files
- ✅ **TypeScript typecheck**: Ensures no type errors
- ✅ **Go vet**: Validates Go code correctness

## 🚀 CI/CD Pipeline

On every push/PR:

1. ✅ Lint all affected projects
2. ✅ Type check entire codebase
3. ✅ Run unit tests with coverage
4. ✅ Run integration tests (Go + PostgreSQL)
5. ✅ Run E2E tests (Playwright)
6. ✅ Build all artifacts
7. ✅ Deploy to AWS (on main branch)
