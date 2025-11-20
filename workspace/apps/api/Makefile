.PHONY: help run build test lint swagger db-up db-down db-reset db-create docker-build

BINARY_NAME=footie-api
MAIN_PATH=cmd/api/main.go
DB_URL=postgres://footie_user:footie_password@localhost:5436/footie?sslmode=disable
MIGRATIONS_PATH=./migrations

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod verify
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

run: ## Run the application
	go run $(MAIN_PATH)

dev: ## Run with hot-reload
	air -c .air.toml

build: ## Build the binary
	CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(BINARY_NAME)-linux $(MAIN_PATH)

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

lint: ## Run linter
	golangci-lint run --config .golangci.yml

lint-fix: ## Run linter and fix issues
	golangci-lint run --config .golangci.yml --fix

swagger: ## Generate swagger documentation
	swag init -g $(MAIN_PATH) -o docs

# Database commands (using golang-migrate)
db-up: ## Run database migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

db-down: ## Rollback last migration
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

db-reset: ## Drop all tables and re-run migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop -f
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

db-create: ## Create new migration (use: make db-create name=add_xg_field)
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

db-status: ## Show migration status
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

db-force: ## Force migration version (use: make db-force version=1)
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(version)

# sqlc commands
sqlc-generate: ## Generate Go code from SQL
	sqlc generate

sqlc-vet: ## Vet SQL queries
	sqlc vet

seed: ## Seed database
	go run cmd/seed/main.go

docker-build: ## Build Docker image
	docker build -t footie-backend:latest .

docker-run: ## Run in Docker
	docker run -p 8080:8080 --env-file ../.env footie-backend:latest

clean: ## Clean build artifacts
	rm -rf bin/ tmp/ coverage.out docs/

tidy: ## Tidy go modules
	go mod tidy

fmt: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

security: ## Run security scan
	gosec ./...

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

.DEFAULT_GOAL := help
