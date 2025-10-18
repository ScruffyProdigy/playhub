#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

echo "🐳 Building and pushing PlayHub Docker images..."

# Configuration
DOCKER_REGISTRY="docker.io"
DOCKER_USERNAME="scruffyprodigy"
BACKEND_IMAGE="playhub-backend"
FRONTEND_IMAGE="playhub-frontend"
TAG="latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if user is logged into Docker Hub
if ! docker info | grep -q "Username"; then
    print_warning "You may need to log in to Docker Hub:"
    echo "docker login"
    echo ""
fi

# Build backend image
print_status "Building backend image..."
cd backend
docker build -t ${DOCKER_REGISTRY}/${DOCKER_USERNAME}/${BACKEND_IMAGE}:${TAG} .
if [ $? -eq 0 ]; then
    print_status "Backend image built successfully"
else
    print_error "Failed to build backend image"
    exit 1
fi

# Build frontend image
print_status "Building frontend image..."
cd ../frontend
docker build -t ${DOCKER_REGISTRY}/${DOCKER_USERNAME}/${FRONTEND_IMAGE}:${TAG} .
if [ $? -eq 0 ]; then
    print_status "Frontend image built successfully"
else
    print_error "Failed to build frontend image"
    exit 1
fi

cd ..

# Push images (optional - uncomment if you want to push to registry)
# print_status "Pushing backend image to registry..."
# docker push ${DOCKER_REGISTRY}/${DOCKER_USERNAME}/${BACKEND_IMAGE}:${TAG}
# 
# print_status "Pushing frontend image to registry..."
# docker push ${DOCKER_REGISTRY}/${DOCKER_USERNAME}/${FRONTEND_IMAGE}:${TAG}

print_status "Docker images built successfully!"
print_status "Backend: ${DOCKER_REGISTRY}/${DOCKER_USERNAME}/${BACKEND_IMAGE}:${TAG}"
print_status "Frontend: ${DOCKER_REGISTRY}/${DOCKER_USERNAME}/${FRONTEND_IMAGE}:${TAG}"
echo ""
print_warning "To push to Docker Hub, uncomment the push commands in this script and run again."
print_warning "Make sure you're logged in with: docker login"
