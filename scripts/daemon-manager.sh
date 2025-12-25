#!/bin/bash
# Palm Rejection Daemon Manager

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

show_help() {
    echo -e "${GREEN}Palm Rejection Daemon Manager${NC}"
    echo "Usage: $0 [command]"
    echo ""
    echo -e "${YELLOW}Commands:${NC}"
    echo "  install     - Install and enable the systemd service (requires sudo)"
    echo "  uninstall   - Remove the systemd service (requires sudo)"
    echo "  start       - Start the daemon"
    echo "  stop        - Stop the daemon"
    echo "  restart     - Restart the daemon"
    echo "  status      - Show daemon status"
    echo "  logs        - Show real-time logs"
    echo "  logs-last   - Show last 50 log entries"
    echo "  enable      - Enable auto-start on boot"
    echo "  disable     - Disable auto-start on boot"
    echo "  run         - Run in foreground with debug logging (safe testing)"
    echo "  run-short   - Run with 10s timeout for quick testing"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  sudo $0 install"
    echo "  $0 status"
    echo "  $0 logs"
}

case "$1" in
    install)
        sudo ./scripts/install-systemd.sh
        ;;
    uninstall)
        sudo ./scripts/uninstall-systemd.sh
        ;;
    start)
        sudo systemctl start palm-reject-daemon
        echo -e "${GREEN}Daemon started${NC}"
        ;;
    stop)
        sudo systemctl stop palm-reject-daemon
        echo -e "${GREEN}Daemon stopped${NC}"
        ;;
    restart)
        sudo systemctl restart palm-reject-daemon
        echo -e "${GREEN}Daemon restarted${NC}"
        ;;
    status)
        systemctl status palm-reject-daemon --no-pager
        ;;
    logs)
        sudo journalctl -u palm-reject-daemon -f
        ;;
    logs-last)
        sudo journalctl -u palm-reject-daemon -n 50 --no-pager
        ;;
    enable)
        sudo systemctl enable palm-reject-daemon
        echo -e "${GREEN}Auto-start enabled${NC}"
        ;;
    disable)
        sudo systemctl disable palm-reject-daemon
        echo -e "${GREEN}Auto-start disabled${NC}"
        ;;
    run)
        echo -e "${YELLOW}Running in foreground mode (Ctrl+C to stop)...${NC}"
        sudo LOG_LEVEL=debug ./bin/palm-reject-daemon run
        ;;
    run-short)
        echo -e "${YELLOW}Running with 10s timeout for testing...${NC}"
        sudo LOG_LEVEL=debug ./bin/palm-reject-daemon run --timeout 10s
        ;;
    *)
        show_help
        exit 1
        ;;
esac