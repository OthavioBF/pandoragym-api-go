# PandoraGym API - Quick Start

## 🚀 Complete Local Deployment Setup

I've set up your PandoraGym API with:
- ✅ Tern migrations (like in your imperium project)
- ✅ Database seeding (like in your pandoragym-api project)
- ✅ Docker deployment with proper service orchestration
- ✅ Sample data with personal trainers, students, exercises, and workouts

## 📁 What was created/updated:

### Migration System (Tern)
- `cmd/tern/main.go` - Migration runner
- `internal/infra/pgstore/migrations/tern.conf` - Tern configuration
- `internal/infra/pgstore/migrations/001_create_tables.sql` - Complete database schema

### Seed System
- `cmd/seed/main.go` - Database seeding with sample data

### Docker Setup
- `docker-compose.yml` - Multi-service setup (db → migrate → seed → app)
- `Dockerfile` - Updated for development
- `deploy.sh` - One-command deployment script

### Configuration
- `.env` - Updated for local development
- `Makefile` - Added migration and seed commands
- `DEPLOYMENT.md` - Complete deployment guide

## 🏃‍♂️ How to run:

### 1. Start Docker Desktop
Make sure Docker is running on your machine.

### 2. Deploy with one command:
```bash
./deploy.sh
```

### 3. Or use individual commands:
```bash
# Start all services
make docker-run

# View logs
make docker-logs

# Stop services
make docker-stop
```

## 🎯 What happens during deployment:

1. **Database starts** (PostgreSQL 15)
2. **Migrations run** (creates all tables with proper schema)
3. **Seeding runs** (populates with sample data)
4. **API starts** (your Go application)

## 📊 Sample Data Created:

### Personal Trainers:
- carlos@pandoragym.com / 123456 (Hipertrofia specialist)
- ana@pandoragym.com / 123456 (Functional training)
- roberto@pandoragym.com / 123456 (Weight loss)

### Students:
- joao@email.com / 123456 (Muscle gain goal)
- maria@email.com / 123456 (Weight loss goal)
- pedro@email.com / 123456 (Conditioning goal)

### Sample Exercises:
- Supino Reto, Agachamento Livre, Puxada Frontal, etc.

### Sample Workouts:
- Treino de Peito e Tríceps
- Treino de Pernas
- Treino Funcional

## 🌐 API Access:
- **URL**: http://localhost:3333
- **Database**: localhost:5432
- **Health Check**: GET http://localhost:3333/health

## 🧪 Test the API:

### Login:
```bash
curl -X POST http://localhost:3333/auth/session \
  -H "Content-Type: application/json" \
  -d '{
    "email": "carlos@pandoragym.com",
    "password": "123456"
  }'
```

### Get Workouts:
```bash
curl -X GET http://localhost:3333/api/workouts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🔧 Development Commands:

```bash
# Database
make migrate-up    # Run migrations
make seed         # Seed database
make setup        # Migrate + seed

# Application
make run          # Run without Docker
make build        # Build binary
make test         # Run tests

# Docker
make docker-run   # Start with Docker
make docker-stop  # Stop containers
make docker-logs  # View logs
```

## 🎉 Ready to go!

Your PandoraGym API is now ready for local development with:
- Complete database schema
- Sample data for testing
- Proper migration system using Tern
- Docker orchestration
- All your original business logic intact

Just start Docker and run `./deploy.sh`!
