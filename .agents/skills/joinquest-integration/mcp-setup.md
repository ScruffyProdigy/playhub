# JoinQuest MCP setup (for agents)

Guide the developer through this once per machine. Full UI copy lives on the developer dashboard → **Connect an AI assistant**.

## Prerequisites

- Node.js 20+
- JoinQuest account (sign in at https://joinquest.cc/developers)

## Steps

1. **Generate API key** — Developer dashboard → Connect an AI assistant → Generate API key. Copy once (`lq_dev_…`).

2. **Add MCP server** — paste into the agent client's MCP config:

```json
{
  "mcpServers": {
    "joinquest-integration": {
      "command": "npx",
      "args": ["-y", "@joinquest/mcp-integration"],
      "env": {
        "JOINQUEST_API_URL": "https://joinquest.cc/graphql",
        "JOINQUEST_API_KEY": "lq_dev_PASTE_HERE"
      }
    }
  }
}
```

Local dev: `JOINQUEST_API_URL` → `http://localhost:8080/graphql`

3. **Restart** the agent client (Cursor, Claude Code, Codex, Claude Desktop, etc.).

4. **Verify** — call `joinquest_integration_list_my_games`. Auth error → recheck key and URL.

## Config file locations

| Client | Path |
|--------|------|
| Cursor | `~/.cursor/mcp.json` or `.cursor/mcp.json` |
| Claude Desktop | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| Claude Code | `.mcp.json` at project root |
| Codex / cross-platform | `.mcp.json` or client-specific path |

Claude Code `.mcp.json` example (add `"type": "stdio"` on the server entry):

```json
{
  "mcpServers": {
    "joinquest-integration": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@joinquest/mcp-integration"],
      "env": {
        "JOINQUEST_API_URL": "https://joinquest.cc/graphql",
        "JOINQUEST_API_KEY": "lq_dev_PASTE_HERE"
      }
    }
  }
}
```

## Alternative install

```bash
curl -fsSL https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts/install-joinquest-mcp.sh | sh
```

Then use `"command": "joinquest-integration-mcp"` with the same `env` block.

## Available tools

See `playbook.md` § MCP tool quick reference, or MCP `joinquest_integration_get_agent_playbook`.
