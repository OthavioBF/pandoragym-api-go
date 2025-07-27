# PandoraGym API - Local Deployment Guide

This guide will help you deploy the PandoraGym API locally using Docker with automatic database migrations and seeding.

## Prerequisites

- Docker and Docker Compose installed
- Git (to clone the repository)
- Make (optional, for using Makefile commands)

## Quick Start

### Option 1: Using the Deployment Script (Recommended)

```bash
./deploy.sh
```

This script will:
1. Stop any existing containers
2. Build and start all services (database, migrations, seed, API)
3. Show you the status and logs
4. Provide useful information about endpoints and credentials

### Option 2: Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app
```

### Option 3: Using Makefile

```bash
# Start services
make docker-run

# View logs
make docker-logs

# Stop services
make docker-stop
```

## Services Architecture

The deployment includes 4 services that run in sequence:

1. **Database (db)**: PostgreSQL 15 with health checks
2. **Migrations (migrate)**: Runs database migrations using Tern
3. **Seed (seed)**: Populates database with sample data
4. **API (app)**: The main application server

## Database Schema

The application creates the following main tables:
- `users` - User accounts (students and personal trainers)
- `personal` - Personal trainer profiles
- `student` - Student profiles
- `workout` - Workout routines
- `exercises_template` - Exercise library
- `exercises_setup` - Exercises within workouts
- `scheduling` - Training sessions
- `message` - Communication between trainers and students

## Sample Data

The seed process creates:

### Personal Trainers
- **Carlos Silva** (carlos@pandoragym.com) - Hipertrofia specialist
- **Ana Costa** (ana@pandoragym.com) - Functional training specialist  
- **Roberto Santos** (roberto@pandoragym.com) - Weight loss specialist

### Students
- **João Silva** (joao@email.com) - Muscle gain goal
- **Maria Santos** (maria@email.com) - Weight loss goal
- **Pedro Costa** (pedro@email.com) - Conditioning goal

### Sample Exercises
- Supino Reto, Agachamento Livre, Puxada Frontal, etc.

### Sample Workouts
- Treino de Peito e Tríceps
- Treino de Pernas  
- Treino Funcional

**Default Password**: `123456` for all users

## API Endpoints

### Authentication
- `POST /auth/session` - User login
- `POST /auth/register/student` - Student registration
- `POST /auth/register/personal` - Personal trainer registration
- `POST /auth/password/recover` - Password recovery
- `POST /auth/password/reset` - Password reset

### User Management
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update user profile
- `POST /api/users/avatar` - Upload avatar

### Workouts
- `GET /api/workouts` - List workouts
- `POST /api/workouts` - Create workout
- `GET /api/workouts/{id}` - Get workout details
- `PUT /api/workouts/{id}` - Update workout
- `DELETE /api/workouts/{id}` - Delete workout

### Exercises
- `GET /api/exercises` - List exercises
- `POST /api/exercises` - Create exercise
- `GET /api/exercises/{id}` - Get exercise details
- `PUT /api/exercises/{id}` - Update exercise
- `DELETE /api/exercises/{id}` - Delete exercise

### Scheduling
- `GET /api/schedulings` - List schedulings
- `POST /api/schedulings` - Create scheduling
- `GET /api/schedulings/{id}` - Get scheduling details
- `PUT /api/schedulings/{id}` - Update scheduling
- `DELETE /api/schedulings/{id}` - Cancel scheduling

## Environment Configuration

The application uses these environment variables:

```env
# Database
DATABASE_URL=postgresql://pandoragym:password@localhost:5432/pandoragym_db?sslmode=disable
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=pandoragym_db
DATABASE_USER=pandoragym
DATABASE_PASSWORD=password

# Server
PORT=3333
BASE_URL=http://localhost:3333

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Supabase (optional)
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-key
```

## Development Commands

```bash
# Database operations
make migrate-up          # Run migrations
make migrate-down        # Rollback migrations
make seed               # Seed database
make setup              # Run migrations + seed

# Application
make run                # Run locally (without Docker)
make build              # Build binary
make test               # Run tests
make dev                # Run with hot reload

# Docker operations
make docker-run         # Start with Docker
make docker-stop        # Stop containers
make docker-logs        # View logs
```

## Troubleshooting

### Database Connection Issues
```bash
# Check if database is running
docker-compose ps

# View database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Migration Issues
```bash
# Run migrations manually
docker-compose run --rm migrate

# Reset database (WARNING: destroys data)
docker-compose down -v
docker-compose up -d
```

### Application Issues
```bash
# View application logs
docker-compose logs app

# Restart application
docker-compose restart app

# Rebuild application
docker-compose up -d --build app
```

### Port Conflicts
If port 3333 or 5432 is already in use:

1. Stop the conflicting service
2. Or modify the ports in `docker-compose.yml`

## Testing the API

### Login Example
```bash
curl -X POST http://localhost:3333/auth/session \
  -H "Content-Type: application/json" \
  -d '{
    "email": "carlos@pandoragym.com",
    "password": "123456"
  }'
```

### Get Workouts Example
```bash
# First login to get token, then:
curl -X GET http://localhost:3333/api/workouts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Production Deployment

For production deployment:

1. Update environment variables in `.env`
2. Use production database credentials
3. Change JWT_SECRET to a secure value
4. Configure proper CORS settings
5. Set up SSL/TLS certificates
6. Use a production-ready Docker setup

## Support

If you encounter issues:

1. Check the logs: `make docker-logs`
2. Verify all services are running: `docker-compose ps`
3. Ensure ports 3333 and 5432 are available
4. Check the database connection in `.env`

The API will be available at `http://localhost:3333` once deployment is complete!
