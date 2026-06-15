#!/usr/bin/env bash
# Write JoinQuest MCP config for Cursor, Claude, Copilot, Windsurf, Cline, and Roo Code.
set -euo pipefail

joinquest_stdio_mcp_server_json() {
  local use_cursor_bin=${1:-false}
  if [ "$use_cursor_bin" = true ]; then
    cat <<EOF
{
  "type": "stdio",
  "command": "npx",
  "args": ["--yes", "--package", "@joinquest/mcp-integration", "joinquest-integration-mcp-cursor"],
  "env": { "JOINQUEST_API_KEY": "$JOINQUEST_API_KEY" }
}
EOF
  else
    cat <<EOF
{
  "type": "stdio",
  "command": "npx",
  "args": ["-y", "@joinquest/mcp-integration"],
  "env": { "JOINQUEST_API_KEY": "$JOINQUEST_API_KEY" }
}
EOF
  fi
}

joinquest_windsurf_stdio_mcp_server_json() {
  cat <<EOF
{
  "type": "stdio",
  "command": "npx",
  "args": ["-y", "@joinquest/mcp-integration"],
  "env": { "JOINQUEST_API_KEY": "\${env:JOINQUEST_API_KEY}" }
}
EOF
}

joinquest_merge_mcp_config() {
  local config_path=$1
  local root_key=$2
  local use_cursor_bin=${3:-false}
  local use_windsurf_env=${4:-false}

  if ! command -v node >/dev/null 2>&1; then
    echo "Node.js is required to write $config_path (merge with existing MCP servers)." >&2
    return 1
  fi

  mkdir -p "$(dirname "$config_path")"
  CONFIG_PATH="$config_path" ROOT_KEY="$root_key" USE_CURSOR_BIN="$use_cursor_bin" USE_WINDSURF_ENV="$use_windsurf_env" JOINQUEST_API_KEY="$JOINQUEST_API_KEY" node <<'NODE'
const fs = require('fs')
const path = process.env.CONFIG_PATH
const rootKey = process.env.ROOT_KEY
const useCursorBin = process.env.USE_CURSOR_BIN === 'true'
const useWindsurfEnv = process.env.USE_WINDSURF_ENV === 'true'

const server = {
  type: 'stdio',
  command: 'npx',
  args: useCursorBin
    ? ['--yes', '--package', '@joinquest/mcp-integration', 'joinquest-integration-mcp-cursor']
    : ['-y', '@joinquest/mcp-integration'],
  env: useWindsurfEnv
    ? { JOINQUEST_API_KEY: '${env:JOINQUEST_API_KEY}' }
    : { JOINQUEST_API_KEY: process.env.JOINQUEST_API_KEY },
}

let config = {}
if (fs.existsSync(path)) {
  config = JSON.parse(fs.readFileSync(path, 'utf8'))
}
if (!config[rootKey] || typeof config[rootKey] !== 'object') {
  config[rootKey] = {}
}
config[rootKey]['joinquest-integration'] = server
fs.writeFileSync(path, `${JSON.stringify(config, null, 2)}\n`)
NODE
  echo "Wrote joinquest-integration MCP to $config_path"
}

joinquest_cline_config_path() {
  case "$(uname -s)" in
    Darwin)
      printf '%s/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json' "$HOME"
      ;;
    MINGW* | MSYS* | CYGWIN*)
      printf '%s/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json' "${APPDATA:-$HOME/AppData/Roaming}"
      ;;
    *)
      printf '%s/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json' "${XDG_CONFIG_HOME:-$HOME/.config}"
      ;;
  esac
}

joinquest_windsurf_config_path() {
  printf '%s/.codeium/windsurf/mcp_config.json' "$HOME"
}

joinquest_write_copilot_mcp() {
  joinquest_merge_mcp_config ".vscode/mcp.json" "servers" false false
}

joinquest_write_roo_mcp() {
  joinquest_merge_mcp_config ".roo/mcp.json" "mcpServers" false false
}

joinquest_write_cline_mcp() {
  joinquest_merge_mcp_config "$(joinquest_cline_config_path)" "mcpServers" false false
}

joinquest_write_windsurf_mcp() {
  joinquest_merge_mcp_config "$(joinquest_windsurf_config_path)" "mcpServers" false true
  echo "Export JOINQUEST_API_KEY in your shell before starting Windsurf, e.g.:"
  echo "  export JOINQUEST_API_KEY=$JOINQUEST_API_KEY"
}

joinquest_copilot_mcp_json() {
  cat <<EOF
{
  "servers": {
    "joinquest-integration": $(joinquest_stdio_mcp_server_json false)
  }
}
EOF
}

joinquest_roo_mcp_json() {
  cat <<EOF
{
  "mcpServers": {
    "joinquest-integration": $(joinquest_stdio_mcp_server_json false)
  }
}
EOF
}

joinquest_windsurf_mcp_json() {
  cat <<EOF
{
  "mcpServers": {
    "joinquest-integration": $(joinquest_windsurf_stdio_mcp_server_json)
  }
}
EOF
}

joinquest_cline_mcp_json() {
  cat <<EOF
{
  "mcpServers": {
    "joinquest-integration": $(joinquest_stdio_mcp_server_json false)
  }
}
EOF
}

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
  if ! joinquest_merge_mcp_config "$config_path" "mcpServers" true false; then
    echo "Paste this into $config_path:" >&2
    joinquest_cursor_mcp_json
    return 1
  fi
}

joinquest_claude_desktop_config_path() {
  case "$(uname -s)" in
    Darwin)
      printf '%s/Library/Application Support/Claude/claude_desktop_config.json' "$HOME"
      ;;
    MINGW* | MSYS* | CYGWIN*)
      printf '%s/Claude/claude_desktop_config.json' "${APPDATA:-$HOME/AppData/Roaming}"
      ;;
    *)
      printf '%s/Claude/claude_desktop_config.json' "${XDG_CONFIG_HOME:-$HOME/.config}"
      ;;
  esac
}

joinquest_claude_desktop_mcp_json() {
  cat <<EOF
{
  "mcpServers": {
    "joinquest-integration": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@joinquest/mcp-integration"],
      "env": {
        "JOINQUEST_API_KEY": "$JOINQUEST_API_KEY"
      }
    }
  }
}
EOF
}

joinquest_write_claude_desktop_mcp() {
  local config_path
  config_path="$(joinquest_claude_desktop_config_path)"

  if ! joinquest_merge_mcp_config "$config_path" "mcpServers" false false; then
    echo "Paste this into $config_path:" >&2
    joinquest_claude_desktop_mcp_json
    return 1
  fi
  echo "Fully quit Claude Desktop, reopen it, and start a new conversation."
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
