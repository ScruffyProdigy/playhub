# Scripts

## Game developers (JoinQuest integration)

Install the agent skill + MCP in your **game repo** or as a Cursor/Claude plugin:

→ **[joinquest-setup/README.md](joinquest-setup/README.md)** — complete file list, review-before-run, plugin vs project install.

Quick start:

```bash
JOINQUEST_API_KEY=lq_dev_... npx joinquest install cursor
# or
npm create joinquest@latest -- --cursor
```

See [packages/joinquest/README.md](../packages/joinquest/README.md).

Legacy `curl | sh` wrappers forward to `npx joinquest` after publish.

## Lobby contributors (this repo)

Local dev, test, deploy — see [AGENTS.md](../AGENTS.md) and [docs/lobby-maintenance.md](../docs/lobby-maintenance.md).

| Script | Purpose |
|--------|---------|
| `setup.sh` / `dev.sh` | Local stack |
| `test.sh` / `test-backend.sh` / `test-frontend.sh` | Tests |
| `ship-joinquest.sh --check` | Pre-ship checks |
| `deploy-joinquest.sh` | Deploy to joinquest.cc (GKE) |
| `build-and-push.sh --push` | Docker images |
| `publish-joinquest-mcp.sh` | npm publish `@joinquest/mcp-integration` |
| `sync-developer-docs.sh` | Sync playbook/guide embeds |
