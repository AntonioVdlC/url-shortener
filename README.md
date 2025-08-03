# URL Shortener

A fast, lightweight URL shortener service built with Go and PostgreSQL. Features automatic link cleanup, secure URL validation, and a simple web interface.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI](https://github.com/avdlc/url-shortener/actions/workflows/test.yml/badge.svg)](https://github.com/avdlc/url-shortener/actions)

## ✨ Features

- **Fast & Lightweight**: Built with Go for optimal performance
- **Secure**: URL validation prevents malicious redirects
- **Auto-cleanup**: Automatic removal of links older than 30 days
- **Clean Architecture**: Modular design with clear separation of concerns
- **Database Persistence**: PostgreSQL with connection pooling
- **Web Interface**: Simple, responsive frontend
- **Docker Ready**: Complete development environment with Docker Compose
- **CI/CD**: GitHub Actions with automated testing and quality checks

## 🚀 Quick Start

### Prerequisites

- [Go 1.23+](https://golang.org/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/) (optional, for convenience commands)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/url-shortener.git
   cd url-shortener
   ```

2. **Set up development environment**
   ```bash
   make dev-setup
   ```
   This will:
   - Start a PostgreSQL database in Docker
   - Create a `.env` file with proper configuration
   - Wait for the database to be ready

3. **Start the application**
   ```bash
   make dev-start
   ```

4. **Visit the application**
   - Open http://localhost:8080 in your browser
   - The database will be automatically initialized on first run

### Manual Setup (without Make)

1. **Start PostgreSQL database**
   ```bash
   docker compose up -d postgres
   ```

2. **Create environment file**
   ```bash
   cp .env.example .env
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

## 📖 API Documentation

### Create Short URL

**POST** `/api/link`

Create a shortened URL from a long URL.

**Request Body:**
```json
{
  "Link": "https://example.com/very/long/url"
}
```

**Response (201 Created):**
```json
{
  "hash": "abc123"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON or missing URL
- `403 Forbidden`: Invalid URL format or unsafe URL
- `500 Internal Server Error`: Database error

### Retrieve Original URL

**GET** `/api/link?hash=abc123`

Retrieve the original URL from a hash.

**Response (200 OK):**
```json
{
  "url": "https://example.com/very/long/url"
}
```

**Error Responses:**
- `404 Not Found`: Hash not found or expired
- `500 Internal Server Error`: Database error

### Example Usage

```bash
# Create a short URL
curl -X POST http://localhost:8080/api/link \
  -H "Content-Type: application/json" \
  -d '{"Link": "https://github.com/golang/go"}'

# Response: {"hash":"abc123"}

# Retrieve the original URL
curl http://localhost:8080/api/link?hash=abc123

# Response: {"url":"https://github.com/golang/go"}
```

## 🛠️ Development

### Available Commands

| Command | Description |
|---------|-------------|
| `make dev-setup` | Set up development environment |
| `make dev-start` | Start the application |
| `make dev-stop` | Stop development database |
| `make test` | Run tests with mock database |
| `make test-integration` | Run integration tests (requires DB) |
| `make build` | Build application binary |
| `make fmt` | Format Go code |
| `make vet` | Run static analysis |
| `make install-hooks` | Install Git pre-commit hooks |

### Database Management

| Command | Description |
|---------|-------------|
| `make db-shell` | Connect to development database |
| `make db-logs` | Show database logs |
| `make pgadmin` | Start pgAdmin web interface |

### Project Structure

```
├── api/                 # HTTP API layer
│   └── handlers/        # Request handlers
├── cron/                # Background jobs
├── db/                  # Database layer
│   └── sql/             # SQL queries
├── public/              # Static web assets
├── scripts/             # Development scripts
│   └── hooks/           # Git hooks
├── utils/               # Utility functions
├── docker-compose.yml   # Development environment
├── Makefile            # Development commands
└── main.go             # Application entry point
```

### Code Quality

This project uses automated code quality checks:

- **Pre-commit hooks**: Automatically run tests, formatting, and linting
- **GitHub Actions**: CI/CD pipeline with comprehensive checks
- **Go standards**: `go fmt`, `go vet`, and module verification

Install pre-commit hooks:
```bash
make install-hooks
```

## 🧪 Testing

### Unit Tests
```bash
# Run tests with mock database
make test

# Run specific package tests
go test ./utils -v
go test ./api/handlers -v
```

### Integration Tests
```bash
# Requires running PostgreSQL database
make test-integration
```

### Test Coverage
```bash
go test -cover ./...
```

## 🏗️ Architecture

### Core Components

- **main.go**: Application entry point and HTTP server setup
- **api/**: REST API implementation with clean handler pattern
- **db/**: Database abstraction layer with connection pooling
- **utils/**: Reusable utilities for URL validation and random generation
- **cron/**: Background job system for automatic cleanup

### Key Design Decisions

- **5-character hashes**: Provides ~380M unique combinations
- **PostgreSQL**: Reliable persistence with ACID compliance
- **Connection pooling**: Efficient database resource usage
- **Modular architecture**: Clear separation of concerns
- **Security first**: URL validation prevents malicious redirects

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `PORT` | HTTP server port | `8080` |

## 🚢 Deployment

### Docker Production

```bash
# Build production image
docker build -t url-shortener .

# Run with external PostgreSQL
docker run -e DATABASE_URL="postgres://..." -p 8080:8080 url-shortener
```

### Environment Setup

1. Set up PostgreSQL database
2. Configure environment variables
3. Deploy application binary or container

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Install pre-commit hooks** (`make install-hooks`)
4. **Make your changes** with proper tests
5. **Commit your changes** (`git commit -m 'Add amazing feature'`)
6. **Push to the branch** (`git push origin feature/amazing-feature`)
7. **Open a Pull Request**

### Development Guidelines

- Follow Go conventions and idioms
- Write tests for new functionality
- Update documentation for API changes
- Ensure all CI checks pass
- Use descriptive commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/url-shortener/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/url-shortener/discussions)
- **Documentation**: See [CLAUDE.md](CLAUDE.md) for detailed development guidance

---

**Built with _boredom_ using Go and PostgreSQL**