#!/bin/bash

# PlayHub Development Script
# This script starts both frontend and backend development servers

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

# Function to cleanup background processes
cleanup() {
    print_info "Shutting down development servers..."
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

echo "🚀 Starting PlayHub development servers..."
echo ""

# Check if we're in the right directory
if [ ! -f "README.md" ] || [ ! -d "backend" ] || [ ! -d "frontend" ]; then
    print_error "Please run this script from the PlayHub project root directory"
    exit 1
fi

# Start backend server
print_info "Starting backend server..."
cd backend

# Check if backend dependencies are installed
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    print_warning "Backend dependencies not found. Running setup..."
    go mod download
    go run github.com/99designs/gqlgen@v0.17.81 generate
fi

# Start backend in background
go run server.go &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Check if backend started successfully
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    print_error "Failed to start backend server"
    exit 1
fi

print_status "Backend server started (PID: $BACKEND_PID) - http://localhost:8080"

# Start frontend server
print_info "Starting frontend server..."
cd ../frontend

# Check if frontend dependencies are installed
if [ ! -d "node_modules" ]; then
    print_warning "Frontend dependencies not found. Installing..."
    npm install
fi

# Start frontend in background
npm run dev &
FRONTEND_PID=$!

# Wait a moment for frontend to start
sleep 3

# Check if frontend started successfully
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    print_error "Failed to start frontend server"
    cleanup
    exit 1
fi

print_status "Frontend server started (PID: $FRONTEND_PID) - http://localhost:5173"

echo ""
echo "🎉 Development servers are running!"
echo ""
echo "📱 Frontend: http://localhost:5173"
echo "🔧 Backend:  http://localhost:8080"
echo "📊 GraphQL:  http://localhost:8080/query"
echo ""
echo "Press Ctrl+C to stop all servers"
echo ""

# Wait for user to stop servers
wait
