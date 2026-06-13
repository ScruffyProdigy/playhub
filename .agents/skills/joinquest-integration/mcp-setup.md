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

### Claude Desktop / other stdio clients

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

3. **Restart the agent client fully** — Cursor: **Cmd+Q**, reopen, new Agent chat.

4. **Verify** — ask the agent to call `joinquest_integration_list_my_games`.

## Agent skill (recommended)

Copy `.agents/skills/joinquest-integration/` into the **game repo** at the same path. See dashboard **Connect an AI assistant** or the [JoinQuest repo](https://github.com/scruffyprodigy/playhub/tree/main/.agents/skills/joinquest-integration).

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
