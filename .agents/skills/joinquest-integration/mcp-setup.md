# JoinQuest MCP setup (for agents)

Guide the developer through this once per machine. Full UI copy lives on the developer dashboard → **Connect an AI assistant**.

## Prerequisites

- Node.js 20+
- JoinQuest account (sign in at https://joinquest.cc/developers)

## Steps

1. **Generate API key** — Developer dashboard → Connect an AI assistant → **Show setup** → Generate API key. Copy once (`lq_dev_…`).

2. **Add MCP server** (uses `npx` — no install step):

### Claude Code (CLI)

```bash
claude mcp add --scope project --transport stdio \
  --env JOINQUEST_API_KEY=lq_dev_PASTE_HERE \
  joinquest-integration -- npx -y @joinquest/mcp-integration
```

Helper script (prompts for key): `bash scripts/add-joinquest-mcp-claude.sh` from the JoinQuest repo, or curl from GitHub once published.

Start a new `claude` session and approve **joinquest-integration**.

### Cursor

```json
{
  "mcpServers": {
    "joinquest-integration": {
      "type": "stdio",
      "command": "npx",
      "args": ["--yes", "--package", "@joinquest/mcp-integration", "joinquest-integration-mcp-cursor"],
      "env": {
        "JOINQUEST_API_KEY": "lq_dev_PASTE_HERE"
      }
    }
  }
}
```

### Claude Desktop

From your game repo (same one-liner as Cursor):

```bash
JOINQUEST_API_KEY=lq_dev_PASTE_HERE curl -fsSL .../install-joinquest-dev.sh | sh -s -- --claude-desktop
```

That merges MCP into Claude Desktop's config (macOS / Windows / Linux). Fully quit Claude Desktop and reopen.

Or paste manually via **Settings → Developer → Edit Config**:

```json
{
  "mcpServers": {
    "joinquest-integration": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@joinquest/mcp-integration"],
      "env": {
        "JOINQUEST_API_KEY": "lq_dev_PASTE_HERE"
      }
    }
  }
}
```

3. **Restart the agent client fully** — Cursor: **Cmd+Q**, reopen, new Agent chat. Claude Desktop: fully quit and reopen.

4. **Verify** — ask the agent: “List my JoinQuest games using the MCP tools.”

5. **Start** — say: “I want to create a game on JoinQuest.” (Same prompt works if your game is already registered — the agent continues integration from there.)

## Agent skill (recommended)

From your **game repo root**. [Read the install script on GitHub](https://github.com/scruffyprodigy/playhub/blob/main/scripts/install-joinquest-dev.sh) before running.

It copies `.agents/skills/joinquest-integration/` from GitHub; with `--cursor` merges MCP into `.cursor/mcp.json`. MCP only calls joinquest.cc when tools run.

**Review first:**

```bash
curl -fsSL https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts/install-joinquest-dev.sh -o install-joinquest-dev.sh
less install-joinquest-dev.sh
JOINQUEST_API_KEY=lq_dev_... bash install-joinquest-dev.sh --cursor
```

**Quick install:** `JOINQUEST_API_KEY=... curl -fsSL .../install-joinquest-dev.sh | sh -s -- --cursor`

## Troubleshooting (Cursor)

| Symptom | Fix |
|---------|-----|
| Disconnected / no tools | Cursor: use `--package` + `joinquest-integration-mcp-cursor` bin (see config above) |
| Auth errors | Regenerate API key on the dashboard |
| Invalid JSON | Validate `.cursor/mcp.json` |
| Still broken | Fully quit Cursor (**Cmd+Q**) |

Roadmap: [docs/developer-ai-setup-roadmap.md](../../docs/developer-ai-setup-roadmap.md)

## Available tools

See `playbook.md` § MCP tool quick reference, or MCP `joinquest_integration_get_agent_playbook`.
