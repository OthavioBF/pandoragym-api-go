#!/bin/bash

echo "🚀 Starting PandoraGym API Local Deployment"
echo "=========================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Stop any existing containers
echo "🛑 Stopping existing containers..."
docker-compose down

# Remove old volumes (optional - uncomment if you want fresh data)
# echo "🗑️  Removing old volumes..."
# docker-compose down -v

# Build and start services
echo "🏗️  Building and starting services..."
docker-compose up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."
docker-compose ps

# Show logs
echo "📋 Showing application logs..."
echo "=========================================="
docker-compose logs app

echo ""
echo "✅ Deployment completed!"
echo "🌐 API is available at: http://localhost:3333"
echo "🗄️  Database is available at: localhost:5432"
echo ""
echo "📚 Available endpoints:"
echo "  - GET  /health                    - Health check"
echo "  - POST /auth/session              - User login"
echo "  - POST /auth/register/student     - Student registration"
echo "  - POST /auth/register/personal    - Personal trainer registration"
echo "  - GET  /api/workouts              - List workouts"
echo "  - GET  /api/exercises             - List exercises"
echo ""
echo "🔧 Useful commands:"
echo "  - make docker-logs    - View application logs"
echo "  - make docker-stop    - Stop all containers"
echo "  - make seed           - Re-run database seeding"
echo ""
echo "👤 Sample login credentials:"
echo "  Personal Trainer: carlos@pandoragym.com / 123456"
echo "  Student: joao@email.com / 123456"
