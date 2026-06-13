#!/usr/bin/env bash
# Write JoinQuest MCP config for Cursor / Claude Code.
set -euo pipefail

joinquest_cursor_mcp_json() {
  cat <<EOF
{
  "mcpServers": {
    "joinquest-integration": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "--yes",
        "--package",
        "@joinquest/mcp-integration",
        "joinquest-integration-mcp-cursor"
      ],
      "env": {
        "JOINQUEST_API_KEY": "$JOINQUEST_API_KEY"
      }
    }
  }
}
EOF
}

joinquest_write_cursor_mcp() {
  local config_path="${1:-.cursor/mcp.json}"

  if ! command -v node >/dev/null 2>&1; then
    echo "Node.js is required to write $config_path (merge with existing MCP servers)." >&2
    echo "Paste this into $config_path:" >&2
    joinquest_cursor_mcp_json
    return 1
  fi

  mkdir -p "$(dirname "$config_path")"
  CONFIG_PATH="$config_path" JOINQUEST_API_KEY="$JOINQUEST_API_KEY" node <<'NODE'
const fs = require('fs')
const path = process.env.CONFIG_PATH
const server = {
  type: 'stdio',
  command: 'npx',
  args: ['--yes', '--package', '@joinquest/mcp-integration', 'joinquest-integration-mcp-cursor'],
  env: { JOINQUEST_API_KEY: process.env.JOINQUEST_API_KEY },
}
let config = { mcpServers: {} }
if (fs.existsSync(path)) {
  config = JSON.parse(fs.readFileSync(path, 'utf8'))
  if (!config.mcpServers || typeof config.mcpServers !== 'object') {
    config.mcpServers = {}
  }
}
config.mcpServers['joinquest-integration'] = server
fs.writeFileSync(path, `${JSON.stringify(config, null, 2)}\n`)
NODE
  echo "Wrote joinquest-integration MCP to $config_path"
}

joinquest_setup_claude_mcp() {
  if ! command -v claude >/dev/null 2>&1; then
    echo "claude CLI not found — install Claude Code or run this manually:" >&2
    joinquest_print_claude_command
    return 1
  fi

  if ! command -v npx >/dev/null 2>&1; then
    echo "npx not found. Install Node.js 20+." >&2
    return 1
  fi

  claude mcp add --scope project --transport stdio \
    --env "JOINQUEST_API_KEY=$JOINQUEST_API_KEY" \
    joinquest-integration -- npx -y @joinquest/mcp-integration

  echo "Added joinquest-integration to .mcp.json (project scope)."
}

joinquest_print_claude_command() {
  cat <<EOF
claude mcp add --scope project --transport stdio \\
  --env JOINQUEST_API_KEY=$JOINQUEST_API_KEY \\
  joinquest-integration -- npx -y @joinquest/mcp-integration
EOF
}
