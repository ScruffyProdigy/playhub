#!/usr/bin/env bash
# Install JoinQuest Claude Code plugin (skill + MCP) into ~/.claude/skills/joinquest-integration
set -euo pipefail

GITHUB_ARCHIVE="${JOINQUEST_GITHUB_ARCHIVE:-https://github.com/scruffyprodigy/playhub/archive/refs/heads/main.tar.gz}"
PLUGIN_DIR="${HOME}/.claude/skills/joinquest-integration"
API_KEY="${JOINQUEST_API_KEY:-}"

if [ -z "$API_KEY" ]; then
  echo "Set JOINQUEST_API_KEY first (lq_dev_… from joinquest.cc/developers)." >&2
  exit 1
fi

mkdir -p "${HOME}/.claude/skills"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading JoinQuest plugin..."
curl -fsSL "$GITHUB_ARCHIVE" | tar -xz -C "$tmp" --strip-components=1

rm -rf "$PLUGIN_DIR"
cp -R "$tmp/plugins/joinquest" "$PLUGIN_DIR"

mcp_json="$PLUGIN_DIR/.mcp.json"
if [ ! -f "$mcp_json" ]; then
  echo "Plugin .mcp.json not found at $mcp_json" >&2
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
echo "Start a new Claude Code session in your game repo (or run /reload-plugins)."
echo "Try: /joinquest-integration:joinquest-integration or ask to integrate with JoinQuest."
