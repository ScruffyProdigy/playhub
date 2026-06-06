# JoinQuest Documentation

Documentation for **JoinQuest** — the platform that connects players to third-party web games and handles lobby, auth, and (roadmap) commerce so developers can focus on their game.

**New here?** Read **[Product vision](vision.md)** first.

## Who should read what

| You are… | Start with |
|----------|------------|
| Evaluating or partnering on a game | [vision.md](vision.md) → [lobby-protocol-handoff.md](lobby-protocol-handoff.md) |
| Implementing game integration | [lobby-protocol-handoff.md](lobby-protocol-handoff.md), [game-catalog-architecture.md](game-catalog-architecture.md), [seat-templates-and-matchmaking.md](seat-templates-and-matchmaking.md) |
| Hacking on this repo | [development.md](development.md), [architecture.md](architecture.md), [contributing.md](contributing.md) |

## Documentation structure

### Product & integration
- **[Product vision](vision.md)** — Problem, promise, mental model, principles
- **[End-to-end partner checklist](end-to-end-partner-checklist.md)** — Sign-in → queue → play → return
- **[Match lifecycle callbacks](match-lifecycle-callbacks.md)** — Per-player finish + match result (planned)
- **[Composition & join options](composition-and-join-options.md)** — Role queues vs game-backed eligibility
- **[Rooms & tables](rooms-and-tables.md)** — Private rooms, forming tables, invite/QR, queue exclusion
- **[Spirit animal avatars](spirit-animal-avatars.md)** — Planned: starter avatars + tarot-guided “Find my spirit animal” flow (TODO)
- **[Lobby ↔ Game Handoff](lobby-protocol-handoff.md)** — Provision, JWT, and integration rationale
- **[Game Catalog Architecture](game-catalog-architecture.md)** — Modes, queues, manifest sync
- **[Seat templates & LFG](seat-templates-and-matchmaking.md)** — Seat map contract and matchmaking spec

### Getting started (contributors)
- **[Development Guide](development.md)** — Local setup and day-to-day workflow
- **[Environment Configuration](environment-configuration.md)** — Runtime env injection (local, k8s)
- **[Testing Guide](testing.md)** — Backend, frontend, and E2E tests

### Architecture & API
- **[Architecture Overview](architecture.md)** — System design and components
- **[API Documentation](api.md)** — GraphQL reference and examples
- **[Pub/Sub](pubsub.md)** — Redis queue notifications

### Operations
- **[Database Migrations](database-migrations.md)** — Schema migrations and CLI
- **[Contributing Guide](contributing.md)** — How to contribute

## Quick start

1. **Prerequisites**: Go 1.25+, Node.js 20+, Docker
2. **Clone**: `git clone https://github.com/scruffyprodigy/playhub.git`
3. **Setup**: `./scripts/setup.sh` from the project root
4. **Start**: `./scripts/dev.sh` (frontend `:5173`, GraphQL `:8080/graphql`)

## Project structure

```text
playhub/                 # repository root (module name unchanged)
├── backend/             # Go GraphQL API (Lobby / JoinQuest server)
├── frontend/            # React + Vite — player-facing JoinQuest UI
├── k8s/                 # Kubernetes manifests
├── docs/                # This documentation
├── scripts/             # Dev, test, and deploy scripts
└── .github/workflows/   # CI/CD
```

## External links

- [GraphQL Schema](../backend/graph/schema/)
- [Frontend source](../frontend/src/)
- [Kubernetes configs](../k8s/)
- [GitHub Actions](../.github/workflows/)
