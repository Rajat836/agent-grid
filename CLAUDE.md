# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ontology Bot — a monorepo with three components:
- **`/app`** — Go backend (Gin HTTP + gRPC, GORM, PostgreSQL + ClickHouse)
- **`/panel`** — Next.js frontend (App Router, shadcn/Radix UI, Tailwind CSS)
- **`/proto`** — Protocol Buffer definitions (source of truth for types)
- **`/scripts`** — Bun/TypeScript utilities (proto codegen, data seeding)

## Build & Development Commands

### Backend (`/app`)
```bash
make run              # Run the app
make run-live         # Run with hot reload (requires fswatch)
make build            # Build binary
make fmt              # Format with gofmt
make vet              # Static analysis
make lint             # fmt + vet + go mod tidy
make test             # Run all tests
go test ./internal/services/... -v           # Test a specific package
go test ./internal/utils/... -run TestName   # Run a single test
make migrate-new      # Create new migration (prompts for name)
```

### Frontend (`/panel`)
```bash
npm run dev           # Dev server on :3000
npm run build         # Production build
npm run lint          # Biome check + auto-fix
npm run format        # Biome format + auto-fix
```

### Proto Code Generation (from repo root)
```bash
./scripts/gen_proto.sh    # Generates Go types → app/internal/types/, TS types → panel/src/types/
```

### Scripts (`/scripts`)
```bash
bun run populate-events                           # Seed fake events
RECORD_COUNT=1000 TPS=50 bun run populate-events  # With options
```

## Architecture

Layered/Clean Architecture in the Go backend:

```
HTTP/gRPC Request → Controller → Service → Repository / Client
```

- **Controllers** (`app/internal/controllers/`) — parse requests, validate, return responses
- **Services** (`app/internal/services/`) — business logic, orchestration
- **Repositories** (`app/internal/repositories/`) — GORM database operations
- **Clients** (`app/internal/clients/`) — external API integrations
- **Models** (`app/internal/models/`) — database entity structs

Each layer has `types.go` (interfaces) and `main.go` (initialization/DI wiring). Controllers must not access repositories directly.

### Proto-First Development

Proto definitions in `/proto` are the **single source of truth** for request/response types. The workflow:

1. Define `.proto` files in `/proto/<feature>/` (add endpoint path as comment: `// POST: /ontology_bot/v1/<feature>/<action>`)
2. Run `./scripts/gen_proto.sh` to generate Go types in `app/internal/types/<feature>/` and TS types in `panel/src/types/<feature>/`
3. Import and use generated types directly in services/controllers — **never create duplicate type definitions**
4. Import convention: `types_events "app/ontology_bot/internal/types/events"`

Files in `app/internal/types/` and `panel/src/types/` are **generated — do not edit manually**.

## Key Conventions

### File Naming
- `controller_<feature>.go`, `service_<feature>.go`, `repo_<feature>.go`
- `routes_<feature>.go`, `requests_<feature>.go`, `response_<feature>.go`, `errors_<feature>.go`
- Models: `<entity>.go`

### Interface Naming
- `Controller<Feature>Methods`, `Service<Feature>Methods`, `Repository<Feature>Methods`, `Client<Feature>Methods`

### Scope-Based Authorization
- Scope constants defined in route files: `ScopeEventsView = "events:view"`
- Format: `[feature]:[action]` (e.g., `events:view`, `events:create`, `events:manage`)
- Middleware: `middlewares.Scopes.HasAnyOneScope(...)` or `HasAllScopes(...)`

### Route Registration
- Feature routes in `cmd/app/routes_<feature>.go`
- Registered from `cmd/app/routes.go`
- Route groups: public, protected (auth required), private

### Database
- PostgreSQL for relational data, ClickHouse for event analytics
- Migrations in `app/migrations/postgres/` and `app/migrations/clickhouse/`
- Migrations run automatically on app startup
- `make migrate-new` creates up/down SQL files with timestamp prefix

### Configuration
- Copy `app/config/local.example.yml` → `app/config/local.yml` (git-ignored)
- HTTP port: 4441, gRPC port: 4442 (configurable)

### Frontend
- Biome for linting/formatting (not ESLint)
- Next.js App Router with server components
- shadcn + Radix UI component library

## Adding a New Feature

1. Create proto files in `/proto/<feature>/`
2. Run `./scripts/gen_proto.sh`
3. Add model in `app/internal/models/`
4. Implement repository → service → controller (bottom-up)
5. Register routes with scopes in `cmd/app/routes_<feature>.go`
6. Create Bruno API collection in `app/api_collection/<feature>/`
7. Add migration if schema changes needed

## CI/CD

- Bitbucket Pipelines (`bitbucket-pipelines.yml`)
- PRs to `staging`: lint + test backend & frontend
- Main branch: deploy to production
- Tags (`v*`): deploy backend service
