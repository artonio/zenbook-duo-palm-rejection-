# Palm Rejection Daemon - SystemD Installation

This guide explains how to install the Palm Rejection Daemon as a systemd service on Fedora Linux.

## Quick Install

1. **Install the daemon as a systemd service:**
   ```bash
   sudo ./scripts/install-systemd.sh
   ```

2. **Check if it's running:**
   ```bash
   systemctl status palm-reject-daemon
   ```

3. **View logs:**
   ```bash
   journalctl -u palm-reject-daemon -f
   ```

## Using the Daemon Manager

For easier management, use the daemon manager script:

```bash
# Show all available commands
./scripts/daemon-manager.sh

# Common commands
sudo ./scripts/daemon-manager.sh install     # Install and start
sudo ./scripts/daemon-manager.sh uninstall  # Remove completely
./scripts/daemon-manager.sh status        # Check status
./scripts/daemon-manager.sh logs        # View real-time logs
./scripts/daemon-manager.sh restart       # Restart the daemon
```

## Manual SystemD Commands

```bash
# Start the daemon
sudo systemctl start palm-reject-daemon

# Stop the daemon
sudo systemctl stop palm-reject-daemon

# Enable auto-start on boot
sudo systemctl enable palm-reject-daemon

# Disable auto-start
sudo systemctl disable palm-reject-daemon

# Check status
systemctl status palm-reject-daemon

# View logs
journalctl -u palm-reject-daemon -f
```

## Testing Safely

The daemon includes a timeout feature for safe testing:

```bash
# Run with 10-second timeout (default in run.sh)
./scripts/run.sh

# Run with custom timeout
./scripts/run.sh --timeout 30s

# Run without timeout (indefinitely)
./scripts/run.sh --timeout 0
```

## Configuration

The systemd service runs with these settings:
- **User**: root (required for input device access)
- **Auto-restart**: Yes, on failure
- **Restart delay**: 5 seconds
- **Security**: Restricted permissions for safety

## Troubleshooting

1. **Permission denied errors:**
   The daemon needs root access to read input devices. The systemd service runs as root.

2. **Daemon not starting:**
   ```bash
   # Check logs
   journalctl -u palm-reject-daemon --no-pager

   # Check if binary exists
   ls -la /usr/local/bin/palm-reject-daemon
   ```

3. **Touchpad not being disabled:**
   ```bash
   # Check if devices are found
   sudo LOG_LEVEL=debug /usr/local/bin/palm-reject-daemon run --timeout 10s
   ```

4. **Remove completely:**
   ```bash
   sudo ./scripts/uninstall-systemd.sh
   ```