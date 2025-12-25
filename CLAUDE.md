# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

This project was developed using Claude Code Router with the model: **moonshotai/kimi-k2-instruct-0905**

## Project Overview

A Go-based daemon for ASUS ZenBook Duo that implements palm rejection by disabling the touchpad while typing. The system monitors keyboard input and temporarily disables touchpad input to prevent accidental cursor movement.

## Architecture

### Core Components

1. **Device Discovery** (`internal/touchpad/`)
   - `FindAllTouchpadDevices()` - Locates all touchpad devices (Zenbook Duo has 3)
   - `FindKeyboardDevice()` - Finds the primary keyboard device
   - `MultiController` - Manages multiple touchpad devices simultaneously
   - `KeyboardMonitor` - Monitors keyboard input events via evdev

2. **Event System** (`internal/events/`)
   - `SystemEventBus` - Broadcast pattern for system-wide events
   - Used for inter-component communication

3. **Typing Detection** (`internal/consumer/`)
   - `TypingDetectionConsumer` - Core logic that disables touchpad on keypress
   - 300ms hard-coded cooldown period
   - Thread-safe with mutex protection

4. **Inter-Process Communication** (`internal/pipe/`)
   - Unix pipe receiver for external commands
   - Default pipe path: `/tmp/zenbook-duo-daemon.pipe`
   - Commands: `touchpad_disable`, `touchpad_enable`, `touchpad_toggle`

5. **Logging** (`pkg/logging/`)
   - Structured logging with zerolog
   - Log level controlled via `LOG_LEVEL` environment variable
   - Levels: trace, debug, info, warn, error, fatal

### Data Flow

1. Daemon starts and discovers touchpad/keyboard devices via evdev
2. Keyboard monitor captures keypress events
3. Typing detection consumer receives events and disables all touchpads
4. Touchpads remain disabled for cooldown period (300ms)
5. System handles graceful shutdown on SIGINT/SIGTERM

## Development Commands

### Build
```bash
# Quick build
go build -o palm-reject-daemon ./cmd/palm-reject-daemon

# Using build script (recommended)
./scripts/build.sh

# Build with race detector for testing
./scripts/build.sh --race
```

### Run
```bash
# Run the daemon
./palm-reject-daemon run

# Run with debug logging
LOG_LEVEL=debug ./palm-reject-daemon run

# Run with timeout for safe testing
./scripts/run.sh --timeout 10s

# Run without timeout (indefinitely)
./scripts/run.sh --timeout 0
```

### Systemd Service Management
```bash
# Install as system service
sudo ./scripts/install-systemd.sh

# Check status
systemctl status palm-reject-daemon

# View logs
journalctl -u palm-reject-daemon -f

# Restart daemon
sudo systemctl restart palm-reject-daemon
# OR
./scripts/daemon-manager.sh restart
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

## Important Limitations

- **Keyboard Support**: Only works with the built-in keyboard attached to the Zenbook Duo
- **No Bluetooth Support**: Does not support Bluetooth keyboards
- **Hot-Plugging**: Detaching and re-attaching the keyboard does not automatically restart the daemon
- **Hard-coded Cooldown**: 300ms cooldown period is not configurable

## Key Technical Details

- **Root Required**: Must run as root to access input devices via evdev
- **Multi-touchpad**: Supports all 3 touchpads found on Zenbook Duo
- **Memory Usage**: ~6MB RAM when running
- **Binary Location**: Installs to `/usr/local/bin/palm-reject-daemon`
- **Service File**: `/etc/systemd/system/palm-reject-daemon.service`

## Common Development Tasks

### Testing Device Discovery
```bash
# Check what devices are found
sudo LOG_LEVEL=debug ./bin/palm-reject-daemon run --timeout 10s
```

### Manual Touchpad Control
```bash
# Disable touchpad
echo "touchpad_disable" > /tmp/zenbook-duo-daemon.pipe

# Enable touchpad
echo "touchpad_enable" > /tmp/zenbook-duo-daemon.pipe

# Toggle touchpad
echo "touchpad_toggle" > /tmp/zenbook-duo-daemon.pipe
```

### Debugging
```bash
# View recent logs
journalctl -u palm-reject-daemon --no-pager -n 20

# Check if service is active
systemctl is-active palm-reject-daemon
```