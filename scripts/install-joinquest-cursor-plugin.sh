#!/usr/bin/env bash
# Install JoinQuest Cursor plugin (skill + MCP) into ~/.cursor/plugins/local/joinquest
set -euo pipefail

GITHUB_ARCHIVE="${JOINQUEST_GITHUB_ARCHIVE:-https://github.com/scruffyprodigy/playhub/archive/refs/heads/main.tar.gz}"
PLUGIN_DIR="${HOME}/.cursor/plugins/local/joinquest"
API_KEY="${JOINQUEST_API_KEY:-}"

if [ -z "$API_KEY" ]; then
  echo "Set JOINQUEST_API_KEY first (lq_dev_… from joinquest.cc/developers)." >&2
  exit 1
fi

mkdir -p "${HOME}/.cursor/plugins/local"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading JoinQuest plugin..."
curl -fsSL "$GITHUB_ARCHIVE" | tar -xz -C "$tmp" --strip-components=1

rm -rf "$PLUGIN_DIR"
cp -R "$tmp/plugins/joinquest" "$PLUGIN_DIR"

mcp_json="$PLUGIN_DIR/mcp.json"
if [ ! -f "$mcp_json" ]; then
  echo "Plugin mcp.json not found at $mcp_json" >&2
  exit 1
fi

if command -v perl >/dev/null 2>&1; then
  perl -i -pe 's/<paste-api-key-from-joinquest-dashboard>/$ENV{JOINQUEST_API_KEY}/' "$mcp_json"
elif sed --version 2>/dev/null | grep -q GNU; then
  sed -i "s/<paste-api-key-from-joinquest-dashboard>/${API_KEY}/" "$mcp_json"
else
  sed -i '' "s/<paste-api-key-from-joinquest-dashboard>/${API_KEY}/" "$mcp_json"
fi

echo "Installed JoinQuest plugin to $PLUGIN_DIR"
echo "Quit Cursor completely (Cmd+Q), reopen your game project, then start a fresh Agent chat."
