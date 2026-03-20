# App

Ontology Bot Backend Service

Built with:
- Go 1.24
- Gin Framework
- gRPC
- GORM
- PostgreSQL & ClickHouse

## Prerequisites

- Go 1.24+
- PostgreSQL 12+
- Redis (optional for local development)
- Make (for running commands)

## Setup

### 1. Install Dependencies

```bash
go mod download
go mod tidy
```

### 2. Configure the Application

Copy the example configuration file:

```bash
cp config/local.example.yml config/local.yml
```

Edit `config/local.yml` and update values as needed:

```yaml
ServerPort: 4441          # HTTP server port
GrpcPort: 4442           # gRPC server port
Environment: "local"     # local, stg, production

Database:
  master_database_dsn: "postgres://user:password@localhost:5432/ontology_bot?sslmode=disable"
  slave_database_dsn: "postgres://user:password@localhost:5432/ontology_bot?sslmode=disable"

Redis:
  host: "localhost"
  port: 6379
```

### 3. Set Up Database

Create the PostgreSQL database:

```bash
createdb -U postgres ontology_bot
```

(Or update the DSN in `config/local.yml` with your database credentials)

## Running the Backend Service

### Option 1: Using Make (Recommended)

Start the development server:

```bash
make run
```

The server will start on:
- **HTTP**: http://localhost:4441
- **gRPC**: localhost:4442

### Option 2: Live Reload (with fswatch)

For hot reload on file changes:

```bash
make run-live
```

**Note**: Requires `fswatch` to be installed. Install with:
- macOS: `brew install fswatch`
- Ubuntu/Debian: `sudo apt-get install fswatch`
- Other: See [fswatch documentation](https://emcrisostomo.github.io/fswatch/)

### Option 3: Direct Go Command

```bash
go run cmd/server/main.go
```

### Option 4: Build and Run Binary

```bash
make build
./server
```

## Health Endpoints

### Liveness Check
```bash
curl http://localhost:4441/ontology_bot/v1/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2026-03-19T10:30:00Z",
  "version": "1.0.0"
}
```

### Readiness Check
```bash
curl http://localhost:4441/ontology_bot/v1/health/ready
```

Response:
```json
{
  "ready": true,
  "timestamp": "2026-03-19T10:30:00Z"
}
```

## Development Commands

### Code Quality

```bash
make fmt           # Format code with gofmt
make vet           # Run static analysis
make lint          # Format + vet + tidy modules
```

### Testing

```bash
make test                                    # Run all tests
go test ./internal/services/... -v          # Test specific package
go test ./internal/utils/... -run TestName  # Run single test
```

### Database Migrations

```bash
make migrate-new  # Create new migration (prompts for name)
```

Creates migration files in:
- PostgreSQL: `migrations/postgres/`
- ClickHouse: `migrations/clickhouse/`

### cleaning Up

```bash
make clean  # Remove temporary files and binaries
```

## Project Structure

```
app/
├── cmd/
│   ├── server/          # Entry point (main.go)
│   └── app/             # Routes registration
├── internal/
│   ├── controllers/     # HTTP/gRPC handlers
│   ├── services/        # Business logic
│   ├── repositories/    # Database access
│   ├── models/          # Database entities
│   ├── types/           # Generated Protocol Buffer types
│   ├── config/          # Configuration management
│   └── clients/         # External API clients
├── migrations/
│   ├── postgres/        # PostgreSQL migrations
│   └── clickhouse/      # ClickHouse migrations
└── config/
    └── local.yml        # Local configuration (git-ignored)
```

## Architecture

The backend follows a **layered/clean architecture pattern**:

```
HTTP/gRPC Request
    ↓
Controller (parse, validate input)
    ↓
Service (business logic, orchestration)
    ↓
Repository / Client (database or external APIs)
```

**Key Rule**: Controllers must NOT access repositories directly—only through services.

## Adding a New Feature

1. Create proto files in `/proto/<feature>/`
2. Run `./scripts/gen_proto.sh` to generate types
3. Add model in `app/internal/models/<entity>.go`
4. Implement **bottom-up**: repository → service → controller
5. Register routes in `cmd/app/routes_<feature>.go`
6. Add migrations if needed: `make migrate-new`

See [../CLAUDE.md](../CLAUDE.md) for detailed conventions and architecture guidelines.
