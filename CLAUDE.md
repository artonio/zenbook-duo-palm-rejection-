# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

This project was developed using Claude Code Router with the model: **moonshotai/kimi-k2-instruct-0905**

## Project Overview

A Go-based daemon for ASUS ZenBook Duo that implements palm rejection by disabling the touchpad while typing. The system monitors keyboard input and temporarily disables touchpad input to prevent accidental cursor movement.

## Development Commands

### Build
```bash
go build -o palm-reject-daemon cmd/palm-reject-daemon/main.go
```

### Run
```bash
# Run the daemon
./palm-reject-daemon run

# Run with custom log level
LOG_LEVEL=debug ./palm-reject-daemon run
```

### Test
No test files currently exist in the project.

### Dependencies
```bash
go mod download
go mod tidy
```

## Architecture

### Core Components

1. **Device Discovery** (`internal/touchpad` - missing implementation)
   - `touchpad.FindAllTouchpadDevices()` - Locates all touchpad devices
   - `touchpad.FindKeyboardDevice()` - Finds the primary keyboard device
   - `touchpad.NewMultiController()` - Manages multiple touchpad devices
   - `touchpad.NewKeyboardMonitor()` - Monitors keyboard input events

2. **Event System** (`internal/events/`)
   - `KeyPressEventBus` - Handles keyboard event distribution
   - `SystemEventBus` - Manages system-wide events for component communication

3. **Typing Detection** (`internal/consumer/typing_detection.go`)
   - `TypingDetectionConsumer` - Core logic that disables touchpad on keypress
   - Configurable cooldown period (default: 300ms)
   - Thread-safe implementation with mutex protection

4. **Inter-Process Communication** (`internal/pipe/`)
   - Unix pipe receiver for external commands
   - Default pipe path: `/tmp/palm-reject-daemon.sock`

5. **Logging** (`pkg/logging/logging.go`)
   - Structured logging with zerolog
   - Log level controlled via `LOG_LEVEL` environment variable
   - Supports levels: trace, debug, info, warn, error, fatal

### Data Flow

1. Daemon starts and discovers touchpad/keyboard devices
2. Keyboard monitor captures keypress events
3. Typing detection consumer receives events and disables touchpad
4. Touchpad remains disabled for cooldown period (default 300ms)
5. System handles graceful shutdown on SIGINT/SIGTERM

### Missing Components

The following components are referenced but not implemented:
- `internal/touchpad` package (device discovery and control)
- Test files
- Build configuration (Makefile)
- Documentation (README.md)

## Environment Variables

- `LOG_LEVEL`: Controls logging verbosity (default: "info")