.PHONY: help run build test lint swagger migrate-up migrate-down migrate-create docker-build

BINARY_NAME=footie-api
MAIN_PATH=cmd/api/main.go

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod verify
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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

migrate-up: ## Run database migrations
	go run cmd/migrate/main.go up

migrate-down: ## Rollback last migration
	go run cmd/migrate/main.go down

migrate-create: ## Create new migration (use: make migrate-create name=migration_name)
	migrate create -ext sql -dir migrations -seq $(name)

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

