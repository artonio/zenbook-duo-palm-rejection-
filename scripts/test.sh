#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "Running tests..."

# Run all tests with verbose output and race detector
go test -v -race ./...

# Run tests with coverage
echo ""
echo "Running tests with coverage..."
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

echo ""
echo "Tests complete!"

# Optional: generate HTML coverage report
if [ "$1" == "--html" ]; then
    go tool cover -html=coverage.out -o coverage.html
    echo "Coverage report generated: coverage.html"
fi
