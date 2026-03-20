# Ontology Bot

Ontology Bot is a monorepo project consisting of:
- **Backend** (`/app`) — Go backend service with Gin HTTP + gRPC
- **Frontend** (`/panel`) — Next.js frontend application
- **Protocol Buffers** (`/proto`) — Shared type definitions
- **Scripts** (`/scripts`) — Utility scripts for code generation and seeding

## Quick Start

### Backend Setup
```bash
cd app
cp config/local.example.yml config/local.yml
# Edit config/local.yml with your settings
make run
```

### Frontend Setup
```bash
cd panel
npm install
npm run dev
```

### Proto Code Generation
```bash
./scripts/gen_proto.sh
```

## Project Structure

See [CLAUDE.md](./CLAUDE.md) for detailed architecture and conventions.
