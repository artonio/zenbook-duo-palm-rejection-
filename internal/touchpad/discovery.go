// Package touchpad provides touchpad control for palm rejection during typing.
package touchpad

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

const (
	// inputDevDir is the directory containing input event devices
	inputDevDir = "/dev/input"
	// sysClassInput is the sysfs path for input device information
	sysClassInput = "/sys/class/input"
)

// DeviceInfo contains information about an input device
type DeviceInfo struct {
	// Path is the device path (e.g., /dev/input/event5)
	Path string
	// Name is the device name from sysfs
	Name string
}

// FindTouchpadDevice finds the first touchpad's evdev device path.
// Returns the path to /dev/input/eventX for the touchpad.
// Note: Use FindAllTouchpadDevices for systems with multiple touchpads.
func FindTouchpadDevice(logger zerolog.Logger) (*DeviceInfo, error) {
	devices, err := FindAllTouchpadDevices(logger)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("no touchpad device found")
	}
	return devices[0], nil
}

// FindAllTouchpadDevices finds ALL touchpad evdev devices.
// The Zenbook Duo has two screens, each with its own touchpad.
// Returns paths to /dev/input/eventX for all touchpads found.
func FindAllTouchpadDevices(logger zerolog.Logger) ([]*DeviceInfo, error) {
	entries, err := os.ReadDir(inputDevDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", inputDevDir, err)
	}

	var devices []*DeviceInfo

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "event") {
			continue
		}

		name := getDeviceNameFromSysfs(entry.Name())
		if name == "" {
			continue
		}

		logger.Debug().
			Str("device", entry.Name()).
			Str("name", name).
			Msg("Checking input device")

		if isTouchpadDevice(name) {
			path := filepath.Join(inputDevDir, entry.Name())
			logger.Info().
				Str("path", path).
				Str("name", name).
				Msg("Found touchpad device")
			devices = append(devices, &DeviceInfo{Path: path, Name: name})
		}
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no touchpad device found")
	}

	logger.Info().Int("count", len(devices)).Msg("Total touchpad devices found")
	return devices, nil
}

// FindKeyboardDevice finds the keyboard's evdev device path.
// This is needed to monitor for typing activity (regular keypresses, not Fn keys).
// Returns the path to /dev/input/eventX for the keyboard.
// Prefers keyd virtual keyboard if present (keyd grabs the physical keyboard).
func FindKeyboardDevice(logger zerolog.Logger) (*DeviceInfo, error) {
	entries, err := os.ReadDir(inputDevDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", inputDevDir, err)
	}

	var fallback *DeviceInfo

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "event") {
			continue
		}

		name := getDeviceNameFromSysfs(entry.Name())
		if name == "" {
			continue
		}

		logger.Debug().
			Str("device", entry.Name()).
			Str("name", name).
			Msg("Checking input device for keyboard")

		// Prefer keyd virtual keyboard (it grabs the physical keyboard)
		if strings.Contains(name, "keyd virtual keyboard") {
			path := filepath.Join(inputDevDir, entry.Name())
			logger.Info().
				Str("path", path).
				Str("name", name).
				Msg("Found keyboard device")
			return &DeviceInfo{Path: path, Name: name}, nil
		}

		// Keep AT Translated as fallback
		if strings.Contains(name, "AT Translated Set 2 keyboard") && fallback == nil {
			fallback = &DeviceInfo{
				Path: filepath.Join(inputDevDir, entry.Name()),
				Name: name,
			}
		}
	}

	if fallback != nil {
		logger.Info().
			Str("path", fallback.Path).
			Str("name", fallback.Name).
			Msg("Found keyboard device")
		return fallback, nil
	}

	return nil, fmt.Errorf("no keyboard device found")
}

// getDeviceNameFromSysfs reads device name from sysfs WITHOUT opening the evdev device.
// This is safe to call on any input device without affecting the input stack.
// Returns empty string if device name cannot be determined.
func getDeviceNameFromSysfs(eventName string) string {
	// Path: /sys/class/input/eventX/device/name
	namePath := filepath.Join(sysClassInput, eventName, "device", "name")
	data, err := os.ReadFile(namePath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// isTouchpadDevice checks if the device name indicates a touchpad.
func isTouchpadDevice(name string) bool {
	nameLower := strings.ToLower(name)
	// Check for common touchpad identifiers
	return strings.Contains(nameLower, "touchpad") ||
		strings.Contains(nameLower, "trackpad") ||
		// ASUS-specific: the integrated touchpad may have specific names
		(strings.Contains(nameLower, "asus") && strings.Contains(nameLower, "touch"))
}

// isKeyboardDevice checks if the device name indicates a keyboard for typing detection.
// Priority order:
// 1. "keyd virtual keyboard" - if keyd is running, it grabs the physical keyboard
// 2. "AT Translated Set 2 keyboard" - standard internal keyboard on most laptops
func isKeyboardDevice(name string) bool {
	// keyd virtual keyboard - if keyd (key remapper) is running, it GRABs the physical
	// keyboard and emits events through this virtual device
	if strings.Contains(name, "keyd virtual keyboard") {
		return true
	}
	// Standard internal keyboard - receives all regular keypresses (a-z, numbers, etc.)
	// This is the primary keyboard device on most laptops
	if strings.Contains(name, "AT Translated Set 2 keyboard") {
		return true
	}
	return false
}

// IsTouchpadPresent checks if a touchpad device exists.
func IsTouchpadPresent() bool {
	entries, err := os.ReadDir(inputDevDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "event") {
			continue
		}
		name := getDeviceNameFromSysfs(entry.Name())
		if isTouchpadDevice(name) {
			return true
		}
	}
	return false
}

// IsKeyboardPresent checks if the ASUS keyboard device exists.
func IsKeyboardPresent() bool {
	entries, err := os.ReadDir(inputDevDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "event") {
			continue
		}
		name := getDeviceNameFromSysfs(entry.Name())
		if isKeyboardDevice(name) {
			return true
		}
	}
	return false
}
