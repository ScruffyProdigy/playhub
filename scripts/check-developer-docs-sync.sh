#!/usr/bin/env bash
# Fail if canonical developer docs are not synced to embeds and skill copies.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "Checking developer doc sync..."

"$ROOT/scripts/sync-developer-docs.sh" >/dev/null

PATHS=(
  backend/internal/developer/integration_guide.md
  backend/internal/developer/agent_playbook.md
  .agents/skills/joinquest-integration/playbook.md
  plugins/joinquest/skills/joinquest-integration/SKILL.md
  plugins/joinquest/skills/joinquest-integration/mcp-setup.md
  plugins/joinquest/skills/joinquest-integration/playbook.md
)

if git diff --quiet -- "${PATHS[@]}"; then
  echo "✅ Developer doc copies are in sync."
  exit 0
fi

echo "❌ Developer doc copies are out of sync." >&2
echo "Run ./scripts/sync-developer-docs.sh and commit the updated copies." >&2
echo "" >&2
git diff -- "${PATHS[@]}" >&2
exit 1
