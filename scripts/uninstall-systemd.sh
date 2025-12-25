#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Uninstalling Palm Rejection Daemon...${NC}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root (use sudo)${NC}"
    exit 1
fi

# Stop the service if running
echo -e "${YELLOW}Stopping service...${NC}"
systemctl stop palm-reject-daemon.service 2>/dev/null || true

# Disable the service
echo -e "${YELLOW}Disabling service...${NC}"
systemctl disable palm-reject-daemon.service 2>/dev/null || true

# Remove systemd service file
echo -e "${YELLOW}Removing systemd service...${NC}"
rm -f /etc/systemd/system/palm-reject-daemon.service

# Reload systemd
echo -e "${YELLOW}Reloading systemd...${NC}"
systemctl daemon-reload

# Remove binary
echo -e "${YELLOW}Removing binary...${NC}"
rm -f /usr/local/bin/palm-reject-daemon

echo -e "${GREEN}Uninstallation complete!${NC}"
echo -e "${YELLOW}Note: The palm rejection daemon has been removed from your system.${NC}"