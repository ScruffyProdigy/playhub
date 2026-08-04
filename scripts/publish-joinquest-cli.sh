#!/usr/bin/env bash
# Publish joinquest and create-joinquest npm packages.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if ! command -v npm >/dev/null 2>&1; then
  echo "npm not found" >&2
  exit 1
fi

if ! npm whoami >/dev/null 2>&1; then
  echo "Not logged in to npm. Run: npm login" >&2
  exit 1
fi

publish_pkg() {
  local dir=$1
  echo "==> Testing $dir"
  (cd "$dir" && npm test)
  echo "==> Publishing $dir"
  (cd "$dir" && npm publish --access public)
}

(cd "$ROOT/packages/joinquest" && npm run sync-assets)
publish_pkg "$ROOT/packages/joinquest"
publish_pkg "$ROOT/packages/create-joinquest"

echo "Done. Try: npx joinquest install cursor --dry-run"
