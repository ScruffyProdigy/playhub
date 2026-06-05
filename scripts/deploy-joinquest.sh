#!/bin/bash
# Deploy Lobby to namespace joinquest (GKE demo). Rewrites playhub -> joinquest in k8s/base.
# Secrets in k8s/secrets/*.yaml must already use metadata.namespace: joinquest.
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
# shellcheck source=lib/lobby-smtp.sh
. "$ROOT/scripts/lib/lobby-smtp.sh"
# shellcheck source=lib/lobby-game-service.sh
. "$ROOT/scripts/lib/lobby-game-service.sh"
# shellcheck source=lib/lobby-auth-peppers.sh
. "$ROOT/scripts/lib/lobby-auth-peppers.sh"

NAMESPACE="joinquest"
CONTEXT="${KUBE_CONTEXT:-}"
LOBBY_PUBLIC_URL="${LOBBY_PUBLIC_URL:-https://joinquest.cc}"
SMTP_FROM="${SMTP_FROM:-noreply@joinquest.cc}"
SMTP_FROM_NAME="${SMTP_FROM_NAME:-JoinQuest}"
INSTALL_CERT_MANAGER="${INSTALL_CERT_MANAGER:-false}"

echo "Deploying Lobby to namespace: $NAMESPACE"

if ! command -v kubectl &> /dev/null; then
  echo "kubectl not found" >&2
  exit 1
fi

if [ -n "$CONTEXT" ]; then
  echo "Using context: $CONTEXT"
  kubectl config use-context "$CONTEXT"
fi

apply_playhub_as_joinquest() {
  sed 's/namespace: playhub/namespace: joinquest/g' "$1" | kubectl apply -f -
}

if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
  kubectl create namespace "$NAMESPACE"
fi

if [ "$INSTALL_CERT_MANAGER" = "true" ]; then
  echo "Installing cert-manager..."
  kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.5/cert-manager.yaml
  kubectl wait --for=condition=available --timeout=300s deployment/cert-manager -n cert-manager
  kubectl wait --for=condition=available --timeout=300s deployment/cert-manager-webhook -n cert-manager
fi

# GKE Autopilot: controller cannot use kube-system for leader election (Warden denies leases).
"$ROOT/scripts/patch-cert-manager-gke.sh"

echo "Applying Let's Encrypt ClusterIssuer (cluster-wide)..."
kubectl apply -f k8s/env/cluster-issuer-letsencrypt-prod.yaml

echo "Applying secrets..."
kubectl apply -f k8s/secrets/pg-auth.yaml
kubectl apply -f k8s/secrets/pg-dsn.yaml
kubectl apply -f k8s/secrets/jwks-secret.yaml

echo "Applying base manifests (namespace joinquest)..."
apply_playhub_as_joinquest k8s/base/postgres.yaml
apply_playhub_as_joinquest k8s/base/backend.yaml
apply_playhub_as_joinquest k8s/base/frontend.yaml

echo "Applying joinquest TLS certificate + ingress..."
kubectl apply -f k8s/env/joinquest-certificate.yaml
kubectl apply -f k8s/env/joinquest-ingress.yaml

echo "Applying joinquest frontend config..."
kubectl apply -f k8s/env/joinquest.yaml

echo "Configuring backend for browser auth at $LOBBY_PUBLIC_URL ..."
kubectl set env deployment/lobby-backend -n "$NAMESPACE" \
  APP_ENV=production \
  CORS_ALLOWED_ORIGINS="$LOBBY_PUBLIC_URL" \
  MAGIC_LINK_BASE_URL="${LOBBY_PUBLIC_URL}/auth/complete?token=" \
  LOBBY_ISSUER_URL="${LOBBY_ISSUER_URL:-$LOBBY_PUBLIC_URL}" \
  SESSION_COOKIE_SECURE=true \
  LOBBY_ADMIN_EMAILS="${LOBBY_ADMIN_EMAILS:-ryan.c.kohler@gmail.com}" \
  --containers=backend
apply_lobby_smtp_secret
apply_lobby_auth_peppers_secret
apply_lobby_game_service_secret

echo "Waiting for Postgres..."
kubectl wait --for=condition=ready --timeout=300s pod -l app=pg -n "$NAMESPACE"

echo "Running database migrations..."
kubectl delete job playhub-db-migrate -n "$NAMESPACE" --ignore-not-found
sed 's/namespace: playhub/namespace: joinquest/g' k8s/jobs/migration.yaml | kubectl apply -f -
kubectl wait --for=condition=complete --timeout=120s job/playhub-db-migrate -n "$NAMESPACE"
kubectl logs job/playhub-db-migrate -n "$NAMESPACE"

echo "Applying stale session cleanup CronJob..."
sed 's/namespace: playhub/namespace: joinquest/g' k8s/jobs/stale-session-cleanup.yaml | kubectl apply -f -

echo "Patching game catalog handoff URLs for production..."
run_patch_game_handoff_urls_job "$NAMESPACE"

echo "Waiting for app deployments..."
kubectl wait --for=condition=available --timeout=300s deployment/lobby-backend -n "$NAMESPACE"
kubectl wait --for=condition=available --timeout=300s deployment/lobby-frontend -n "$NAMESPACE"

kubectl get pods,svc,ingress -n "$NAMESPACE"
echo "TLS: kubectl get certificate -n $NAMESPACE"
kubectl get certificate -n "$NAMESPACE" 2>/dev/null || true
echo "Done. Public URL: $LOBBY_PUBLIC_URL"
echo "DNS: point joinquest.cc (A or proxied CNAME) at the ingress ADDRESS below."
echo "Frontend API URL: $(kubectl get configmap lobby-frontend-config -n "$NAMESPACE" -o jsonpath='{.data.REACT_APP_API_BASE_URL}')"
