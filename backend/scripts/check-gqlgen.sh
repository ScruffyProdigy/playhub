#!/bin/bash

# GQLGen Drift Detection Script
# This script checks if GraphQL schema files are newer than generated files
# and reminds developers to run 'gqlgen generate' if needed.

set -e

echo "🔍 Checking for gqlgen drift..."

# Change to the backend directory
cd "$(dirname "$0")/.."

# Run the drift detection test
if go test ./graph -run="TestGqlgenDrift" -v; then
    echo "✅ No drift detected - generated code is up to date."
    exit 0
else
    echo ""
    echo "❌ Drift detected! Some schema files are newer than generated files."
    echo ""
    echo "To fix this, run:"
    echo "  make generate"
    echo ""
    echo "Or manually:"
    echo "  go run github.com/99designs/gqlgen@v0.17.81 generate"
    echo ""
    exit 1
fi
