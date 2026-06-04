#!/bin/bash
# Deploy JoinQuest Lobby + demo-game-rps to GKE (joinquest.cc + rpsls-duel.win).
#
# Prereqs:
#   - kubectl context pointing at your GKE cluster
#   - k8s/secrets/*.yaml filled in for lobby (joinquest namespace)
#   - ../demo-game-rps/k8s/secrets/pg-dsn.yaml for the game DB
#   - Recommended: k8s/secrets/lobby-auth-peppers.yaml (MAGIC_LINK_PEPPER, LOBBY_GAME_TOKEN_PEPPER)
#   - Optional legacy: lobby-game-service.yaml (global LOBBY_GAME_SERVICE_TOKEN for dev)
#
# Usage:
#   ./scripts/deploy-gke.sh              # build, push, deploy both
#   BUILD_PUSH=false ./scripts/deploy-gke.sh   # deploy only (images already on registry)
#
# RPS images default to docker.io/scruffyprodigy (see demo-game-rps k8s/base). Override:
#   RPS_REGISTRY=ghcr.io RPS_IMAGE_OWNER=playhub ./scripts/deploy-gke.sh

set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BUILD_PUSH="${BUILD_PUSH:-true}"
DEMO_GAME_RPS="${DEMO_GAME_RPS:-$ROOT/../demo-game-rps}"

echo "=== JoinQuest GKE deploy ==="

if [ "$BUILD_PUSH" = "true" ]; then
  echo "Building and pushing Lobby images..."
  "$ROOT/scripts/build-and-push.sh" --push
fi

echo "Deploying Lobby (namespace joinquest)..."
"$ROOT/scripts/deploy-joinquest.sh"

if [ ! -d "$DEMO_GAME_RPS" ]; then
  echo "demo-game-rps not found at $DEMO_GAME_RPS — skipping game deploy" >&2
  exit 1
fi

RPS_SECRETS="$DEMO_GAME_RPS/k8s/secrets"
if [ ! -f "$RPS_SECRETS/pg-dsn.yaml" ]; then
  echo "Creating $RPS_SECRETS/pg-dsn.yaml from example (in-cluster Postgres)..."
  PW="$(openssl rand -hex 12)"
  sed -e "s/namespace: rps-game/namespace: rps-game/" \
      -e "s/CHANGE_ME/${PW}/g" \
      "$RPS_SECRETS/pg-dsn.example.yaml" > "$RPS_SECRETS/pg-dsn.yaml"
  echo "  (new Postgres password written to pg-dsn.yaml — back up locally)"
fi

if [ "$BUILD_PUSH" = "true" ]; then
  echo "Building and pushing RPS game images (Docker Hub)..."
  REGISTRY="${RPS_REGISTRY:-docker.io}" IMAGE_OWNER="${RPS_IMAGE_OWNER:-scruffyprodigy}" \
    "$DEMO_GAME_RPS/scripts/build-and-push.sh" --push
fi

echo "Deploying demo-game-rps (namespace rps-game, host rpsls-duel.win)..."
"$DEMO_GAME_RPS/scripts/deploy-gke.sh"

echo ""
echo "=== Done ==="
echo "Lobby:  https://joinquest.cc"
echo "Game:   https://rpsls-duel.win"
echo ""
kubectl get ingress -n joinquest
kubectl get ingress -n rps-game
echo ""
echo "DNS: joinquest.cc and rpsls-duel.win -> ingress ADDRESS (kubectl get ingress -A)"
