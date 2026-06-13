#!/usr/bin/env bash
# Shared defaults for JoinQuest MCP setup scripts.
set -euo pipefail

joinquest_load_mcp_env() {
  if [ -z "${JOINQUEST_API_KEY:-}" ]; then
    if [ -t 0 ]; then
      read -r -p "JoinQuest API key (lq_dev_… from developer dashboard): " JOINQUEST_API_KEY
    else
      echo "JOINQUEST_API_KEY is required (generate one on the developer dashboard)." >&2
      exit 1
    fi
  fi

  if [ -z "$JOINQUEST_API_KEY" ]; then
    echo "JOINQUEST_API_KEY cannot be empty." >&2
    exit 1
  fi

  export JOINQUEST_API_KEY
}
