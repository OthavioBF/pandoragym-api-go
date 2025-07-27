# PandoraGym API - Separate Deployment Guide

This guide shows how to deploy the database and Go application separately, giving you more control over each component.

## 🗄️ Step 1: Deploy Database

### Option A: Using the deployment script (Recommended)
```bash
./deploy-db.sh
```

### Option B: Using Makefile
```bash
make deploy-db
# or
make db-start
```

### Option C: Using Docker Compose directly
```bash
docker-compose -f docker-compose.db.yml up -d
```

### What happens:
- ✅ PostgreSQL 15 container starts
- ✅ Database `pandoragym_db` is created
- ✅ User `pandoragym` with password `password` is created
- ✅ Health checks ensure database is ready
- ✅ Data persists in Docker volume `postgres_data`

### Database Details:
- **Host**: localhost
- **Port**: 5432
- **Database**: pandoragym_db
- **User**: pandoragym
- **Password**: password
- **Connection String**: `postgresql://pandoragym:password@localhost:5432/pandoragym_db?sslmode=disable`

## 🚀 Step 2: Deploy Go Application

### Option A: Using the deployment script (Recommended)
```bash
./deploy-app.sh
```

### Option B: Using individual commands
```bash
# Install dependencies
make deps

# Run migrations
make migrate-up

# Seed database
make seed

# Build and run application
make build
make run
```

### What happens:
- ✅ Checks database connectivity
- ✅ Installs Go dependencies
- ✅ Runs database migrations (creates all tables)
- ✅ Seeds database with sample data
- ✅ Builds the Go application
- ✅ Starts the API server on port 3333

## 🎯 Complete Deployment Process

### 1. Deploy Database
```bash
./deploy-db.sh
```
**Output:**
```
🗄️  Starting PandoraGym Database Deployment
==========================================
🛑 Stopping existing database container...
🚀 Starting PostgreSQL database...
⏳ Waiting for database to be ready...
🔍 Checking database health...
✅ Database is healthy and ready!
🔌 Testing database connection...
✅ Database connection successful!

🎉 Database deployment completed!
```

### 2. Deploy Go Application
```bash
./deploy-app.sh
```
**Output:**
```
🚀 Starting PandoraGym Go Application Deployment
==============================================
🔍 Checking if database is running...
✅ Database connection successful!
✅ Go is installed: go version go1.24.2
📦 Installing Go dependencies...
🗄️  Running database migrations...
✅ Migrations completed successfully!
🌱 Seeding database with sample data...
✅ Database seeding completed successfully!
🏗️  Building Go application...
✅ Application built successfully!
🚀 Starting Go application...
```

## 🔧 Management Commands

### Database Management
```bash
# Start database
make db-start

# Stop database
make db-stop

# View database logs
make db-logs

# Remove database and all data (⚠️ DESTRUCTIVE)
make db-remove
```

### Application Management
```bash
# Run migrations only
make migrate-up

# Seed database only
make seed

# Run migrations + seed
make setup

# Build application
make build

# Run application
make run

# Run with hot reload
make dev
```

### Development Workflow
```bash
# 1. Start database (once)
make db-start

# 2. Setup database (once or when schema changes)
make setup

# 3. Run application
make run
# or for development with hot reload
make dev
```

## 🧪 Testing the Deployment

### 1. Check Database
```bash
# Test connection
docker exec pandoragym_db pg_isready -U pandoragym -d pandoragym_db

# Connect to database
docker exec -it pandoragym_db psql -U pandoragym -d pandoragym_db

# List tables
\dt

# Check sample data
SELECT name, email, role FROM users;
```

### 2. Test API
```bash
# Health check
curl http://localhost:3333/health

# Login with sample user
curl -X POST http://localhost:3333/auth/session \
  -H "Content-Type: application/json" \
  -d '{
    "email": "carlos@pandoragym.com",
    "password": "123456"
  }'

# Get workouts (use token from login response)
curl -X GET http://localhost:3333/api/workouts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📊 Sample Data Available

After seeding, you'll have:

### Personal Trainers:
- **Carlos Silva** (carlos@pandoragym.com / 123456) - Hipertrofia specialist
- **Ana Costa** (ana@pandoragym.com / 123456) - Functional training specialist
- **Roberto Santos** (roberto@pandoragym.com / 123456) - Weight loss specialist

### Students:
- **João Silva** (joao@email.com / 123456) - Muscle gain goal
- **Maria Santos** (maria@email.com / 123456) - Weight loss goal
- **Pedro Costa** (pedro@email.com / 123456) - Conditioning goal

### Sample Exercises:
- Supino Reto, Agachamento Livre, Puxada Frontal, Desenvolvimento com Halteres, Rosca Direta

### Sample Workouts:
- Treino de Peito e Tríceps
- Treino de Pernas
- Treino Funcional

## 🔄 Common Workflows

### Fresh Start (Clean Database)
```bash
# Remove all data and start fresh
make db-remove
make db-start
make setup
make run
```

### Schema Changes
```bash
# Create new migration
make migrate-create name=add_new_table

# Edit the migration file in internal/infra/pgstore/migrations/
# Then run migrations
make migrate-up
```

### Reset Sample Data
```bash
# Re-run seeding (will add duplicate data)
make seed

# For clean seed, reset database first
make db-remove
make db-start
make setup
```

### Development Mode
```bash
# Start database once
make db-start

# Run app with hot reload
make dev
```

## 🚨 Troubleshooting

### Database Issues
```bash
# Check if database container is running
docker ps | grep pandoragym_db

# View database logs
make db-logs

# Test database connection
docker exec pandoragym_db pg_isready -U pandoragym -d pandoragym_db

# Connect to database manually
docker exec -it pandoragym_db psql -U pandoragym -d pandoragym_db
```

### Application Issues
```bash
# Check Go installation
go version

# Check dependencies
go mod tidy

# Test database connection from Go
go run ./cmd/tern

# Build application manually
make build

# Run with verbose logging
go run cmd/server/main.go
```

### Port Conflicts
If ports 3333 or 5432 are in use:

1. **For database (port 5432)**: Edit `docker-compose.db.yml` and change the port mapping
2. **For application (port 3333)**: Edit `.env` file and change `PORT=3333`

### Migration Issues
```bash
# Check migration status
tern status --migrations ./internal/infra/pgstore/migrations --config ./internal/infra/pgstore/migrations/tern.conf

# Rollback all migrations
make migrate-down

# Run migrations again
make migrate-up
```

## 🎉 Benefits of Separate Deployment

1. **Better Control**: Start/stop database and application independently
2. **Development Friendly**: Keep database running while restarting application
3. **Debugging**: Easier to isolate issues between database and application
4. **Resource Management**: Database runs in Docker, application runs natively
5. **Performance**: Native Go application performance without Docker overhead

## 🔗 Next Steps

Once deployed:
1. **API Documentation**: Available at the endpoints listed in the main README
2. **Frontend Integration**: Connect your frontend to `http://localhost:3333`
3. **Production**: Use the production database URL in `.env` for production deployment

Your PandoraGym API is now running with separate database and application deployments! 🚀
