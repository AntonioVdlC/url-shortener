#!/bin/bash

# Development setup script for url-shortener
# Sets up local PostgreSQL database and environment

set -e

echo "🚀 Setting up url-shortener development environment..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${2}${1}${NC}"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_status "❌ Docker is not installed. Please install Docker first." "$RED"
    exit 1
fi

# Check if Docker Compose is available
if ! docker compose version &> /dev/null; then
    print_status "❌ Docker Compose is not available. Please install Docker Compose." "$RED"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
    print_status "📝 Creating .env file from .env.example..." "$YELLOW"
    cp .env.example .env
    print_status "✅ .env file created" "$GREEN"
else
    print_status "ℹ️  .env file already exists" "$BLUE"
fi

# Start PostgreSQL database
print_status "🐘 Starting PostgreSQL database..." "$YELLOW"
docker compose up -d postgres

# Wait for database to be ready
print_status "⏳ Waiting for database to be ready..." "$YELLOW"
until docker compose exec postgres pg_isready -U urluser -d url_shortener > /dev/null 2>&1; do
    sleep 1
done

print_status "✅ PostgreSQL database is ready!" "$GREEN"

# Display connection information
print_status "📋 Development environment ready!" "$GREEN"
echo ""
print_status "Database connection details:" "$BLUE"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: url_shortener"
echo "  Username: urluser"
echo "  Password: urlpass"
echo ""
print_status "To start the application:" "$BLUE"
echo "  make dev-start"
echo ""
print_status "To stop the database:" "$BLUE"
echo "  make dev-stop"
echo ""
print_status "To start pgAdmin (optional):" "$BLUE"
echo "  make pgadmin"
echo "  Then visit: http://localhost:8081 (admin@admin.com / admin)"