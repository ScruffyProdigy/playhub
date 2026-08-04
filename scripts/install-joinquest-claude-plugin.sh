#!/usr/bin/env bash
set -euo pipefail
API_KEY="${JOINQUEST_API_KEY:-}"
if [ -z "$API_KEY" ]; then
  echo "Set JOINQUEST_API_KEY first (lq_dev_… from joinquest.cc/developers)." >&2
  exit 1
fi
exec env JOINQUEST_API_KEY="$API_KEY" npx -y joinquest install claude --plugin
