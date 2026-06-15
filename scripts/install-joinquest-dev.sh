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
    curl -fsSL "$GITHUB_RAW/lib/joinquest-platform.sh" -o "$lib_dir/joinquest-platform.sh"
  fi

  # shellcheck source=lib/joinquest-skill.sh
  . "$lib_dir/joinquest-skill.sh"
  # shellcheck source=lib/joinquest-mcp-env.sh
  . "$lib_dir/joinquest-mcp-env.sh"
  # shellcheck source=lib/joinquest-mcp-config.sh
  . "$lib_dir/joinquest-mcp-config.sh"
  # shellcheck source=lib/joinquest-platform.sh
  . "$lib_dir/joinquest-platform.sh"
}

_joinquest_load_libs

MODE="${JOINQUEST_SETUP_MODE:-auto}"
SETUP_CURSOR=false
SETUP_CLAUDE=false
SETUP_CLAUDE_DESKTOP=false
SETUP_COPILOT=false
SETUP_ROO=false
SETUP_WINDSURF=false
SETUP_CLINE=false

usage() {
  cat <<EOF
JoinQuest developer setup — agent skill + MCP (Node.js 20+).

What this script does:
  1. Downloads .agents/skills/joinquest-integration/ from GitHub into this project
  2. With a platform flag + JOINQUEST_API_KEY: merges joinquest-integration MCP config
  3. Adds platform-specific rules/skills where helpful (Copilot, Roo, Cline, Windsurf)
  MCP uses npx @joinquest/mcp-integration and calls joinquest.cc when tools run.

Source (read before running):
  https://github.com/scruffyprodigy/playhub/blob/main/scripts/install-joinquest-dev.sh

Options:
  --cursor          Skill + .cursor/mcp.json
  --claude          Skill + claude mcp add (needs claude CLI)
  --claude-desktop  Skill + Claude Desktop config
  --copilot         Skill + .vscode/mcp.json + Copilot skill/instructions
  --roo             Skill + .roo/mcp.json + Roo rules
  --windsurf        Skill + ~/.codeium/windsurf/mcp_config.json + Windsurf rules
  --cline           Skill + Cline MCP settings + .clinerules
  --all             Skill + Cursor + Claude Code MCP
  --skill-only      Agent skill only
  -h, --help        Show this help

Generate an API key: https://joinquest.cc/developers → Connect an AI assistant

ChatGPT is not supported yet — it requires a hosted HTTPS MCP connector, not local stdio.
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    --cursor) MODE=manual; SETUP_CURSOR=true ;;
    --claude) MODE=manual; SETUP_CLAUDE=true ;;
    --claude-desktop) MODE=manual; SETUP_CLAUDE_DESKTOP=true ;;
    --copilot) MODE=manual; SETUP_COPILOT=true ;;
    --roo) MODE=manual; SETUP_ROO=true ;;
    --windsurf) MODE=manual; SETUP_WINDSURF=true ;;
    --cline) MODE=manual; SETUP_CLINE=true ;;
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
  echo "Next: generate an API key on https://joinquest.cc/developers and re-run with a platform flag."
  exit 0
fi

NEED_KEY=false
if [ "$MODE" = "manual" ]; then
  if $SETUP_CURSOR || $SETUP_CLAUDE || $SETUP_CLAUDE_DESKTOP || $SETUP_COPILOT || $SETUP_ROO || $SETUP_WINDSURF || $SETUP_CLINE; then
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
  echo "  JOINQUEST_API_KEY=lq_dev_... curl -fsSL $GITHUB_RAW/install-joinquest-dev.sh | sh -s -- --copilot"
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

if $SETUP_COPILOT; then
  joinquest_install_copilot_artifacts
  joinquest_write_copilot_mcp
  echo "Restart VS Code, enable Agent mode, and approve joinquest-integration MCP tools."
fi

if $SETUP_ROO; then
  joinquest_install_roo_artifacts
  joinquest_write_roo_mcp
  echo "Reload VS Code or start a new Roo Code chat in this project."
fi

if $SETUP_WINDSURF; then
  joinquest_install_windsurf_artifacts
  joinquest_write_windsurf_mcp
  echo "Fully quit Windsurf (Cmd+Q), reopen, and start a fresh Cascade chat."
fi

if $SETUP_CLINE; then
  joinquest_install_cline_artifacts
  joinquest_write_cline_mcp
  echo "Reload the Cline panel in VS Code and approve joinquest-integration when prompted."
fi

if ! $SETUP_CURSOR && ! $SETUP_CLAUDE && ! $SETUP_CLAUDE_DESKTOP && ! $SETUP_COPILOT && ! $SETUP_ROO && ! $SETUP_WINDSURF && ! $SETUP_CLINE; then
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
  echo ""
  echo "Copilot — add to .vscode/mcp.json:"
  joinquest_copilot_mcp_json
  echo ""
  echo "Roo Code — add to .roo/mcp.json:"
  joinquest_roo_mcp_json
  echo ""
  echo "Windsurf — add to $(joinquest_windsurf_config_path):"
  joinquest_windsurf_mcp_json
  echo ""
  echo "Cline — add to $(joinquest_cline_config_path):"
  joinquest_cline_mcp_json
fi

echo ""
echo "Verify: ask your agent — \"List my JoinQuest games using the MCP tools.\""
echo "Then start: \"I want to create a game on JoinQuest.\""
