# Makefile for url-shortener development

.PHONY: help dev-setup dev-start dev-stop test test-integration build clean fmt vet install-hooks

# Default target
help: ## Show this help message
	@echo "url-shortener development commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development environment
dev-setup: ## Set up local development environment (PostgreSQL + .env)
	@./scripts/dev-setup.sh

dev-start: ## Start the application locally
	@echo "🚀 Starting url-shortener..."
	@go run main.go

dev-stop: ## Stop local development database
	@echo "🛑 Stopping development database..."
	@docker compose down

# Testing
test: ## Run tests with mock database
	@echo "🧪 Running tests..."
	@go test ./... -tags mock

test-integration: ## Run integration tests (requires running database)
	@echo "🔗 Running integration tests..."
	@go test ./...

# Code quality
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	@go fmt ./...

vet: ## Run static analysis
	@echo "🔍 Running static analysis..."
	@go vet ./...

# Build
build: ## Build the application binary
	@echo "🔨 Building application..."
	@go build -o bin/url-shortener

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/

# Git hooks
install-hooks: ## Install git pre-commit hooks
	@./scripts/install-hooks.sh

# Database management
db-shell: ## Connect to development database shell
	@docker compose exec postgres psql -U urluser -d url_shortener

db-logs: ## Show database logs
	@docker compose logs -f postgres

# Tools
pgadmin: ## Start pgAdmin (database management UI)
	@echo "🔧 Starting pgAdmin..."
	@docker compose --profile tools up -d pgadmin
	@echo "Visit http://localhost:8081 (admin@admin.com / admin)"