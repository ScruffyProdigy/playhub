#!/bin/bash
# Launch JoinQuest integration MCP under Cursor.
# Cursor injects ELECTRON_RUN_AS_NODE and a bundled "node" on PATH, which breaks stdio MCP.
set -eu

unset ELECTRON_RUN_AS_NODE
export PATH="/usr/local/bin:/opt/homebrew/bin:/usr/bin:/bin:${PATH:-}"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
NODE="${JOINQUEST_NODE:-/usr/local/bin/node}"

exec "$NODE" "$ROOT/mcp/joinquest-integration/src/index.js"
