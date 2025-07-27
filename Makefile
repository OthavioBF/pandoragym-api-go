.PHONY: build run test clean docker-build docker-run docker-stop migrate-up migrate-down seed setup deploy-db deploy-app local-start local-stop local-restart local-status

# Build the application
build:
	go build -o bin/pandoragym-api cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Docker commands - Full deployment (database + app)
docker-build:
	docker build -t pandoragym-api .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f app

# Separate deployment commands
deploy-db:
	./deploy-db.sh

deploy-app:
	./deploy-app.sh

# Database-only Docker commands
db-start:
	docker-compose -f docker-compose.db.yml up -d

db-stop:
	docker-compose -f docker-compose.db.yml stop

db-logs:
	docker-compose -f docker-compose.db.yml logs -f db

db-remove:
	docker-compose -f docker-compose.db.yml down -v

# Local services management (PostgreSQL + Go App)
local-start:
	@echo "🚀 Starting local PandoraGym services..."
	@echo "📊 Starting PostgreSQL..."
	@brew services start postgresql@15
	@sleep 3
	@echo "🔗 Testing database connection..."
	@export PATH="/opt/homebrew/opt/postgresql@15/bin:$$PATH" && \
	if psql -h localhost -U pandoragym -d pandoragym_db -c "SELECT 1;" > /dev/null 2>&1; then \
		echo "✅ Database connection successful!"; \
	else \
		echo "❌ Database connection failed!"; \
		exit 1; \
	fi
	@echo "🚀 Starting Go application..."
	@nohup make run > app.log 2>&1 & echo $$! > app.pid
	@sleep 5
	@echo "🔍 Testing application..."
	@if curl -s http://localhost:3333/health > /dev/null; then \
		echo "✅ Application is running successfully!"; \
		echo "🌐 API available at: http://localhost:3333"; \
		echo "🗄️  Database available at: localhost:5432"; \
	else \
		echo "❌ Application failed to start!"; \
		exit 1; \
	fi

local-stop:
	@echo "🛑 Stopping local PandoraGym services..."
	@echo "📊 Stopping Go application..."
	@if [ -f app.pid ]; then \
		kill `cat app.pid` 2>/dev/null || true; \
		rm -f app.pid; \
	fi
	@pkill -f "go run cmd/server/main.go" 2>/dev/null || true
	@pkill -f "pandoragym-api" 2>/dev/null || true
	@lsof -ti:3333 | xargs kill -9 2>/dev/null || true
	@echo "🗄️  Stopping PostgreSQL..."
	@brew services stop postgresql@15
	@echo "✅ All services stopped!"

local-restart: local-stop local-start

local-status:
	@echo "📊 Local Services Status:"
	@echo "========================="
	@printf "PostgreSQL: "
	@if brew services list | grep postgresql@15 | grep started > /dev/null; then \
		echo "✅ Running"; \
	else \
		echo "❌ Stopped"; \
	fi
	@printf "Go Application: "
	@if curl -s http://localhost:3333/health > /dev/null 2>&1; then \
		echo "✅ Running (http://localhost:3333)"; \
	else \
		echo "❌ Stopped"; \
	fi
	@printf "Database Connection: "
	@export PATH="/opt/homebrew/opt/postgresql@15/bin:$$PATH" && \
	if psql -h localhost -U pandoragym -d pandoragym_db -c "SELECT 1;" > /dev/null 2>&1; then \
		echo "✅ Connected"; \
	else \
		echo "❌ Failed"; \
	fi

# Database migrations using tern
migrate-up:
	go run ./cmd/tern

migrate-down:
	tern migrate --migrations ./internal/infra/pgstore/migrations --config ./internal/infra/pgstore/migrations/tern.conf --destination 0

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	tern new --migrations ./internal/infra/pgstore/migrations $(name)

# Seed database with sample data
seed:
	go run ./cmd/seed

# Complete setup (migrations + seed)
setup: migrate-up seed

# Development commands
dev:
	air

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/jackc/tern/v2@latest

# Generate API documentation
docs:
	swag init -g cmd/server/main.go

# Security scan
security:
	gosec ./...

# Benchmark tests
bench:
	go test -bench=. -benchmem ./...

# Coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "🚀 Local Services:"
	@echo "  local-start      - Start PostgreSQL and Go application"
	@echo "  local-stop       - Stop PostgreSQL and Go application"
	@echo "  local-restart    - Restart all local services"
	@echo "  local-status     - Check status of all local services"
	@echo ""
	@echo "🚀 Deployment:"
	@echo "  deploy-db        - Deploy database only (Docker)"
	@echo "  deploy-app       - Deploy Go application (requires database)"
	@echo "  docker-run       - Full deployment (database + app in Docker)"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  db-start         - Start database container only"
	@echo "  db-stop          - Stop database container"
	@echo "  db-logs          - View database logs"
	@echo "  db-remove        - Remove database container and data"
	@echo ""
	@echo "🔄 Migrations & Seeding:"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback all migrations"
	@echo "  migrate-create name=<name> - Create new migration"
	@echo "  seed             - Seed database with sample data"
	@echo "  setup            - Run migrations and seed"
	@echo ""
	@echo "🏗️  Development:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  test             - Run tests"
	@echo "  dev              - Run with hot reload"
	@echo "  clean            - Clean build artifacts"
	@echo "  deps             - Install dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code"
	@echo ""
	@echo "📊 Analysis:"
	@echo "  coverage         - Generate test coverage report"
	@echo "  bench            - Run benchmarks"
	@echo "  security         - Run security scan"
