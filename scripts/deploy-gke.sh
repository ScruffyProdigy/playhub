#!/bin/bash
# Deploy JoinQuest Lobby + Word Hunt + RPSLS demo to GKE.
#
# Prereqs:
#   - kubectl context pointing at your GKE cluster
#   - k8s/secrets/*.yaml filled in for lobby (joinquest namespace)
#   - ../wordhunt/k8s/secrets/pg-dsn.yaml and ../demo-game-rps/k8s/secrets/pg-dsn.yaml
#
# Usage:
#   ./scripts/deploy-gke.sh              # build, push, deploy all
#   BUILD_PUSH=false ./scripts/deploy-gke.sh   # deploy only (images already on registry)

set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BUILD_PUSH="${BUILD_PUSH:-true}"
WORDHUNT="${WORDHUNT:-$ROOT/../wordhunt}"
DEMO_GAME_RPS="${DEMO_GAME_RPS:-$ROOT/../demo-game-rps}"

echo "=== JoinQuest GKE deploy ==="

if [ "$BUILD_PUSH" = "true" ]; then
  echo "Building and pushing Lobby images..."
  "$ROOT/scripts/build-and-push.sh" --push
fi

echo "Deploying Lobby (namespace joinquest)..."
"$ROOT/scripts/deploy-joinquest.sh"

ensure_game_secret() {
  local game_dir="$1"
  local namespace="$2"
  local secrets="$game_dir/k8s/secrets"
  if [ ! -f "$secrets/pg-dsn.yaml" ]; then
    echo "Creating $secrets/pg-dsn.yaml from example (in-cluster Postgres)..."
    PW="$(openssl rand -hex 12)"
    sed -e "s/namespace: rps-game/namespace: ${namespace}/" \
        -e "s/namespace: rpsls-duel/namespace: ${namespace}/" \
        -e "s/CHANGE_ME/${PW}/g" \
        "$secrets/pg-dsn.example.yaml" > "$secrets/pg-dsn.yaml"
    echo "  (new Postgres password written to pg-dsn.yaml — back up locally)"
  fi
}

if [ -d "$WORDHUNT" ]; then
  ensure_game_secret "$WORDHUNT" "rps-game"
  if [ "$BUILD_PUSH" = "true" ]; then
    echo "Building and pushing Word Hunt images (word-hunt-api:wordhunt, word-hunt-client:wordhunt)..."
    TAG=wordhunt "$WORDHUNT/scripts/build-and-push.sh" --push
  fi
  echo "Deploying Word Hunt (namespace rps-game, host word-hunt-arena.win)..."
  IMAGE_TAG=wordhunt "$WORDHUNT/scripts/deploy-gke.sh"
else
  echo "wordhunt not found at $WORDHUNT — skipping" >&2
fi

if [ -d "$DEMO_GAME_RPS" ]; then
  ensure_game_secret "$DEMO_GAME_RPS" "rpsls-duel"
  if [ "$BUILD_PUSH" = "true" ]; then
    echo "Building and pushing RPSLS images (rps-game-api:rpsls, rps-game-client:rpsls)..."
    TAG=rpsls "$DEMO_GAME_RPS/scripts/build-and-push.sh" --push
  fi
  echo "Deploying RPSLS (namespace rpsls-duel, host rpsls-duel.win)..."
  IMAGE_TAG=rpsls GAME_NAMESPACE=rpsls-duel "$DEMO_GAME_RPS/scripts/deploy-gke.sh"
else
  echo "demo-game-rps not found at $DEMO_GAME_RPS — skipping" >&2
fi

echo ""
echo "=== Verifying game identity on public hosts ==="
"$ROOT/scripts/verify-game-deployments.sh" || {
  echo "WARNING: game verification failed — check IMAGE_TAG and registry repos above." >&2
  exit 1
}

echo ""
echo "=== Done ==="
echo "Lobby:     https://joinquest.cc"
echo "Word Hunt: https://word-hunt-arena.win"
echo "RPSLS:     https://rpsls-duel.win"
echo ""
kubectl get ingress -n joinquest
kubectl get ingress -n rps-game 2>/dev/null || true
kubectl get ingress -n rpsls-duel 2>/dev/null || true
