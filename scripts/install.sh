#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SERVICE_NAME="zenbook-duo-daemon"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
INSTALL_DIR="/opt/zenbook-duo-daemon"
CONFIG_DIR="/etc/zenbook-duo-daemon"
BINARY_PATH="${INSTALL_DIR}/zenbook-duo-daemon"

check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo "Error: This script must be run as root" >&2
        exit 1
    fi
}

install() {
    check_root

    cd "$PROJECT_DIR"

    echo "Building..."
    ./scripts/build.sh

    echo "Stopping existing service..."
    systemctl stop "${SERVICE_NAME}" 2>/dev/null || true
    systemctl disable "${SERVICE_NAME}" 2>/dev/null || true

    echo "Installing binary to ${INSTALL_DIR}..."
    mkdir -p "${INSTALL_DIR}"
    cp bin/zenbook-duo-daemon "${BINARY_PATH}"
    chmod +x "${BINARY_PATH}"

    echo "Installing systemd service..."
    cp systemd/zenbook-duo-daemon.service "${SERVICE_FILE}"

    echo "Creating config directory..."
    mkdir -p "${CONFIG_DIR}"
    if [ ! -f "${CONFIG_DIR}/config.toml" ]; then
        cp configs/config.example.toml "${CONFIG_DIR}/config.toml"
        echo "Created default config at ${CONFIG_DIR}/config.toml"
    else
        echo "Config file already exists, not overwriting"
    fi

    echo "Reloading systemd..."
    systemctl daemon-reload

    echo "Enabling and starting service..."
    systemctl enable "${SERVICE_NAME}"
    systemctl start "${SERVICE_NAME}"

    echo ""
    echo "Installation complete!"
    echo ""
    systemctl status "${SERVICE_NAME}" --no-pager || true
}

install "$@"
