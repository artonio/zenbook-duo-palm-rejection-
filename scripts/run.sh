#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Warn if running with sudo (build should be done as normal user)
if [ "$EUID" -eq 0 ]; then
    echo "Warning: This script should NOT be run with sudo directly."
    echo "It will use sudo only for running the binary."
    echo ""
    echo "Usage: ./scripts/run.sh [--test-mode]"
    echo ""
    # If binary exists, skip build and just run
    if [ -f "./bin/palm-reject-daemon" ]; then
        echo "Binary exists, skipping build..."
    else
        echo "Error: Binary not found and cannot build as root."
        echo "Please build first as normal user: ./scripts/build.sh"
        exit 1
    fi
else
    # Build first (as normal user)
    ./scripts/build.sh
fi

echo ""
echo "Running daemon (requires sudo for USB/evdev access)..."
echo "Press Ctrl+C to stop"
echo ""

# Default timeout for safety (10 seconds)
# Override with --timeout 0 to disable
DEFAULT_TIMEOUT="--timeout 10s"

# Check if user provided their own timeout
if [[ "$*" == *"--timeout"* ]]; then
    DEFAULT_TIMEOUT=""
fi

# Run with debug logging
# Pass through any arguments (like --test-mode)
if [ "$EUID" -eq 0 ]; then
    # Already root
    LOG_LEVEL=debug ./bin/palm-reject-daemon run $DEFAULT_TIMEOUT "$@"
else
    # Need sudo
    sudo LOG_LEVEL=debug ./bin/palm-reject-daemon run $DEFAULT_TIMEOUT "$@"
fi
