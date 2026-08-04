# JoinQuest lobby — agent guide

This repository is the **JoinQuest platform** (player shell + GraphQL API + developer dashboard). It is **not** a game repo. Game logic and provision endpoints live in separate game repositories.

**Maintainers (you):** follow this file and [docs/lobby-maintenance.md](docs/lobby-maintenance.md).

**Game integration:** use [`.agents/skills/joinquest-integration/`](.agents/skills/joinquest-integration/) — typically installed into a *game* repo via `install-joinquest-dev.sh`, not the primary workflow here.

## Naming

| Term | Meaning |
|------|---------|
| `lobby` | Local folder / this repo |
| **JoinQuest** | Product name; production at [joinquest.cc](https://joinquest.cc) |
| `playhub` | Legacy internal name — still used for Docker images (`playhub-backend`), DB names (`playhub`, `playhub_test`), and the GitHub repo [`scruffyprodigy/playhub`](https://github.com/scruffyprodigy/playhub) |

Do not rename legacy `playhub` identifiers unless explicitly asked.

## Repo map

```text
backend/                          Go GraphQL API (gqlgen)
  graph/schema/*.graphqls         GraphQL schema (edit here first)
  graph/*.resolvers.go            Resolvers (hand-written)
  internal/store/                 PostgreSQL persistence
  internal/developer/             Embedded agent docs (synced copies — do not edit)
  migrations/                     SQL migrations
frontend/src/                     React + Vite player + developer UI
  lib/developers.js               Developer dashboard GraphQL client
  components/developers/          Developer portal components
mcp/joinquest-integration/        @joinquest/mcp-integration npm package
docs/                             Human docs; some are canonical for agent/MCP embeds
k8s/                              Kubernetes manifests (namespace `playhub` in files; deploy script rewrites to `joinquest`)
scripts/                          setup, test, deploy, doc sync
.agents/skills/joinquest-integration/   Agent skill source (copied to game repos + plugins)
plugins/joinquest/                Cursor/Claude plugin bundle
```

## Local development

```bash
./scripts/setup.sh
./scripts/dev.sh          # frontend :5173, GraphQL :8080/graphql
./scripts/test.sh         # backend + frontend unit tests
```

- GraphQL codegen: `cd backend && make generate` (required after schema edits)
- Backend tests only: `./scripts/test-backend.sh` (uses isolated `playhub_test` DB)
- MCP tests: `cd mcp/joinquest-integration && npm test`

See [docs/development.md](docs/development.md) for details.

## Change checklists

### Developer dashboard / GraphQL feature

Touch every layer that exposes the same operation:

1. `backend/graph/schema/developer.graphqls`
2. `cd backend && make generate`
3. `backend/internal/store/developer.go` (+ `*_test.go`)
4. `backend/graph/developer.resolvers.go`
5. `frontend/src/lib/developers.js` (+ component tests)
6. `frontend/src/components/developers/*.jsx` as needed
7. `mcp/joinquest-integration/src/graphql.js`
8. `mcp/joinquest-integration/src/tools.js`
9. Bump `mcp/joinquest-integration/package.json` version when MCP surface changes
10. `docs/developer-self-service.md` and/or `docs/developer-agent-playbook.md`
11. `./scripts/sync-developer-docs.sh`
12. Run tests (backend, frontend, MCP)

### Playbook or integration guide only

1. Edit **canonical** file in `docs/`:
   - `docs/developer-agent-playbook.md`
   - `docs/developer-integration-guide.md`
2. Run `./scripts/sync-developer-docs.sh` (updates backend embeds, skill copies, plugin copies)
3. CI enforces sync via `./scripts/check-developer-docs-sync.sh`

Never edit `backend/internal/developer/*.md` or `.agents/.../playbook.md` directly.

### GraphQL schema (any domain)

1. Edit `backend/graph/schema/*.graphqls`
2. `cd backend && make generate`
3. Implement resolver + store changes
4. Add/update tests; CI runs gqlgen drift detection

### Database schema

1. Add migration under `backend/migrations/`
2. `cd backend && make migrate-up` locally
3. Update `backend/internal/store/` + tests

## Shipping (human approval required)

Do **not** commit, push, deploy, or publish unless the user explicitly asks.

Pre-ship checks:

```bash
./scripts/ship-joinquest.sh --check
```

Production deploy ([joinquest.cc](https://joinquest.cc)):

```bash
./scripts/build-and-push.sh --push
gcloud container clusters get-credentials joinquest --region us-east1 --project joinquest-demo  # refresh if Unauthorized
./scripts/deploy-joinquest.sh
```

MCP npm publish (only when `mcp/` changed and version bumped):

```bash
./scripts/publish-joinquest-mcp.sh
```

Full runbook: [docs/lobby-maintenance.md](docs/lobby-maintenance.md).

## Conventions

- **Minimal scope** — smallest correct diff; don't refactor unrelated code
- **Match existing style** — Go: `gofmt`, store layer for DB; frontend: JSX + Vitest
- **No secrets** — never commit `.env`, API keys, k8s secrets, or `.DS_Store`
- **Tests** — add store/resolver tests for backend logic; component/lib tests for dashboard changes
- **Human gates** — production deploy, npm publish, and git push require explicit user approval

## Key docs

| Topic | Doc |
|-------|-----|
| Maintainer runbook | [docs/lobby-maintenance.md](docs/lobby-maintenance.md) |
| Local dev | [docs/development.md](docs/development.md) |
| Architecture | [docs/architecture.md](docs/architecture.md) |
| Testing | [docs/testing.md](docs/testing.md) |
| Developer self-service spec | [docs/developer-self-service.md](docs/developer-self-service.md) |
| Agent playbook (canonical) | [docs/developer-agent-playbook.md](docs/developer-agent-playbook.md) |
| Integration guide (canonical) | [docs/developer-integration-guide.md](docs/developer-integration-guide.md) |
