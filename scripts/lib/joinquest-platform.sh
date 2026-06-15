#!/usr/bin/env bash
# Platform-specific skill / rules files for Copilot, Roo, Cline, Windsurf.
set -euo pipefail

joinquest_skill_source_dir() {
  local script_dir
  script_dir="$(CDPATH= cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  if [ -f "$script_dir/../../.agents/skills/joinquest-integration/SKILL.md" ]; then
    printf '%s/../../.agents/skills/joinquest-integration' "$script_dir"
    return 0
  fi
  if [ -d ".agents/skills/joinquest-integration" ]; then
    printf '%s' ".agents/skills/joinquest-integration"
    return 0
  fi
  return 1
}

joinquest_copy_skill_tree() {
  local src=$1
  local dest=$2
  mkdir -p "$(dirname "$dest")"
  rm -rf "$dest"
  cp -R "$src" "$dest"
}

joinquest_write_mcp_first_instructions() {
  local dest=$1
  cat >"$dest" <<'EOF'
# JoinQuest integration

Integrate a multiplayer game with [JoinQuest](https://joinquest.cc).

**Use MCP first:** at the start of any JoinQuest task, call `joinquest_integration_get_agent_playbook`.
Copilot and other agents may truncate long instruction files — the MCP playbook is the source of truth.

Verify: "List my JoinQuest games using the MCP tools."
Start: "I want to create a game on JoinQuest."

Requires the **joinquest-integration** MCP server (Node.js 20+, `npx @joinquest/mcp-integration`).
EOF
}

joinquest_install_copilot_artifacts() {
  local skill_src
  skill_src="$(joinquest_skill_source_dir)" || {
    echo "Install the agent skill first (missing .agents/skills/joinquest-integration)." >&2
    return 1
  }

  joinquest_copy_skill_tree "$skill_src" ".github/skills/joinquest-integration"
  mkdir -p .github
  joinquest_write_mcp_first_instructions ".github/copilot-instructions.md"
  echo "Installed Copilot skill at .github/skills/joinquest-integration/ and .github/copilot-instructions.md"
}

joinquest_install_roo_artifacts() {
  local skill_src
  skill_src="$(joinquest_skill_source_dir)" || {
    echo "Install the agent skill first (missing .agents/skills/joinquest-integration)." >&2
    return 1
  }

  joinquest_copy_skill_tree "$skill_src" ".roo/rules/joinquest-integration"
  echo "Installed Roo rules at .roo/rules/joinquest-integration/"
}

joinquest_install_cline_artifacts() {
  local skill_src
  skill_src="$(joinquest_skill_source_dir)" || {
    echo "Install the agent skill first (missing .agents/skills/joinquest-integration)." >&2
    return 1
  }

  joinquest_copy_skill_tree "$skill_src" ".cline/rules/joinquest-integration"
  joinquest_write_mcp_first_instructions ".clinerules"
  echo "Installed Cline rules at .cline/rules/joinquest-integration/ and .clinerules"
}

joinquest_install_windsurf_artifacts() {
  local skill_src
  skill_src="$(joinquest_skill_source_dir)" || {
    echo "Install the agent skill first (missing .agents/skills/joinquest-integration)." >&2
    return 1
  }

  joinquest_copy_skill_tree "$skill_src" ".windsurf/rules/joinquest-integration"
  echo "Installed Windsurf rules at .windsurf/rules/joinquest-integration/"
}
