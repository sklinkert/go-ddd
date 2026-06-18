.DEFAULT_GOAL := help

# Local database URL used by migrate-* targets. Override on the CLI:
#   make migrate-up DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable
DATABASE_URL ?= postgres://marketplace:marketplace@localhost:5432/marketplace?sslmode=disable

# Packages whose tests require Docker (testcontainers).
DOCKER_TEST_PKGS := github.com/sklinkert/go-ddd/internal/infrastructure/db/postgres \
                    github.com/sklinkert/go-ddd/internal/testhelpers

.PHONY: help build run test test-unit lint fmt tidy vendor sqlc migrate-up migrate-down docker-up docker-down

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary
	go build -o bin/marketplace ./cmd/marketplace

run: ## Run the server locally
	go run ./cmd/marketplace

test: ## Run all tests with the race detector (needs Docker)
	go test -race ./...

test-unit: ## Run only tests that don't need Docker
	go test -race $(filter-out $(DOCKER_TEST_PKGS),$(shell go list ./...))

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format code (gofmt + goimports)
	golangci-lint fmt ./...

tidy: ## Tidy and re-vendor modules
	go mod tidy
	go mod vendor

vendor: ## Re-vendor modules
	go mod vendor

sqlc: ## Regenerate sqlc code from sql/queries
	sqlc generate

migrate-up: ## Apply all pending migrations
	go run migrate.go -database-url "$(DATABASE_URL)" -command up

migrate-down: ## Roll back the last migration
	go run migrate.go -database-url "$(DATABASE_URL)" -command down -steps 1

docker-up: ## Start Postgres + app via docker compose
	docker compose up --build

docker-down: ## Stop and remove docker compose resources
	docker compose down -v
