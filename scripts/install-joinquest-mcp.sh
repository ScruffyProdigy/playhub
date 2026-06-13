#!/usr/bin/env sh
set -eu

# One-time install of the JoinQuest integration MCP CLI.
# Requires Node.js 20+ and npm.

REPO_ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
PKG_DIR="$REPO_ROOT/mcp/joinquest-integration"

if [ ! -f "$PKG_DIR/package.json" ]; then
  echo "Could not find mcp/joinquest-integration in $REPO_ROOT" >&2
  exit 1
fi

cd "$PKG_DIR"
npm install --omit=dev
npm link
echo "Installed joinquest-integration-mcp globally. Re-run your MCP client (Cursor or Claude)."
