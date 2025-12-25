#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo -e "${GREEN}Installing Palm Rejection Daemon...${NC}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root (use sudo)${NC}"
    exit 1
fi

# Build the daemon first
echo -e "${YELLOW}Building daemon...${NC}"
cd "$PROJECT_DIR"
./scripts/build.sh

# Copy binary to system location
echo -e "${YELLOW}Installing binary...${NC}"
cp bin/palm-reject-daemon /usr/local/bin/
chmod 755 /usr/local/bin/palm-reject-daemon

# Copy systemd service file
echo -e "${YELLOW}Installing systemd service...${NC}"
cp scripts/palm-reject-daemon.service /etc/systemd/system/

# Reload systemd
echo -e "${YELLOW}Reloading systemd...${NC}"
systemctl daemon-reload

# Enable and start the service
echo -e "${YELLOW}Enabling service...${NC}"
systemctl enable palm-reject-daemon.service

echo -e "${YELLOW}Starting service...${NC}"
systemctl start palm-reject-daemon.service

# Check status
echo -e "${YELLOW}Checking service status...${NC}"
sleep 2
systemctl status palm-reject-daemon.service --no-pager

echo -e "${GREEN}Installation complete!${NC}"
echo -e "${YELLOW}Commands:${NC}"
echo "  - Check status: systemctl status palm-reject-daemon"
echo "  - View logs: journalctl -u palm-reject-daemon -f"
echo "  - Stop: systemctl stop palm-reject-daemon"
echo "  - Start: systemctl start palm-reject-daemon"
echo "  - Disable: systemctl disable palm-reject-daemon"