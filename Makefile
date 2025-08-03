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
test: ## Run tests (requires database)
	@echo "🧪 Running tests..."
	@if ! docker compose ps postgres 2>/dev/null | grep -q "Up"; then \
		echo "❌ PostgreSQL database is not running"; \
		echo "🚀 To start the database, run: make dev-setup"; \
		echo "📖 Or manually: docker compose up -d postgres"; \
		exit 1; \
	fi
	@if ! docker compose exec postgres pg_isready -U urluser -d url_shortener >/dev/null 2>&1; then \
		echo "❌ PostgreSQL database is not ready"; \
		echo "⏳ Please wait for the database to start up"; \
		echo "🔍 Check status with: docker compose logs postgres"; \
		exit 1; \
	fi
	@echo "📊 Database ready - running tests"
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