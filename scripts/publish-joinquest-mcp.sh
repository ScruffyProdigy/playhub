#!/usr/bin/env bash
# Publish @joinquest/mcp-integration to the public npm registry.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PKG_DIR="$ROOT/mcp/joinquest-integration"

if ! command -v npm >/dev/null 2>&1; then
  echo "npm not found" >&2
  exit 1
fi

if ! npm whoami >/dev/null 2>&1; then
  echo "Not logged in to npm. Run: npm login" >&2
  echo "You must have publish access to the @joinquest scope." >&2
  exit 1
fi

# npm publish requires account 2FA (auth-and-writes) + OTP, OR a granular token with bypass 2FA.
TFA_STATUS=$(npm profile get tfa 2>/dev/null | awk '{print $NF}' || true)
if [[ "${TFA_STATUS:-}" == "disabled" && -z "${NPM_TOKEN:-}" ]]; then
  echo "npm account 2FA is disabled ($(npm whoami))." >&2
  echo "OTP alone will not work. Choose one:" >&2
  echo "  1. Enable 2FA: https://www.npmjs.com/settings/$(npm whoami)/tfa" >&2
  echo "     (mode: Authorization and publishing), then:" >&2
  echo "     NPM_OTP=123456 ./scripts/publish-joinquest-mcp.sh" >&2
  echo "  2. Granular access token (Publish on @joinquest/*, bypass 2FA enabled):" >&2
  echo "     NPM_TOKEN=npm_... ./scripts/publish-joinquest-mcp.sh" >&2
  echo >&2
  echo "Also create the @joinquest org first if it does not exist:" >&2
  echo "  https://www.npmjs.com/org/create" >&2
  exit 1
fi

if ! npm org ls joinquest >/dev/null 2>&1; then
  echo "The @joinquest npm organization does not exist or you are not a member." >&2
  echo "Create it at https://www.npmjs.com/org/create (name: joinquest), then retry." >&2
  exit 1
fi

cd "$PKG_DIR"
echo "Running tests..."
npm test

NAME=$(node --input-type=module -e "import p from './package.json' with { type: 'json' }; console.log(p.name + '@' + p.version)")
echo "Publishing $NAME..."

if [[ -n "${NPM_TOKEN:-}" ]]; then
  npm publish --access public --userconfig <(printf '//registry.npmjs.org/:_authToken=%s\n' "$NPM_TOKEN")
else
  PUBLISH_ARGS=(--access public)
  if [[ -n "${NPM_OTP:-}" ]]; then
    PUBLISH_ARGS+=(--otp "$NPM_OTP")
  else
    echo "Tip: pass NPM_OTP=123456 (authenticator code) for interactive publish." >&2
  fi
  if ! npm publish "${PUBLISH_ARGS[@]}"; then
    echo >&2
    echo "Publish failed." >&2
    echo "See https://docs.npmjs.com/requiring-2fa-for-package-publishing-and-settings-modification" >&2
    exit 1
  fi
fi

echo "Done. Test: npx -y @joinquest/mcp-integration (stdio MCP; Ctrl+C to exit)"
