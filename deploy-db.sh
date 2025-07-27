#!/bin/bash

echo "🗄️  Starting PandoraGym Database Deployment"
echo "=========================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Stop existing database container if running
echo "🛑 Stopping existing database container..."
docker-compose -f docker-compose.db.yml down

# Start database
echo "🚀 Starting PostgreSQL database..."
docker-compose -f docker-compose.db.yml up -d

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 10

# Check if database is healthy
echo "🔍 Checking database health..."
if docker-compose -f docker-compose.db.yml ps | grep -q "healthy"; then
    echo "✅ Database is healthy and ready!"
else
    echo "⚠️  Database might still be starting up. Checking connection..."
fi

# Test database connection
echo "🔌 Testing database connection..."
if docker exec pandoragym_db pg_isready -U pandoragym -d pandoragym_db > /dev/null 2>&1; then
    echo "✅ Database connection successful!"
else
    echo "❌ Database connection failed. Please check the logs."
    docker-compose -f docker-compose.db.yml logs db
    exit 1
fi

echo ""
echo "🎉 Database deployment completed!"
echo "📊 Database Details:"
echo "  - Host: localhost"
echo "  - Port: 5432"
echo "  - Database: pandoragym_db"
echo "  - User: pandoragym"
echo "  - Password: password"
echo ""
echo "🔗 Connection String:"
echo "  postgresql://pandoragym:password@localhost:5432/pandoragym_db?sslmode=disable"
echo ""
echo "🔧 Next steps:"
echo "  1. Run migrations: make migrate-up"
echo "  2. Seed database: make seed"
echo "  3. Start Go application: make run"
echo ""
echo "💡 Useful commands:"
echo "  - docker-compose -f docker-compose.db.yml logs db  # View database logs"
echo "  - docker-compose -f docker-compose.db.yml stop     # Stop database"
echo "  - docker-compose -f docker-compose.db.yml down -v  # Remove database and data"
