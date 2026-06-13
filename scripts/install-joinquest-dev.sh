#!/usr/bin/env bash
# Install JoinQuest agent skill + MCP setup for game developers.
#
# From your game repo root:
#   JOINQUEST_API_KEY=lq_dev_... curl -fsSL .../install-joinquest-dev.sh | sh -s -- --cursor
set -euo pipefail

GITHUB_RAW="${JOINQUEST_GITHUB_RAW:-https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts}"

_joinquest_load_libs() {
  local lib_dir=""
  local src="${BASH_SOURCE[0]:-$0}"

  if [ "$src" != "sh" ] && [ "$src" != "bash" ] && [ -f "$(dirname "$src")/lib/joinquest-skill.sh" ]; then
    lib_dir="$(CDPATH= cd -- "$(dirname "$src")" && pwd)/lib"
  else
    lib_dir="$(mktemp -d)"
    curl -fsSL "$GITHUB_RAW/lib/joinquest-skill.sh" -o "$lib_dir/joinquest-skill.sh"
    curl -fsSL "$GITHUB_RAW/lib/joinquest-mcp-env.sh" -o "$lib_dir/joinquest-mcp-env.sh"
    curl -fsSL "$GITHUB_RAW/lib/joinquest-mcp-config.sh" -o "$lib_dir/joinquest-mcp-config.sh"
  fi

  # shellcheck source=lib/joinquest-skill.sh
  . "$lib_dir/joinquest-skill.sh"
  # shellcheck source=lib/joinquest-mcp-env.sh
  . "$lib_dir/joinquest-mcp-env.sh"
  # shellcheck source=lib/joinquest-mcp-config.sh
  . "$lib_dir/joinquest-mcp-config.sh"
}

_joinquest_load_libs

MODE="${JOINQUEST_SETUP_MODE:-auto}"
SETUP_CURSOR=false
SETUP_CLAUDE=false
SETUP_CLAUDE_DESKTOP=false

usage() {
  cat <<EOF
JoinQuest developer setup — agent skill + MCP (Node.js 20+).

What this script does:
  1. Downloads .agents/skills/joinquest-integration/ from GitHub into this project
     (agent instructions — markdown only, no code execution from the skill itself).
  2. With --cursor + JOINQUEST_API_KEY: merges joinquest-integration into .cursor/mcp.json
  3. With --claude + JOINQUEST_API_KEY: runs "claude mcp add" for this repo
  4. With --claude-desktop + JOINQUEST_API_KEY: merges into Claude Desktop config
  MCP uses npx @joinquest/mcp-integration and calls joinquest.cc when tools run.

Source (read before running):
  https://github.com/scruffyprodigy/playhub/blob/main/scripts/install-joinquest-dev.sh

Safer install (download, review, then run):
  curl -fsSL $GITHUB_RAW/install-joinquest-dev.sh -o install-joinquest-dev.sh
  less install-joinquest-dev.sh
  JOINQUEST_API_KEY=lq_dev_... bash install-joinquest-dev.sh --cursor

Quick install (pipes into shell):
  JOINQUEST_API_KEY=lq_dev_... curl -fsSL $GITHUB_RAW/install-joinquest-dev.sh | sh -s -- --cursor

Options:
  --cursor          Skill + .cursor/mcp.json (needs JOINQUEST_API_KEY)
  --claude          Skill + claude mcp add (needs JOINQUEST_API_KEY + claude CLI)
  --claude-desktop  Skill + Claude Desktop claude_desktop_config.json (needs JOINQUEST_API_KEY)
  --all             Skill + Cursor + Claude Code MCP
  --skill-only      Agent skill only
  -h, --help    Show this help

Generate an API key: https://joinquest.cc/developers → Connect an AI assistant
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    --cursor) MODE=manual; SETUP_CURSOR=true ;;
    --claude) MODE=manual; SETUP_CLAUDE=true ;;
    --claude-desktop) MODE=manual; SETUP_CLAUDE_DESKTOP=true ;;
    --all) MODE=manual; SETUP_CURSOR=true; SETUP_CLAUDE=true ;;
    --skill-only) MODE=skill-only ;;
    -h | --help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 1 ;;
  esac
  shift
done

joinquest_install_skill

if [ "$MODE" = "skill-only" ]; then
  echo ""
  echo "Next: generate an API key on https://joinquest.cc/developers and re-run with --cursor, --claude, or --claude-desktop."
  exit 0
fi

NEED_KEY=false
if [ "$MODE" = "manual" ]; then
  if $SETUP_CURSOR || $SETUP_CLAUDE || $SETUP_CLAUDE_DESKTOP; then
    NEED_KEY=true
  fi
else
  if [ -n "${JOINQUEST_API_KEY:-}" ]; then
    SETUP_CURSOR=true
  fi
fi

if $NEED_KEY || { [ "$MODE" = "auto" ] && [ -n "${JOINQUEST_API_KEY:-}" ]; }; then
  joinquest_load_mcp_env
fi

if [ -z "${JOINQUEST_API_KEY:-}" ]; then
  echo ""
  echo "Skill installed. For MCP, set JOINQUEST_API_KEY and re-run, e.g.:"
  echo "  JOINQUEST_API_KEY=lq_dev_... curl -fsSL $GITHUB_RAW/install-joinquest-dev.sh | sh -s -- --cursor"
  exit 0
fi

if $SETUP_CURSOR; then
  joinquest_write_cursor_mcp ".cursor/mcp.json"
  echo "Fully quit Cursor (Cmd+Q), reopen this project, and start a new Agent chat."
fi

if $SETUP_CLAUDE; then
  joinquest_setup_claude_mcp || true
fi

if $SETUP_CLAUDE_DESKTOP; then
  joinquest_write_claude_desktop_mcp
fi

if ! $SETUP_CURSOR && ! $SETUP_CLAUDE && ! $SETUP_CLAUDE_DESKTOP; then
  echo ""
  echo "MCP manual setup (API key included):"
  echo ""
  echo "Cursor — add to .cursor/mcp.json:"
  joinquest_cursor_mcp_json
  echo ""
  echo "Claude Code — run:"
  joinquest_print_claude_command
  echo ""
  echo "Claude Desktop — add to $(joinquest_claude_desktop_config_path):"
  joinquest_claude_desktop_mcp_json
fi

echo ""
echo "Verify: ask your agent — \"List my JoinQuest games using the MCP tools.\""
echo "Then start: \"I want to create a game on JoinQuest.\""
