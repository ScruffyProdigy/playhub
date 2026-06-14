#!/usr/bin/env bash
# Sync canonical developer docs from docs/ to backend embeds and agent skill copies.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

embed_body() {
  local src="$1"
  local canonical="$2"
  head -1 "$src"
  echo ""
  printf '> **Embedded copy:** Do not edit here. Canonical source: `%s`. Run `./scripts/sync-developer-docs.sh` after editing.\n' "$canonical"
  echo ""
  awk 'NR==1 { next }
       NR==2 && /^$/ { next }
       /^> \*\*/ && !body { next }
       { body=1; print }' "$src"
}

GUIDE_SRC="$ROOT/docs/developer-integration-guide.md"
PLAYBOOK_SRC="$ROOT/docs/developer-agent-playbook.md"

embed_body "$GUIDE_SRC" "docs/developer-integration-guide.md" \
  > "$ROOT/backend/internal/developer/integration_guide.md"
echo "Synced docs/developer-integration-guide.md -> backend/internal/developer/integration_guide.md"

embed_body "$PLAYBOOK_SRC" "docs/developer-agent-playbook.md" \
  > "$ROOT/backend/internal/developer/agent_playbook.md"
echo "Synced docs/developer-agent-playbook.md -> backend/internal/developer/agent_playbook.md"

embed_body "$PLAYBOOK_SRC" "docs/developer-agent-playbook.md" \
  > "$ROOT/.agents/skills/joinquest-integration/playbook.md"
echo "Synced docs/developer-agent-playbook.md -> .agents/skills/joinquest-integration/playbook.md"

PLUGIN_SKILL="$ROOT/plugins/joinquest/skills/joinquest-integration"
mkdir -p "$PLUGIN_SKILL"
cp "$ROOT/.agents/skills/joinquest-integration/SKILL.md" "$PLUGIN_SKILL/"
cp "$ROOT/.agents/skills/joinquest-integration/mcp-setup.md" "$PLUGIN_SKILL/"
cp "$ROOT/.agents/skills/joinquest-integration/playbook.md" "$PLUGIN_SKILL/"
echo "Synced agent skill -> plugins/joinquest/skills/joinquest-integration/"

echo "Done."
