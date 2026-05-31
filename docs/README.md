# JoinQuest Documentation

This directory contains documentation for the JoinQuest gaming lobby platform.

## 📚 Documentation Structure

### Getting started
- **[Development Guide](development.md)** — Local setup and day-to-day workflow
- **[Environment Configuration](environment-configuration.md)** — Runtime env injection (local, k8s)
- **[Testing Guide](testing.md)** — How to run backend, frontend, and E2E tests

### Architecture & API
- **[Architecture Overview](architecture.md)** — System design and components
- **[Game Catalog Architecture](game-catalog-architecture.md)** — Modes, queues, manifest sync
- **[Lobby ↔ Game Handoff](lobby-protocol-handoff.md)** — Provision, JWT, and integration rationale
- **[API Documentation](api.md)** — GraphQL reference and examples
- **[Pub/Sub](pubsub.md)** — Redis queue notifications

### Operations
- **[Database Migrations](database-migrations.md)** — Schema migrations and CLI
- **[Contributing Guide](contributing.md)** — How to contribute

## 🚀 Quick Start

1. **Prerequisites**: Go 1.25+, Node.js 20+, Docker
2. **Clone**: `git clone https://github.com/scruffyprodigy/playhub.git`
3. **Setup**: `./scripts/setup.sh` from the project root
4. **Start**: `./scripts/dev.sh` (frontend `:5173`, GraphQL `:8080/graphql`)

## 📁 Project Structure

```
playhub/                 # repository root (module name unchanged)
├── backend/             # Go GraphQL API
├── frontend/            # React + Vite application
├── k8s/                 # Kubernetes manifests
├── docs/                # This documentation
├── scripts/             # Dev, test, and deploy scripts
└── .github/workflows/   # CI/CD
```

## 🔗 External Links

- [GraphQL Schema](../backend/graph/schema/)
- [Frontend source](../frontend/src/)
- [Kubernetes configs](../k8s/)
- [GitHub Actions](../.github/workflows/)
