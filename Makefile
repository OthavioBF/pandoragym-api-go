.PHONY: help build run test clean deps fmt lint docker-setup docker-run docker-stop docker-restart docker-migrate docker-seed docker-logs docker-status docker-shell docker-clean db-start db-stop db-shell migrate-up migrate-down migrate-create seed setup dev-start dev-stop install-tools

# Default target
.DEFAULT_GOAL := help

# Build the application
build:
	@echo "🔨 Building application..."
	go build -o bin/pandoragym-api cmd/server/main.go
	@echo "✅ Build completed!"

# Run the application locally
run:
	@echo "🚀 Starting application locally..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	@echo "✅ Clean completed!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run
	@echo "✅ Linting completed!"

# =============================================================================
# Docker Commands
# =============================================================================

docker-setup:
	@echo "🚀 Setting up complete Docker environment..."
	@make docker-run
	@echo "⏳ Waiting for database to be ready..."
	@sleep 15
	@make docker-migrate
	@make docker-seed
	@echo "🎉 Complete Docker environment ready!"
	@echo "🌐 App available at http://localhost:3333"

docker-run:
	@echo "🚀 Starting Docker containers"
	@if ! docker-compose up -d; then \
		echo "❌ Docker failed to start trying to recover..."; \
		echo "🛑 Running docker-stop to clean up..."; \
		make docker-stop > /dev/null 2>&1 || true; \
		echo "🔄 Retrying to start containers..."; \
		docker-compose up -d; \
	fi
	@echo "⏳ Waiting for services to be ready..."
	@sleep 5
	@echo "🔍 Checking if API is running..."
	@for i in $$(seq 1 5); do \
		if curl -s -f http://localhost:3333/health > /dev/null 2>&1; then \
			echo "✅ API is running at http://localhost:3333"; \
			exit 0; \
		fi; \
		echo "⏳ Waiting for API... ($$i/5)"; \
		sleep 3; \
	done; \
	echo "❌ API failed to start"; \
	echo "📋 Application logs:"; \
	docker-compose logs --tail=20 app; \
	echo "🛑 Stopping containers due to API failure..."; \
	docker-compose down; \
	exit 1

docker-stop:
	@echo "🛑 Stopping Docker containers..."
	docker-compose down

docker-clean:
	@echo "🧹 Cleaning up Docker resources..."
	@echo "🛑 Stopping containers..."
	docker-compose down
	@echo "🗑️  Removing orphaned volumes..."
	docker volume prune -f
	@echo "🗑️  Removing unused networks..."
	docker network prune -f
	@echo "🗑️  Removing unused images..."
	docker image prune -f
	@echo "✅ Docker cleanup complete!"
	@echo "✅ Containers stopped!"

docker-restart:
	@echo "🔄 Restarting Docker containers..."
	@make docker-stop
	@make docker-run

docker-migrate:
	@echo "📊 Running migrations in Docker..."
	docker-compose exec -e DATABASE_HOST=db -e DATABASE_PORT=5432 -e DATABASE_NAME=pandoragym_db -e DATABASE_USER=pandoragym -e DATABASE_PASSWORD=password app go run ./cmd/tern
	@echo "✅ Migrations completed!"

docker-seed:
	@echo "🌱 Running seed in Docker..."
	docker-compose exec app go run ./cmd/seed
	@echo "✅ Seed completed!"

docker-logs:
	@echo "📋 Showing application logs..."
	docker-compose logs -f app

docker-status:
	@echo "📊 Docker Services Status:"
	@echo "=========================="
	@docker-compose ps

# Connect to app container shell
docker-shell:
	@echo "🐚 Connecting to application container..."
	docker-compose exec app sh

# =============================================================================
# Database Commands (for local development)
# =============================================================================

# Start database container only
db-start:
	@echo "📊 Starting PostgreSQL container..."
	docker run -d \
		--name pandoragym_db \
		-e POSTGRES_DB=pandoragym_db \
		-e POSTGRES_USER=pandoragym \
		-e POSTGRES_PASSWORD=password \
		-p 5432:5432 \
		postgres:15-alpine || echo "Container already exists"
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@echo "✅ Database ready!"

# Stop database container
db-stop:
	@echo "🛑 Stopping PostgreSQL container..."
	docker stop pandoragym_db || true
	docker rm pandoragym_db || true
	@echo "✅ Database stopped!"

