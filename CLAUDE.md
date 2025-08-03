# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building and Running
- `go run main.go` - Run the application locally (requires DATABASE_URL environment variable)
- `go build` - Build the binary
- `go mod tidy` - Clean up module dependencies

### Testing
- `go test ./...` - Run all tests in the project
- `go test ./api/handlers` - Run tests for API handlers specifically
- `go test ./utils` - Run utility function tests
- `go test -v ./...` - Run tests with verbose output

### Code Quality
- `go fmt ./...` - Format all Go code
- `go vet ./...` - Run Go's built-in static analysis
- `go mod verify` - Verify module dependencies

## Architecture Overview

This is a URL shortener service built in Go with PostgreSQL as the database. The application follows a clean layered architecture:

### Core Components
- **main.go**: Entry point that sets up HTTP routes, initializes database and cron jobs
- **api/**: HTTP API layer with a single `/api/link` endpoint that handles both POST (create) and GET (retrieve) operations
- **api/handlers/**: Request handlers for creating short URLs and retrieving original URLs
- **db/**: Database layer with connection management and SQL query execution
- **utils/**: Utility functions for URL validation and random string generation
- **cron/**: Background job system for automatic link cleanup
- **public/**: Static frontend assets (HTML, CSS, JS)

### Database Design
- Uses PostgreSQL with connection pooling via pgx/v4
- SQL queries are stored as separate files in `db/sql/`
- Database initialization and migration handled automatically on startup

### Key Patterns
- Handler functions return (status int, data string) tuples
- SQL queries are loaded once at startup and reused
- Random hash generation uses 5-character strings (~380M possible combinations)
- URL validation prevents malicious redirects and ensures proper formatting

### Environment Requirements
- `DATABASE_URL`: PostgreSQL connection string (required)
- `PORT`: Server port (defaults to 8080)
- `.env` file support via godotenv

### API Endpoints
- `POST /api/link` - Create shortened URL (expects JSON: `{"Link": "https://example.com"}`)
- `GET /api/link?hash=XXXXX` - Retrieve original URL by hash
- `/` - Serves static frontend
- `/static/*` - Serves static assets from public/ directory