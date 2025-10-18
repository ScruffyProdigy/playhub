#!/bin/bash

# PlayHub Setup Script
# This script sets up the development environment for PlayHub

set -e

echo "🚀 Setting up PlayHub development environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    print_status "Go $GO_VERSION found"
else
    print_error "Go not found. Please install Go 1.25+ from https://golang.org/dl/"
    exit 1
fi

# Check Node.js
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version | sed 's/v//')
    print_status "Node.js $NODE_VERSION found"
else
    print_error "Node.js not found. Please install Node.js 20+ from https://nodejs.org/"
    exit 1
fi

# Check Docker
if command -v docker &> /dev/null; then
    print_status "Docker found"
else
    print_warning "Docker not found. You'll need it for database and deployment."
fi

# Setup backend
echo "🔧 Setting up backend..."
cd backend

# Install Go dependencies
print_status "Installing Go dependencies..."
go mod download

# Generate GraphQL code
print_status "Generating GraphQL code..."
go run github.com/99designs/gqlgen@v0.17.81 generate

# Run backend tests
print_status "Running backend tests..."
go test ./...

print_status "Backend setup complete!"

# Setup frontend
echo "🔧 Setting up frontend..."
cd ../frontend

# Install npm dependencies
print_status "Installing npm dependencies..."
npm install

# Run frontend tests
print_status "Running frontend tests..."
npm run test:run

print_status "Frontend setup complete!"

# Make scripts executable
echo "🔧 Making scripts executable..."
cd ..
chmod +x scripts/*.sh

print_status "All scripts are now executable!"

echo ""
echo "🎉 PlayHub setup complete!"
echo ""
echo "Next steps:"
echo "  • Run './scripts/dev.sh' to start development servers"
echo "  • Run './scripts/test.sh' to run all tests"
echo "  • Read docs/development.md for more information"
echo ""
echo "Happy coding! 🚀"