# Connect to database shell
db-shell:
	@echo "🐚 Connecting to PostgreSQL shell..."
	docker exec -it pandoragym_db psql -U pandoragym -d pandoragym_db

# =============================================================================
# Migration and Seed Commands (for local development)
# =============================================================================

# Run migrations locally
migrate-up:
	@echo "📊 Running migrations locally..."
	go run ./cmd/tern
	@echo "✅ Migrations completed!"

# Rollback migrations
migrate-down:
	@echo "⬇️ Rolling back migrations..."
	tern migrate --migrations ./internal/infra/pgstore/migrations --config ./internal/infra/pgstore/migrations/tern.conf --destination 0
	@echo "✅ Migrations rolled back!"

# Create new migration
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	@echo "📝 Creating new migration: $(name)"
	tern new --migrations ./internal/infra/pgstore/migrations $(name)
	@echo "✅ Migration created!"

# Run seed locally
seed:
	@echo "🌱 Running seed locally..."
	go run ./cmd/seed
	@echo "✅ Seed completed!"

# Complete setup (migrations + seed) locally
setup: migrate-up seed

# =============================================================================
# Development Workflows
# =============================================================================

# Start development environment (DB in Docker, app local)
dev-start:
	@echo "🚀 Starting development environment..."
	@make db-start
	@echo "📊 Running migrations..."
	@make migrate-up
	@echo "🌱 Running seed..."
	@make seed
	@echo "✅ Development environment ready!"
	@echo "💡 Run 'make run' to start the application locally"

# Stop development environment
dev-stop:
	@make db-stop

# =============================================================================
# Tools
# =============================================================================

# Install development tools
install-tools:
	@echo "🛠️ Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/jackc/tern/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Tools installed!"

# =============================================================================
# Help
# =============================================================================

help:
	@echo "🚀 PandoraGym API - Available Commands"
	@echo "====================================="
	@echo ""
	@echo "🏗️  Build & Run:"
	@echo "  build            - Build the application binary"
	@echo "  run              - Run the application locally"
	@echo "  test             - Run all tests"
	@echo "  clean            - Clean build artifacts"
	@echo ""
	@echo "🐳 Docker Commands (Full Environment):"
	@echo "  docker-setup     - Complete Docker setup with hot reload (recommended)"
	@echo "  docker-run       - Start containers with hot reload enabled"
	@echo "  docker-stop      - Stop containers"
	@echo "  docker-clean     - Clean up Docker resources (volumes, networks, images)"
	@echo "  docker-restart   - Restart containers"
	@echo "  docker-migrate   - Run migrations in Docker"
	@echo "  docker-seed      - Run seed in Docker"
	@echo "  docker-watch     - Watch logs for hot reload changes"
	@echo "  docker-test      - Test API endpoint"
	@echo "  docker-logs      - Show app logs"
	@echo "  docker-status    - Show container status"
	@echo "  docker-shell     - Connect to app container"
	@echo "  docker-clean     - Clean up Docker resources"
	@echo ""
	@echo "🗄️  Database (Local Development):"
	@echo "  db-start         - Start PostgreSQL container"
	@echo "  db-stop          - Stop PostgreSQL container"
	@echo "  db-shell         - Connect to database shell"
	@echo ""
	@echo "🔄 Migrations & Seeding (Local):"
	@echo "  migrate-up       - Run migrations"
	@echo "  migrate-down     - Rollback migrations"
	@echo "  migrate-create name=<name> - Create new migration"
	@echo "  seed             - Run database seed"
	@echo "  setup            - Run migrations + seed"
	@echo ""
	@echo "🧪 Development Workflows:"
	@echo "  dev-start        - Start dev environment (DB in Docker, app local)"
	@echo "  dev-stop         - Stop dev environment"
	@echo ""
	@echo "🛠️  Tools & Quality:"
	@echo "  deps             - Install dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code"
	@echo "  install-tools    - Install development tools"
	@echo ""
	@echo "📊 Quick Start:"
	@echo "  🐳 Full Docker:   make docker-setup"
	@echo "  💻 Development:   make dev-start && make run"
	@echo ""
	@echo "🌐 After setup, API available at: http://localhost:3333"
