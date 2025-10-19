#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

echo "🚀 Deploying PlayHub to Staging..."

# Configuration
NAMESPACE="playhub-staging"
CONTEXT="staging-cluster"  # Change this to your staging cluster context

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed or not in PATH"
    exit 1
fi

# Set context if specified
if [ ! -z "$CONTEXT" ]; then
    print_status "Setting Kubernetes context to: $CONTEXT"
    kubectl config use-context $CONTEXT
fi

# Create namespace
print_step "Creating namespace..."
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Apply base configurations (database, backend, frontend, ingress)
print_step "Applying base configurations..."
kubectl apply -f k8s/base/postgres.yaml -n $NAMESPACE
kubectl apply -f k8s/base/backend.yaml -n $NAMESPACE
kubectl apply -f k8s/base/frontend.yaml -n $NAMESPACE
kubectl apply -f k8s/base/ingress.yaml -n $NAMESPACE

# Apply database secrets and initialization
print_step "Applying database secrets and initialization..."
kubectl apply -f k8s/secrets/pg-auth.yaml -n $NAMESPACE
kubectl apply -f k8s/secrets/pg-dsn.yaml -n $NAMESPACE
kubectl apply -f k8s/init/db-init.yaml -n $NAMESPACE

# Apply staging environment configuration
print_step "Applying staging environment configuration..."
kubectl apply -f k8s/env/staging.yaml -n $NAMESPACE

# Wait for database to be ready
print_step "Waiting for database to be ready..."
kubectl wait --for=condition=ready --timeout=300s pod -l app=pg -n $NAMESPACE

# Wait for deployments to be ready
print_step "Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/lobby-backend -n $NAMESPACE
kubectl wait --for=condition=available --timeout=300s deployment/lobby-frontend -n $NAMESPACE

# Show deployment status
print_step "Deployment status:"
kubectl get pods -n $NAMESPACE
kubectl get services -n $NAMESPACE
kubectl get ingress -n $NAMESPACE

print_status "Staging deployment completed successfully!"
echo ""
print_status "Access your application:"
echo "  Frontend: https://staging.playhub.com"
echo "  Backend API: https://api-staging.playhub.com"
echo "  GraphQL: https://api-staging.playhub.com/graphql"
echo ""
print_status "Environment configuration:"
kubectl get configmap lobby-frontend-config -n $NAMESPACE -o yaml
