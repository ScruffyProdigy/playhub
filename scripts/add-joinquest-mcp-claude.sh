#!/usr/bin/env bash
# Register JoinQuest integration MCP with Claude Code (project scope).
# Requires: claude CLI, Node.js 20+.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=lib/joinquest-mcp-env.sh
. "$ROOT/scripts/lib/joinquest-mcp-env.sh"

if ! command -v claude >/dev/null 2>&1; then
  echo "claude CLI not found. Install Claude Code: https://code.claude.com/docs/en/setup" >&2
  exit 1
fi

if ! command -v npx >/dev/null 2>&1; then
  echo "npx not found. Install Node.js 20+." >&2
  exit 1
fi

joinquest_load_mcp_env

claude mcp add --scope project --transport stdio \
  --env "JOINQUEST_API_KEY=$JOINQUEST_API_KEY" \
  joinquest-integration -- npx -y @joinquest/mcp-integration

echo "Added joinquest-integration to .mcp.json (project scope)."
echo "Start a new claude session in this directory and approve the server when prompted."
