# JoinQuest integration MCP server

stdio MCP server for **game developers** integrating with JoinQuest. Exposes developer dashboard operations to Cursor, Claude, and other agents via the same GraphQL API as the portal.

## Quick setup (recommended)

1. Sign in at [JoinQuest](https://joinquest.cc) and open your **developer dashboard**.
2. Expand **Connect an AI assistant** → **Generate API key** (copy it once).
3. Add this to Cursor (`.cursor/mcp.json`) or Claude Desktop config:

```json
{
  "mcpServers": {
    "joinquest-integration": {
      "command": "npx",
      "args": ["-y", "@joinquest/mcp-integration"],
      "env": {
        "JOINQUEST_API_URL": "https://joinquest.cc/graphql",
        "JOINQUEST_API_KEY": "lq_dev_..."
      }
    }
  }
}
```

Node.js 20+ required. The first run downloads the package via `npx`.

### Alternative: global CLI install

```bash
curl -fsSL https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts/install-joinquest-mcp.sh | sh
```

Then use `"command": "joinquest-integration-mcp"` with the same `env` block (no `args`).

## Environment

| Variable | Required | Description |
|----------|----------|-------------|
| `JOINQUEST_API_KEY` | Yes* | From developer dashboard → Connect AI assistant |
| `JOINQUEST_API_URL` | No | GraphQL URL (default `http://localhost:8080/graphql`) |
| `JOINQUEST_ISSUER_URL` | No | Lobby JWT issuer / provision `lobbyId` |
| `JOINQUEST_PUBLIC_URL` | No | Browser Lobby URL for example provision payloads |
| `JOINQUEST_SESSION` | Legacy | Session cookie (prefer API key) |

## Client-specific paths

| Client | Config file |
|--------|-------------|
| Cursor | `~/.cursor/mcp.json` or `.cursor/mcp.json` |
| Claude Desktop | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| Claude Code | `.mcp.json` at project root (add `"type": "stdio"`) |

## Tools

| Tool | Purpose |
|------|---------|
| `joinquest_integration_get_agent_playbook` | **Start here** — end-to-end agent workflow (phases 1–8) |
| `joinquest_integration_get_integration_guide` | Full integration guide (markdown) |
| `joinquest_integration_get_discovery_prompt` | Agent interview script |
| `joinquest_integration_get_catalog_tag_taxonomy` | Valid tag IDs |
| `joinquest_integration_list_my_games` | Owner's games + visibility |
| `joinquest_integration_get_game_checks` | Checklist + metadata |
| `joinquest_integration_run_game_checks` | Run manifest / provision / JWT suite |
| `joinquest_integration_update_game_metadata` | Save catalog copy + tags |
| `joinquest_integration_get_game_credentials` | serviceToken + webhook secret |
| `joinquest_integration_get_example_provision_payload` | Sample provision JSON |
| `joinquest_integration_request_public_release` | Submit for catalog review |

## Run manually

```bash
cd mcp/joinquest-integration
npm install
JOINQUEST_API_KEY=your-key node src/index.js
```

## Tests

```bash
npm test
```
