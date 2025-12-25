#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Find go binary - check common locations if not in PATH
GO_BIN="go"
if ! command -v go &> /dev/null; then
    # Common Go installation paths
    for path in /usr/local/go/bin/go /usr/bin/go "$HOME/go/bin/go" /snap/bin/go; do
        if [ -x "$path" ]; then
            GO_BIN="$path"
            break
        fi
    done
fi

if ! command -v "$GO_BIN" &> /dev/null && [ "$GO_BIN" = "go" ]; then
    echo "Error: go not found in PATH or common locations"
    echo "Please install Go or add it to your PATH"
    exit 1
fi

echo "Building palm-reject-daemon..."
echo "Using Go: $GO_BIN"

# Ensure bin directory exists
mkdir -p bin

# Build for current platform
"$GO_BIN" build -ldflags="-s -w" -o bin/palm-reject-daemon ./cmd/palm-reject-daemon

echo "Build complete: bin/palm-reject-daemon"

# Optionally build with race detector for testing
if [ "$1" == "--race" ]; then
    echo "Building with race detector..."
    "$GO_BIN" build -race -o bin/palm-reject-daemon-race ./cmd/palm-reject-daemon
    echo "Race build complete: bin/palm-reject-daemon-race"
fi
