#!/bin/bash
# GKE Autopilot/Warden blocks coordination.k8s.io/leases in kube-system.
# cert-manager must elect its leader in the cert-manager namespace with explicit leases RBAC.
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if ! kubectl get deployment cert-manager -n cert-manager >/dev/null 2>&1; then
  echo "cert-manager deployment not found in cert-manager namespace; skip patch."
  exit 0
fi

echo "Applying cert-manager GKE leader-election RBAC..."
kubectl apply -f "$ROOT/k8s/env/cert-manager-gke-rbac.yaml"

NEEDS_ROLLOUT=false

if ! kubectl get deployment cert-manager -n cert-manager -o json \
  | grep -q 'leader-election-namespace=cert-manager'; then
  echo "Patching cert-manager: --leader-election-namespace=cert-manager"
  kubectl patch deployment cert-manager -n cert-manager --type=json -p='[
    {"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--leader-election-namespace=cert-manager"}
  ]'
  NEEDS_ROLLOUT=true
else
  echo "cert-manager deployment already has --leader-election-namespace=cert-manager"
fi

if kubectl get deployment cert-manager-cainjector -n cert-manager >/dev/null 2>&1; then
  if ! kubectl get deployment cert-manager-cainjector -n cert-manager -o json \
    | grep -q 'leader-election-namespace=cert-manager'; then
    echo "Patching cert-manager-cainjector: --leader-election-namespace=cert-manager"
    kubectl patch deployment cert-manager-cainjector -n cert-manager --type=json -p='[
      {"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--leader-election-namespace=cert-manager"}
    ]' 2>/dev/null || true
    NEEDS_ROLLOUT=true
  fi
fi

# Restart so controller picks up RBAC + args (safe if already running).
kubectl rollout restart deployment/cert-manager -n cert-manager
NEEDS_ROLLOUT=true

if [ "$NEEDS_ROLLOUT" = true ]; then
  kubectl rollout status deployment/cert-manager -n cert-manager --timeout=120s
fi

echo "Done. Expect: successfully acquired lease cert-manager/cert-manager-controller"
echo "  kubectl logs -n cert-manager deployment/cert-manager --tail=15"
