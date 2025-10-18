#!/bin/bash

# PlayHub Test Script
# This script runs all tests for the PlayHub project

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

# Function to run tests and capture exit code
run_tests() {
    local test_name="$1"
    local test_command="$2"
    local test_dir="$3"
    
    print_info "Running $test_name..."
    
    if [ ! -z "$test_dir" ]; then
        cd "$test_dir"
    fi
    
    if eval "$test_command"; then
        print_status "$test_name passed"
        return 0
    else
        print_error "$test_name failed"
        return 1
    fi
}

echo "🧪 Running PlayHub test suite..."
echo ""

# Check if we're in the right directory
if [ ! -f "README.md" ] || [ ! -d "backend" ] || [ ! -d "frontend" ]; then
    print_error "Please run this script from the PlayHub project root directory"
    exit 1
fi

# Track overall test results
OVERALL_RESULT=0

# Run backend tests
echo "🔧 Backend Tests"
echo "=================="

# Backend unit tests
if ! run_tests "Backend Unit Tests" "go test ./..." "backend"; then
    OVERALL_RESULT=1
fi

# Backend drift detection
if ! run_tests "Backend Drift Detection" "go test -v -run=TestGqlgenDrift ./graph" "backend"; then
    OVERALL_RESULT=1
fi

echo ""

# Run frontend tests
echo "🎨 Frontend Tests"
echo "=================="

# Frontend unit and integration tests
if ! run_tests "Frontend Unit & Integration Tests" "npm run test:run" "frontend"; then
    OVERALL_RESULT=1
fi

echo ""

# Run E2E tests (optional, can be slow)
if [ "$1" = "--e2e" ] || [ "$1" = "-e" ]; then
    echo "🌐 End-to-End Tests"
    echo "===================="
    
    if ! run_tests "Frontend E2E Tests" "npm run test:e2e" "frontend"; then
        OVERALL_RESULT=1
    fi
    
    echo ""
fi

# Summary
echo "📊 Test Summary"
echo "==============="

if [ $OVERALL_RESULT -eq 0 ]; then
    print_status "All tests passed! 🎉"
    echo ""
    echo "Your code is ready for commit!"
else
    print_error "Some tests failed. Please fix the issues before committing."
    echo ""
    echo "Run individual test suites for more details:"
    echo "  • Backend: cd backend && go test ./..."
    echo "  • Frontend: cd frontend && npm run test:run"
    echo "  • E2E: cd frontend && npm run test:e2e"
fi

exit $OVERALL_RESULT
