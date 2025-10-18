#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

echo "🚀 Deploying PlayHub to Local Kubernetes..."

# Configuration
NAMESPACE="playhub"
CONTEXT="minikube"

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
kubectl apply -f k8s/base/namespace.yaml

# Apply base configurations (backend, ingress)
print_step "Applying base configurations..."
kubectl apply -f k8s/base/backend.yaml
kubectl apply -f k8s/base/ingress.yaml

# Apply local environment configuration
print_step "Applying local environment configuration..."
kubectl apply -f k8s/env/local.yaml

# Wait for deployments to be ready
print_step "Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/lobby-backend -n $NAMESPACE
kubectl wait --for=condition=available --timeout=300s deployment/lobby-frontend -n $NAMESPACE

# Show deployment status
print_step "Deployment status:"
kubectl get pods -n $NAMESPACE
kubectl get services -n $NAMESPACE

print_status "Local deployment completed successfully!"
echo ""
print_status "Access your application:"
echo "  Frontend: http://localhost:8080 (with port-forward)"
echo "  Backend API: http://localhost:8081 (with port-forward)"
echo "  GraphQL: http://localhost:8081/graphql"
echo ""
print_warning "To access the application, run:"
echo "  kubectl port-forward -n $NAMESPACE svc/lobby-frontend 8080:80 &"
echo "  kubectl port-forward -n $NAMESPACE svc/lobby-backend 8081:80 &"
echo ""
print_status "Environment configuration:"
kubectl get configmap lobby-frontend-config -n $NAMESPACE -o yaml
