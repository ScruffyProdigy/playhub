# Testing Guide

Testing strategies for the JoinQuest platform ([product vision](./vision.md)): API, player UI, and game-integration paths (catalog, queues, handoff).

## Testing Status

### ✅ Implemented
- **Frontend unit tests**: Vitest (components, hooks, lib)
- **Frontend integration tests**: App-level flows with mocked GraphQL
- **Frontend E2E tests**: Playwright against a running stack
- **Backend unit tests**: Resolvers, mappers, auth, gameclient
- **Backend integration tests**: PostgreSQL-backed tests for auth, catalog, queue/handoff, games, player lookup
- **Environment configuration tests**: Runtime `window.env` validation
- **Drift detection**: gqlgen schema vs generated code (`gqlgen_drift_test.go`)

### 🚧 In Development
- **Performance tests**: Load testing and benchmarking

### 📋 Planned
- **API contract tests**: Schema validation against external consumers
- **Security tests**: Vulnerability and penetration testing

## Testing Philosophy

JoinQuest uses a layered testing strategy:

- **Unit tests**: Individual functions and components
- **Integration tests**: Real database (backend) or mocked GraphQL (frontend)
- **End-to-end tests**: Full browser workflows against dev stack

## Backend Testing

### Test Structure

```
backend/
├── graph/
│   ├── resolvers_test.go              # Resolver unit tests (no DB)
│   ├── healthz_test.go                # Basic GraphQL smoke tests
│   ├── gqlgen_drift_test.go           # Schema / codegen drift detection
│   ├── auth_integration_test.go       # Magic-link auth + session
│   ├── catalog_integration_test.go    # registerGame + manifest sync
│   ├── queue_integration_test.go      # joinQueue, subscriptions, handoff
│   ├── handoff_integration_test.go    # Provision payload to game API
│   ├── games_integration_test.go      # Catalog games + sessions
│   └── player_integration_test.go     # Game-service player lookup
├── internal/
│   ├── store/*_test.go                # Store layer (uses playhub_test DB)
│   └── auth/*_test.go                 # JWT, cookies, issuer URLs
```

Integration tests require PostgreSQL. `./scripts/test-backend.sh` migrates and uses the isolated `playhub_test` database.

### Running Backend Tests

```bash
# Recommended: isolated test DB
./scripts/test-backend.sh

# All packages (requires TEST_DATABASE_URL or test-backend.sh setup)
cd backend && go test ./...

# GraphQL package only
cd backend && go test ./graph -v

# With coverage
cd backend && go test -cover ./...
```

## Frontend Testing

### Running Frontend Tests

```bash
# Unit and integration tests
cd frontend && npm run test:run

# Watch mode
cd frontend && npm run test:watch

# E2E tests (backend + frontend must be running)
cd frontend && npm run test:e2e
```

See [frontend/README_TESTING.md](../frontend/README_TESTING.md) for component-level detail.

## CI/CD Testing

GitHub Actions run backend tests, frontend unit tests, gqlgen drift checks, and E2E where configured. Path filters limit runs to relevant changes.

### Test Scripts

```bash
./scripts/test.sh           # All tests
./scripts/test-backend.sh   # Backend (playhub_test DB)
./scripts/test-frontend.sh  # Frontend unit tests
```

## Best Practices

1. **Backend integration tests** must not mutate the dev `playhub` database — use `playhub_test` via `./scripts/test-backend.sh`.
2. **Queue/handoff tests** restore seeded demo game URLs after runs (`RestoreDemoGameHandoffURLs`).
3. After GraphQL schema changes, run `go run github.com/99designs/gqlgen@v0.17.81 generate` in `backend/` and commit generated files.
4. **Frontend tests** mock `fetch` / GraphQL in `src/test/setup.js`; E2E uses the real API.

## Debugging Tests

### Backend

```bash
go test -v ./graph -run TestJoinQueue
```

### Frontend

```bash
cd frontend && npm run test:watch
cd frontend && npm run test:e2e:headed
```
