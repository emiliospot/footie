# 🚀 Quick Start - Footie Monorepo

## ⚡ 3-Minute Setup

```bash
# 1. Navigate to workspace
cd workspace

# 2. Install Node dependencies
npm install

# 3. Install Go dependencies
cd apps/api && go mod download && cd ../..

# 4. Install Air for hot-reload (if not already installed)
go install github.com/air-verse/air@latest

# 5. Setup environment
cp .env.example .env

# 6. Start infrastructure (PostgreSQL + Redis)
npm run docker:up

# 7. Start development servers (with hot-reload!)
npm run dev
```

**That's it!** 🎉

## 🌐 Access Your App

- **Frontend**: http://localhost:4200
- **Backend API**: http://localhost:8081
- **API Docs**: http://localhost:8081/swagger
- **PgAdmin**: http://localhost:5050 (optional, with `--profile tools`)

## 🧪 Verify Everything Works

```bash
# Run all tests
npm test

# Or run individually:
npm run test:api              # Backend unit tests
npm run test:api:integration  # Backend integration tests (Docker)
npm run test:web              # Frontend tests
npm run test:e2e              # Playwright E2E tests

# Check linting
npm run lint
```

## 🎮 Common Commands

```bash
# Development
npm run dev          # Start API + Web with hot-reload
npm run api          # Backend only (Air hot-reload)
npm run web          # Frontend only (HMR)

# Building
npm run build        # Build everything
npm run build:api    # Build backend
npm run build:web    # Build frontend

# Testing
npm test             # Test all
npm run test:api     # Backend tests
npm run test:web     # Frontend tests
npm run test:e2e     # E2E tests

# Linting
npm run lint         # Lint all
npm run lint:fix     # Auto-fix issues

# Docker
npm run docker:up    # Start infrastructure
npm run docker:down  # Stop infrastructure
npm run docker:logs  # View logs

# Nx Commands
npx nx graph                # Visualize dependencies
npx nx affected:test        # Test only changed code
npx nx affected:build       # Build only changed code
```

## 🔥 Hot-Reload in Action

### Backend (Golang with Air)

1. Edit any `.go` file in `apps/api/`
2. Watch the terminal - Air rebuilds automatically
3. Changes appear in < 1 second!

### Frontend (Angular with HMR)

1. Edit any `.ts`, `.html`, or `.scss` file in `apps/web/`
2. Browser updates instantly
3. No manual refresh needed!

## 📁 Project Structure

```
workspace/
├── apps/
│   ├── api/          # Golang backend
│   ├── web/          # Angular 19 frontend
│   └── web-e2e/      # Playwright E2E tests
├── libs/
│   └── shared/       # Shared TypeScript types
└── infra/
    ├── docker/       # Docker Compose
    └── terraform/    # AWS infrastructure
```

## 🎯 What's Included

✅ **Air** hot-reload for Golang
✅ **Angular 19** with HMR
✅ **Repository pattern** (abstracted DB)
✅ **PostgreSQL 16** + **Redis 7**
✅ **Complete testing** (unit, integration, E2E)
✅ **Nx monorepo** (caching, affected commands)
✅ **Docker** ready
✅ **AWS** deployment configured
✅ **CI/CD** with GitHub Actions

## 📚 Next Steps

- **Testing**: See [TESTING_STRATEGY.md](TESTING_STRATEGY.md)
- **Architecture**: See [ARCHITECTURE.md](ARCHITECTURE.md)
- **Deployment**: See [DEPLOYMENT.md](DEPLOYMENT.md)
- **Full Guide**: See [../README.md](../README.md)

## 🐛 Troubleshooting

### Air not found?

```bash
go install github.com/air-verse/air@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### Port already in use?

```bash
# Kill process on port 8088
lsof -ti:8088 | xargs kill -9

# Or change port in .env
API_PORT=8088
```

### Database connection error?

```bash
# Ensure PostgreSQL is running
docker ps | grep postgres

# Restart infrastructure
npm run docker:down
npm run docker:up
```

## 💡 Pro Tips

1. Keep Air terminal open - see rebuild status
2. Use `npx nx graph` to visualize dependencies
3. Use `nx affected:*` to only build/test changes
4. Create `.env` from `.env.example` for local config

## 🎉 You're Ready

Start building amazing football analytics features! ⚽

**Happy coding!** 🚀
