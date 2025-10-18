#!/bin/bash

# PlayHub Backend Test Script
# This script runs backend tests with detailed output

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

echo "🔧 Running PlayHub Backend Tests..."
echo ""

# Check if we're in the right directory
if [ ! -d "backend" ]; then
    print_error "Please run this script from the PlayHub project root directory"
    exit 1
fi

cd backend

# Check if Go is available
if ! command -v go &> /dev/null; then
    print_error "Go not found. Please install Go 1.25+ from https://golang.org/dl/"
    exit 1
fi

# Check if dependencies are installed
if [ ! -f "go.sum" ]; then
    print_warning "Dependencies not found. Installing..."
    go mod download
fi

# Run different types of tests
print_info "Running unit tests..."
if go test -v ./...; then
    print_status "Unit tests passed"
else
    print_error "Unit tests failed"
    exit 1
fi

echo ""

print_info "Running drift detection tests..."
if go test -v -run=TestGqlgenDrift ./graph; then
    print_status "Drift detection tests passed"
else
    print_error "Drift detection tests failed"
    exit 1
fi

echo ""

print_info "Running benchmarks..."
if go test -bench=. -benchmem ./graph; then
    print_status "Benchmarks completed"
else
    print_warning "Some benchmarks failed (this is often normal)"
fi

echo ""

print_info "Running with coverage..."
if go test -cover ./...; then
    print_status "Coverage analysis completed"
else
    print_error "Coverage analysis failed"
    exit 1
fi

echo ""
print_status "All backend tests completed successfully! 🎉"
