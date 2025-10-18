#!/usr/bin/env bash
set -e

# List of environment variables you want to expose
VARS=(
  REACT_APP_API_BASE_URL
  REACT_APP_ENV
)

OUT_FILE="/usr/share/nginx/html/env.js"

echo "Creating $OUT_FILE ..."
echo "window.env = window.env || {};" > "$OUT_FILE"

for VAR in "${VARS[@]}"; do
  VALUE="${!VAR}"
  if [ -n "$VALUE" ]; then
    SAFE_VALUE=$(printf '%s' "$VALUE" | sed 's/"/\\"/g')
    echo "window.env.$VAR = \"$SAFE_VALUE\";" >> "$OUT_FILE"
    echo "  Added $VAR"
  fi
done

echo "✅ Generated $OUT_FILE"
