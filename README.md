# Palm Rejection Daemon for ASUS Zenbook Duo 2024

A lightweight system daemon that automatically disables the touchpad while typing to prevent accidental cursor movements (palm rejection).

## Features

- **Automatic palm rejection** - Disables touchpad while typing
- **Multi-touchpad support** - Works with multiple touchpad devices
- **Configurable cooldown** - 300ms default (hard-coded)
- **Systemd integration** - Runs as a system service
- **Lightweight** - Uses only ~6MB RAM
- **Safe timeout** - Includes timeout feature for testing
- **Pipe commands** - Manual touchpad control via Unix pipe

## Installation

### Quick Install (Systemd Service)

```bash
# Install as a system service (requires sudo)
sudo ./scripts/install-systemd.sh

# Check if it's running
systemctl status palm-reject-daemon
```

### Manual Build

```bash
# Build the daemon
./scripts/build.sh

# Run with timeout for testing
./scripts/run.sh --timeout 10s
```

## Usage

### Managing the Daemon

Use the daemon manager script for easy control:

```bash
# Show all commands
./scripts/daemon-manager.sh

# Common operations
sudo ./scripts/daemon-manager.sh install     # Install and enable
./scripts/daemon-manager.sh status          # Check status
./scripts/daemon-manager.sh logs            # View real-time logs
./scripts/daemon-manager.sh restart         # Restart daemon
```

### Systemd Commands

```bash
# Check status
systemctl status palm-reject-daemon

# View logs
journalctl -u palm-reject-daemon -f

# Control the service
sudo systemctl start palm-reject-daemon
sudo systemctl stop palm-reject-daemon
sudo systemctl restart palm-reject-daemon

# Enable/disable auto-start
sudo systemctl enable palm-reject-daemon
sudo systemctl disable palm-reject-daemon
```

### Safe Testing

```bash
# Run with 10-second timeout (default)
./scripts/run.sh

# Run with custom timeout
./scripts/run.sh --timeout 30s

# Run without timeout (indefinitely)
./scripts/run.sh --timeout 0
```

## Pipe Commands

The daemon accepts commands via Unix pipe for manual touchpad control:

```bash
# Send commands to the daemon
echo "touchpad_disable" > /tmp/zenbook-duo-daemon.pipe
echo "touchpad_enable" > /tmp/zenbook-duo-daemon.pipe
echo "touchpad_toggle" > /tmp/zenbook-duo-daemon.pipe
```

## How It Works

1. **Device Discovery** - Automatically finds all touchpad and keyboard devices
2. **Event Monitoring** - Monitors keyboard events in real-time
3. **Touchpad Control** - Disables touchpad when keys are pressed
4. **Cooldown Period** - Re-enables touchpad after 300ms of no typing
5. **Multi-device Support** - Handles multiple touchpads simultaneously

## Project Structure

```
zenbook-duo-palm-rejection/
├── cmd/palm-reject-daemon/     # Main daemon entry point
├── internal/
│   ├── consumer/              # Typing detection logic
│   ├── events/                # Event system
│   ├── pipe/                  # Unix pipe receiver
│   └── touchpad/              # Touchpad control
├── pkg/logging/               # Logging utilities
├── scripts/
│   ├── build.sh              # Build script
│   ├── run.sh                # Run with timeout
│   ├── install-systemd.sh     # Install as service
│   ├── uninstall-systemd.sh   # Remove service
│   ├── daemon-manager.sh      # Easy management
│   └── palm-reject-daemon.service  # Systemd unit file
├── go.mod                     # Go dependencies
├── SYSTEMD.md                 # Systemd installation guide
└── README.md                  # This file
```

## Requirements

- Linux with systemd
- Go 1.19+ (for building)
- Root access (for input device access)

## Troubleshooting

### Permission Denied Errors

The daemon needs root access to read input devices. The systemd service runs as root.

### Daemon Not Starting

```bash
# Check logs
journalctl -u palm-reject-daemon --no-pager

# Check if binary exists
ls -la /usr/local/bin/palm-reject-daemon
```

### Touchpad Not Being Disabled

```bash
# Check if devices are found
sudo LOG_LEVEL=debug /usr/local/bin/palm-reject-daemon run --timeout 10s

# Check if service is running
systemctl is-active palm-reject-daemon
```

### Remove Completely

```bash
sudo ./scripts/uninstall-systemd.sh
```

## Development

### Building

```bash
./scripts/build.sh
```

### Testing

```bash
# Run with debug logging
LOG_LEVEL=debug ./bin/palm-reject-daemon run --timeout 10s
```

### Dependencies

- [go-evdev](https://github.com/holoplot/go-evdev) - Input device access
- [zerolog](https://github.com/rs/zerolog) - Structured logging
- [cobra](https://github.com/spf13/cobra) - CLI framework

## License

This project is part of the zenbook-duo-keyboard-daemon project.

## See Also

- [SYSTEMD.md](SYSTEMD.md) - Detailed systemd installation guide
- [Keyboard Daemon](https://github.com/artonio/zenbook-duo-keyboard-daemon-go) - Full keyboard daemon with palm rejection