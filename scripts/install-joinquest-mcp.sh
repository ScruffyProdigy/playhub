#!/usr/bin/env sh
set -eu

# Optional global install of JoinQuest MCP CLIs (offline / pinned versions).
# Most developers should use npx instead — no install required:
#   npx -y @joinquest/mcp-integration
#   npx --yes --package @joinquest/mcp-integration joinquest-integration-mcp-cursor  # Cursor

REPO_ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
PKG_DIR="$REPO_ROOT/mcp/joinquest-integration"

if [ ! -f "$PKG_DIR/package.json" ]; then
  echo "Could not find mcp/joinquest-integration in $REPO_ROOT" >&2
  exit 1
fi

cd "$PKG_DIR"
npm install --omit=dev
npm link
echo "Linked joinquest-integration-mcp and joinquest-integration-mcp-cursor globally."
echo ""
echo "Recommended (no global install): use npx in MCP config:"
echo "  Cursor:  npx --yes --package @joinquest/mcp-integration joinquest-integration-mcp-cursor"
echo "  Claude:  npx -y @joinquest/mcp-integration"
echo ""
echo "Claude Code: bash scripts/install-joinquest-dev.sh --claude"
echo "See docs/developer-ai-setup-roadmap.md"
