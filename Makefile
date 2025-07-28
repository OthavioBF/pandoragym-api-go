.PHONY: build run test clean docker-build docker-run docker-stop migrate-up migrate-down seed setup deploy-db deploy-app local-start local-stop local-restart local-status local-setup

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
	@echo "📊 Starting PostgreSQL Docker container..."
	@docker-compose -f docker-compose.db.yml up -d
	@echo "⏳ Waiting for database to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if docker-compose -f docker-compose.db.yml exec -T db pg_isready -U pandoragym -d pandoragym_db > /dev/null 2>&1; then \
			echo "✅ Database is ready!"; \
			exit 0; \
		fi; \
		echo "Waiting... ($$i/10)"; \
		sleep 3; \
	done; \
	echo "❌ Database failed to start!" && exit 1

db-stop:
	@echo "🛑 Stopping PostgreSQL Docker container..."
	@docker-compose -f docker-compose.db.yml stop
	@echo "✅ Database stopped!"

db-restart:
	@echo "🔄 Restarting PostgreSQL Docker container..."
	@make db-stop
	@make db-start

db-logs:
	@echo "📋 Showing PostgreSQL logs..."
	@docker-compose -f docker-compose.db.yml logs -f db

db-shell:
	@echo "🐚 Connecting to PostgreSQL shell..."
	@docker-compose -f docker-compose.db.yml exec db psql -U pandoragym -d pandoragym_db

db-status:
	@echo "📊 Database Status:"
	@echo "=================="
	@printf "Container: "
	@if docker-compose -f docker-compose.db.yml ps | grep pandoragym_db | grep Up > /dev/null; then \
		echo "✅ Running"; \
	else \
		echo "❌ Stopped"; \
	fi
	@printf "Health: "
	@if docker-compose -f docker-compose.db.yml exec -T db pg_isready -U pandoragym -d pandoragym_db > /dev/null 2>&1; then \
		echo "✅ Healthy"; \
	else \
		echo "❌ Unhealthy"; \
	fi

db-remove:
	@echo "🗑️  Removing PostgreSQL Docker container and volumes..."
	@docker-compose -f docker-compose.db.yml down -v
	@echo "✅ Database removed!"

db-reset: db-remove db-start
	@echo "🔄 Database reset complete!"

# Local services management (PostgreSQL + Go App)
local-start:
	@echo "🚀 Starting local PandoraGym services..."
	@echo "📊 Starting PostgreSQL with Docker..."
	@docker-compose -f docker-compose.db.yml up -d
	@echo "⏳ Waiting for database to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if docker-compose -f docker-compose.db.yml exec -T db pg_isready -U pandoragym -d pandoragym_db > /dev/null 2>&1; then \
			echo "✅ Database is ready!"; \
			break; \
		fi; \
		echo "Waiting... ($$i/10)"; \
		sleep 3; \
		if [ $$i -eq 10 ]; then \
			echo "❌ Database failed to start!" && exit 1; \
		fi; \
	done
	@echo "🚀 Starting Go application with hot reload..."
	@nohup air > app.log 2>&1 & echo $$! > app.pid
	@sleep 5
	@echo "🔍 Testing application..."
	@if curl -s http://localhost:3333/health > /dev/null; then \
		echo "✅ Application is running successfully!"; \
		echo "🌐 API available at: http://localhost:3333"; \
		echo "🗄️  Database available at: localhost:5432"; \
		echo ""; \
		echo "💡 Tip: If you need database setup, run 'make local-setup'"; \
	else \
		echo "❌ Application failed to start!"; \
		echo "💡 Try running 'make local-setup' for complete setup"; \
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
	@echo "🗄️  Stopping PostgreSQL Docker container..."
	@docker-compose -f docker-compose.db.yml down
	@echo "✅ All services stopped!"

local-restart: local-stop local-start

local-status:
	@echo "📊 Local Services Status:"
	@echo "========================="
	@printf "PostgreSQL: "
	@if docker-compose -f docker-compose.db.yml ps | grep pandoragym_db | grep Up > /dev/null; then \
		echo "✅ Running (Docker)"; \
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
	@if docker-compose -f docker-compose.db.yml exec -T db pg_isready -U pandoragym -d pandoragym_db > /dev/null 2>&1; then \
		echo "✅ Connected (port 5432)"; \
	else \
		echo "❌ Failed"; \
	fi

local-setup:
	@echo "🎯 PandoraGym API - Complete Local Setup"
	@echo "========================================"
	@echo ""
	@echo "📊 Starting PostgreSQL with Docker..."
	@docker-compose -f docker-compose.db.yml up -d
	@echo "⏳ Waiting for database to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if docker-compose -f docker-compose.db.yml exec -T db pg_isready -U pandoragym -d pandoragym_db > /dev/null 2>&1; then \
			echo "✅ Database is ready!"; \
			break; \
		fi; \
		echo "Waiting... ($$i/10)"; \
		sleep 3; \
		if [ $$i -eq 10 ]; then \
			echo "❌ Database failed to start!" && exit 1; \
		fi; \
	done
	@echo "🔗 Testing database connection..."
	@if docker-compose -f docker-compose.db.yml exec -T db psql -U pandoragym -d pandoragym_db -c "SELECT 1;" > /dev/null 2>&1; then \
		echo "✅ Database connection successful!"; \
	else \
		echo "❌ Database connection failed!"; \
		exit 1; \
	fi
	@echo "🚀 Starting Go application with hot reload..."
	@nohup air > app.log 2>&1 & echo $$! > app.pid
	@sleep 5
	@echo "🔍 Testing application..."
	@if curl -s http://localhost:3333/health > /dev/null; then \
		echo "✅ Application is running successfully!"; \
	else \
		echo "❌ Application failed to start!"; \
		exit 1; \
	fi
	@echo ""
	@echo "🔄 Running database migrations..."
	@make migrate-up
	@echo ""
	@echo "🌱 Seeding database with sample data..."
	@make seed
	@echo ""
	@echo "🎉 Local setup completed successfully!"
	@echo ""
	@echo "📋 Your PandoraGym API is ready:"
	@echo "  🌐 API: http://localhost:3333"
	@echo "  🗄️  Database: localhost:5433"
	@echo "  📊 Status: make local-status"
	@echo ""
	@echo "🧪 Test credentials:"
	@echo "  👨‍💼 Personal Trainer: carlos@pandoragym.com / 123456"
	@echo "  🎓 Student: joao@email.com / 123456"

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
	@echo "  local-setup      - Complete setup (create DB + start + migrate + seed)"
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
