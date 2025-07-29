# PandoraGym Go API

A complete gym management API built with Go, migrated from the original Node.js/TypeScript version. This API provides comprehensive functionality for managing gym operations, including user management, workout planning, scheduling, and personal trainer services.

## 🚀 Features

### User Management
- **Student Registration & Authentication**: Complete student onboarding with profile management
- **Personal Trainer Registration**: Professional trainer accounts with specialized features
- **JWT Authentication**: Secure token-based authentication system
- **Profile Management**: Avatar upload, profile updates, and user preferences

### Workout Management
- **Workout Creation**: Create custom workout routines with exercises
- **Exercise Library**: Comprehensive exercise database with videos and instructions
- **Workout Templates**: Reusable workout templates for efficiency
- **Exercise Configuration**: Detailed exercise setup with sets, reps, and rest times

### Scheduling System
- **Appointment Booking**: Students can book sessions with personal trainers
- **Schedule Management**: Personal trainers can manage their availability
- **Status Tracking**: Complete scheduling lifecycle management
- **24-Hour Cancellation Rule**: Business rule enforcement for cancellations

### Personal Trainer Features
- **Student Management**: Track and manage assigned students
- **Custom Workouts**: Create personalized workout plans for students
- **Progress Tracking**: Monitor student evolution and performance
- **Communication**: Direct messaging system with students

### Business Logic
- **Role-Based Access**: Different permissions for students, trainers, and admins
- **Subscription Management**: Free and premium tier access control
- **QR Code Integration**: Easy trainer-student pairing system
- **Workout Rating System**: Student feedback and trainer evaluation

## 🏗️ Architecture

This project follows Clean Architecture principles with clear separation of concerns:

```
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and routes
│   ├── core/           # Dependency injection
│   ├── middleware/     # HTTP middleware (auth, CORS, etc.)
│   ├── models/         # Domain models and DTOs
│   ├── services/       # Business logic layer
│   ├── infra/          # Infrastructure (database, external services)
│   └── utils/          # Utility functions
└── pkg/                # Public packages (if any)
```

## 🛠️ Tech Stack

- **Language**: Go 1.24.2
- **Web Framework**: Chi Router
- **Database**: PostgreSQL with pgx driver
- **Authentication**: JWT tokens
- **Validation**: Custom validation with struct tags
- **Containerization**: Docker & Docker Compose
- **Database Migrations**: golang-migrate

## 📋 Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Complete setup in one command**
   ```bash
   make docker-setup
   ```

2. **The API will be available at**: `http://localhost:3333`

### Manual Development Setup

1. **Start development environment**
   ```bash
   make dev-start
   ```

2. **Start the application**
   ```bash
   make run
   ```

## 🔧 Available Commands

Run `make help` to see all available commands:

### 🐳 Docker Commands (Full Environment)
- `make docker-setup` - Complete Docker setup (recommended)
- `make docker-run` - Start containers
- `make docker-stop` - Stop containers
- `make docker-migrate` - Run migrations in Docker
- `make docker-seed` - Run seed in Docker
- `make docker-status` - Show container status

### 💻 Development Commands
- `make dev-start` - Start dev environment (DB in Docker, app local)
- `make run` - Run application locally
- `make build` - Build application binary
- `make test` - Run tests

### 🗄️ Database & Migrations
- `make migrate-up` - Run migrations
- `make seed` - Run database seed
- `make db-shell` - Connect to database shell

## 🔧 Configuration

Create a `.env` file in the root directory:

```env
DATABASE_URL=postgresql://username:password@localhost:5432/pandoragym_db?sslmode=disable
PORT=3333
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-supabase-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key
```

## 📚 API Documentation

### Authentication Endpoints
- `POST /auth/session` - User login (returns access + refresh tokens)
- `POST /auth/register/student` - Student registration
- `POST /auth/register/personal` - Personal trainer registration
- `POST /auth/password/recover` - Password recovery
- `POST /auth/password/reset` - Password reset
- `POST /auth/refresh` - Refresh access token using refresh token
- `POST /auth/revoke` - Revoke refresh token (logout)

### User Endpoints
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update user profile
- `POST /api/users/avatar` - Upload avatar

### Workout Endpoints
- `GET /api/workouts` - List workouts
- `POST /api/workouts` - Create workout
- `GET /api/workouts/{id}` - Get workout details
- `PUT /api/workouts/{id}` - Update workout
- `DELETE /api/workouts/{id}` - Delete workout

### Exercise Endpoints
- `GET /api/exercises` - List exercises
- `POST /api/exercises` - Create exercise
- `GET /api/exercises/{id}` - Get exercise details
- `PUT /api/exercises/{id}` - Update exercise
- `DELETE /api/exercises/{id}` - Delete exercise

### Scheduling Endpoints
- `GET /api/schedulings` - List schedulings
- `POST /api/schedulings` - Create scheduling
- `GET /api/schedulings/{id}` - Get scheduling details
- `PUT /api/schedulings/{id}` - Update scheduling
- `DELETE /api/schedulings/{id}` - Cancel scheduling

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run benchmarks
make bench
```

## 🔨 Development

### Available Make Commands

```bash
make build        # Build the application
make run          # Run the application
make test         # Run tests
make dev          # Run with hot reload (requires air)
make docker-run   # Run with Docker Compose
make docker-stop  # Stop Docker containers
make migrate-up   # Run database migrations
make migrate-down # Rollback migrations
make fmt          # Format code
make lint         # Lint code
```

### Database Migrations

Create a new migration:
```bash
make migrate-create name=create_users_table
```

Run migrations:
```bash
make migrate-up
```

Rollback migrations:
```bash
make migrate-down
```

## 🏢 Business Rules

1. **User Access**: Users can access free exercises and create personal workouts
2. **Workout Ownership**: Users can only modify workouts they created
3. **Premium Features**: Subscription required for personal trainer access
4. **Trainer Assignment**: Users choose trainers via QR code or platform method
5. **Trainer Access**: Personal trainers can only access their contracted students
6. **Scheduling Visibility**: Trainers see only their scheduled appointments
7. **Workout Types**: Trainers can create both exclusive and free workouts
8. **Edit Permissions**: Users can only edit their own created content
9. **Category Management**: Trainers can create new categories if they don't exist
10. **Communication**: Trainers can message students respecting privacy settings
11. **Workout Assignment**: Trainers can assign workouts to specific students
12. **Cancellation Policy**: 24-hour advance notice required for cancellations
13. **Pre-loaded Content**: Application includes 10 pre-registered workout programs
14. **Rating System**: Students can only rate workouts from contracted trainers
15. **Public Profiles**: Trainer profiles are publicly visible for selection
16. **Single Trainer Rule**: Students can only have one personal trainer

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or need help, please:

1. Check the [Issues](../../issues) page
2. Create a new issue if your problem isn't already reported
3. Provide detailed information about your environment and the issue

## 🗺️ Roadmap

- [ ] Complete database migration scripts
- [ ] Implement file upload to Supabase/S3
- [ ] Add comprehensive test coverage
- [ ] Implement real-time notifications
- [ ] Add API rate limiting
- [ ] Implement caching layer
- [ ] Add monitoring and logging
- [ ] Create admin dashboard endpoints
- [ ] Implement WhatsApp integration
- [ ] Add workout analytics and reporting
- [x] **Refresh Token API** - Secure token refresh system with rotation

---

**Note**: This is a migration from the original Node.js/TypeScript version. Some features may still be in development. Please check the roadmap for current status.
