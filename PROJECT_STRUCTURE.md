# PandoraGym Go API - Project Structure

This document outlines the complete project structure of the migrated PandoraGym API from Node.js/TypeScript to Go.

## 📁 Directory Structure

```
pandoragym-go-api/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/                        # HTTP layer
│   │   ├── api.go                  # API struct definition
│   │   ├── auth_handlers.go        # Authentication endpoints
│   │   ├── health_handlers.go      # Health check and utility endpoints
│   │   ├── health_test.go          # Health endpoint tests
│   │   ├── middleware.go           # HTTP middleware (CORS, etc.)
│   │   ├── routes.go               # Route definitions and binding
│   │   ├── scheduling_handlers.go  # Scheduling management endpoints
│   │   ├── user_handlers.go        # User management endpoints
│   │   └── workout_handlers.go     # Workout and exercise endpoints
│   ├── core/
│   │   └── di.go                   # Dependency injection container
│   ├── infra/                      # Infrastructure layer
│   │   ├── migrations/
│   │   │   └── 001_initial_schema.sql # Database schema migration
│   │   └── pgstore/
│   │       └── pgstore.go          # PostgreSQL connection setup
│   ├── middleware/
│   │   └── auth/
│   │       └── jwt.go              # JWT authentication middleware
│   ├── models/                     # Domain models and DTOs
│   │   ├── scheduling.go           # Scheduling-related models
│   │   ├── user.go                 # User, Student, Personal models
│   │   └── workout.go              # Workout and Exercise models
│   ├── services/                   # Business logic layer
│   │   ├── auth_service.go         # Authentication business logic
│   │   ├── scheduling_service.go   # Scheduling business logic
│   │   ├── user_service.go         # User management business logic
│   │   └── workout_service.go      # Workout management business logic
│   └── utils/                      # Utility functions
│       ├── response.go             # HTTP response helpers
│       └── validation.go           # Input validation utilities
├── .env                            # Environment variables
├── .gitignore                      # Git ignore rules
├── docker-compose.yml              # Docker Compose configuration
├── Dockerfile                      # Docker image definition
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── Makefile                        # Build and development commands
├── PROJECT_STRUCTURE.md            # This file
└── README.md                       # Project documentation
```

## 🏗️ Architecture Overview

### Clean Architecture Layers

1. **Presentation Layer** (`internal/api/`)
   - HTTP handlers for all endpoints
   - Request/response processing
   - Route definitions and middleware

2. **Business Logic Layer** (`internal/services/`)
   - Core business rules implementation
   - Data validation and processing
   - Service interfaces and implementations

3. **Data Layer** (`internal/infra/`)
   - Database connections and queries
   - External service integrations
   - Infrastructure concerns

4. **Domain Layer** (`internal/models/`)
   - Business entities and value objects
   - Data transfer objects (DTOs)
   - Domain-specific types and enums

### Key Components

#### API Handlers
- **Auth Handlers**: Login, registration, password recovery
- **User Handlers**: Profile management, avatar upload
- **Workout Handlers**: Workout and exercise CRUD operations
- **Scheduling Handlers**: Appointment booking and management
- **Health Handlers**: System health checks and utilities

#### Services
- **Auth Service**: JWT token management, authentication logic
- **User Service**: User account management, profile operations
- **Workout Service**: Workout and exercise business logic
- **Scheduling Service**: Appointment scheduling and status management

#### Models
- **User Models**: User, Student, Personal trainer entities
- **Workout Models**: Workout, Exercise, ExerciseConfig entities
- **Scheduling Models**: Scheduling, PersonalSchedule, Message entities

#### Infrastructure
- **Database**: PostgreSQL with pgx driver
- **Migrations**: SQL schema migrations
- **Authentication**: JWT-based authentication
- **Validation**: Custom struct tag validation

## 🔄 Migration from Node.js

### Original Structure (Node.js/TypeScript)
```
src/
├── @types/
├── http/
│   ├── routes/
│   │   ├── auth/
│   │   ├── mobile/
│   │   ├── web/
│   │   └── admin/
│   ├── middlewares/
│   └── server.ts
├── lib/
├── schema/
└── utils/
```

### New Structure (Go)
```
internal/
├── api/           # Equivalent to http/routes/
├── middleware/    # Equivalent to http/middlewares/
├── models/        # Equivalent to @types/ and schema/
├── services/      # New business logic layer
├── infra/         # Equivalent to lib/
└── utils/         # Enhanced utilities
```

## 🚀 Key Features Implemented

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (Student, Personal, Admin)
- Password hashing with bcrypt
- Protected route middleware

### User Management
- Student and Personal trainer registration
- Profile management and avatar upload
- User role differentiation

### Workout System
- Workout creation and management
- Exercise library with templates
- Exercise configuration for workouts
- Workout assignment to students

### Scheduling System
- Appointment booking between students and trainers
- Schedule status management
- 24-hour cancellation policy enforcement
- Scheduling history tracking

### Business Rules Implementation
- Workout ownership validation
- Trainer-student relationship management
- Premium feature access control
- Category management for exercises

## 🛠️ Development Workflow

### Building and Running
```bash
# Build the application
make build

# Run in development mode
make run

# Run with Docker
make docker-run

# Run tests
make test
```

### Database Operations
```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create name=migration_name
```

### Code Quality
```bash
# Format code
make fmt

# Lint code
make lint

# Generate coverage report
make coverage
```

## 📊 Database Schema

The database schema includes the following main entities:

- **users**: Base user information
- **personal**: Personal trainer specific data
- **student**: Student specific data
- **workouts**: Workout definitions
- **exercises_templates**: Exercise library
- **exercises_setup**: Exercise configurations in workouts
- **schedulings**: Appointment bookings
- **workouts_history**: Workout execution history
- **workouts_rating**: Workout ratings and feedback

## 🔐 Security Features

- JWT token authentication
- Password hashing with bcrypt
- Role-based access control
- Input validation and sanitization
- CORS middleware
- Request rate limiting (planned)

## 📈 Performance Considerations

- Connection pooling for database
- Efficient database queries with indexes
- Structured logging for monitoring
- Graceful error handling
- Resource cleanup and connection management

## 🧪 Testing Strategy

- Unit tests for business logic
- Integration tests for API endpoints
- Database transaction testing
- Mock services for external dependencies
- Coverage reporting and benchmarking

This structure provides a solid foundation for the PandoraGym API with clear separation of concerns, maintainable code organization, and scalable architecture patterns.
