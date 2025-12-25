#!/bin/bash
set -e

SERVICE_NAME="zenbook-duo-daemon"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
INSTALL_DIR="/opt/zenbook-duo-daemon"
CONFIG_DIR="/etc/zenbook-duo-daemon"

check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo "Error: This script must be run as root" >&2
        exit 1
    fi
}

uninstall() {
    check_root

    echo "Stopping service..."
    systemctl stop "${SERVICE_NAME}" 2>/dev/null || true
    systemctl disable "${SERVICE_NAME}" 2>/dev/null || true

    echo "Removing service file..."
    rm -f "${SERVICE_FILE}"
    systemctl daemon-reload

    echo "Removing installation directory..."
    rm -rf "${INSTALL_DIR}"

    echo ""
    echo "Uninstallation complete."
    echo "Note: Config files in ${CONFIG_DIR} were NOT removed."
    echo "To remove config: sudo rm -rf ${CONFIG_DIR}"
}

uninstall "$@"
