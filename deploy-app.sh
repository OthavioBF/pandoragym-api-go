#!/bin/bash

echo "🚀 Starting PandoraGym Go Application Deployment"
echo "=============================================="

# Check if database is running
echo "🔍 Checking if database is running..."
if ! docker ps | grep -q "pandoragym_db"; then
    echo "❌ Database is not running. Please run ./deploy-db.sh first."
    exit 1
fi

# Test database connection
echo "🔌 Testing database connection..."
if ! docker exec pandoragym_db pg_isready -U pandoragym -d pandoragym_db > /dev/null 2>&1; then
    echo "❌ Cannot connect to database. Please check if database is healthy."
    exit 1
fi

echo "✅ Database connection successful!"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24.2 or later."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy

# Run migrations
echo "🗄️  Running database migrations..."
if go run ./cmd/tern; then
    echo "✅ Migrations completed successfully!"
else
    echo "❌ Migration failed. Please check the error above."
    exit 1
fi

# Seed database
echo "🌱 Seeding database with sample data..."
if go run ./cmd/seed; then
    echo "✅ Database seeding completed successfully!"
else
    echo "❌ Database seeding failed. Please check the error above."
    exit 1
fi

# Build application
echo "🏗️  Building Go application..."
if make build; then
    echo "✅ Application built successfully!"
else
    echo "❌ Build failed. Please check the error above."
    exit 1
fi

# Start application
echo "🚀 Starting Go application..."
echo "=========================================="
echo ""
echo "🎉 Application is starting!"
echo "🌐 API will be available at: http://localhost:3333"
echo "🗄️  Database is running at: localhost:5432"
echo ""
echo "📚 Sample login credentials:"
echo "  Personal Trainer: carlos@pandoragym.com / 123456"
echo "  Student: joao@email.com / 123456"
echo ""
echo "🔧 To stop the application: Press Ctrl+C"
echo "🔧 To stop the database: docker-compose -f docker-compose.db.yml stop"
echo ""
echo "=========================================="

# Run the application
make run
