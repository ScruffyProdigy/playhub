#!/usr/bin/env bash
# Verify production game hosts serve the expected API identity.
# Run after deploy-gke.sh or any manual game rollout.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

game_from_status() {
  local url=$1
  local json
  json=$(curl -sfS "$url" 2>/dev/null) || return 1
  if command -v python3 >/dev/null 2>&1; then
    GAME_STATUS_JSON="$json" python3 - <<'PY'
import json, os
print(json.loads(os.environ["GAME_STATUS_JSON"])["game"])
PY
    return
  fi
  echo "$json" | sed -n 's/.*"game":"\([^"]*\)".*/\1/p' | head -1
}

verify_host() {
  local host=$1 expected=$2
  local actual
  actual=$(game_from_status "https://${host}/api/v1/status") || {
    echo "FAIL: could not reach https://${host}/api/v1/status" >&2
    return 1
  }
  if [ "$actual" != "$expected" ]; then
    echo "FAIL: https://${host} reports game=${actual}, expected ${expected}" >&2
    return 1
  fi
  echo "OK  https://${host} → ${actual}"
}

echo "=== Verifying production game deployments ==="
failed=0
verify_host word-hunt-arena.win word-hunt || failed=1
verify_host rpsls-duel.win rock-paper-scissors-lizard-robot || failed=1

if [ "$failed" -ne 0 ]; then
  echo "" >&2
  echo "Game identity check failed. Common causes:" >&2
  echo "  - Wrong IMAGE_TAG (never use :latest for demo games)" >&2
  echo "  - Word Hunt must use word-hunt-api/word-hunt-client:wordhunt" >&2
  echo "  - RPS must use rps-game-api/rps-game-client:rpsls" >&2
  exit 1
fi

echo "All game hosts verified."
