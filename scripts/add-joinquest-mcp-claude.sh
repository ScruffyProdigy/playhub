#!/usr/bin/env bash
# Deprecated: use install-joinquest-dev.sh --claude
ROOT="$(CDPATH= cd -- "$(dirname "$0")" && pwd)"
exec "$ROOT/install-joinquest-dev.sh" --claude "$@"
