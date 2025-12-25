package touchpad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTouchpadDevice(t *testing.T) {
	tests := []struct {
		name     string
		device   string
		expected bool
	}{
		{"Synaptics Touchpad", "SynPS/2 Synaptics TouchPad", true},
		{"Elan Touchpad", "ELAN1200:00 04F3:307A Touchpad", true},
		{"Trackpad", "Apple Inc. Apple Internal Keyboard / Trackpad", true},
		{"ASUS Touch", "ASUE140A:00 04F3:3134 Touchpad", true},
		{"Uppercase TOUCHPAD", "DELL Touchpad", true},
		{"Not a touchpad", "AT Translated Set 2 keyboard", false},
		{"Mouse device", "Logitech USB Mouse", false},
		{"Empty string", "", false},
		{"ASUS but no touch", "ASUE140A:00 04F3:3134 Keyboard", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTouchpadDevice(tt.device)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsKeyboardDevice(t *testing.T) {
	tests := []struct {
		name     string
		device   string
		expected bool
	}{
		{"keyd virtual keyboard", "keyd virtual keyboard", true},
		{"AT Translated", "AT Translated Set 2 keyboard", true},
		{"Not keyboard", "SynPS/2 Synaptics TouchPad", false},
		{"Empty string", "", false},
		{"Partial match", "AT Translated", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKeyboardDevice(tt.device)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDeviceNameFromSysfs(t *testing.T) {
	// Skip this test if running in an environment with real devices
	// as we can't easily mock the filesystem path
	t.Skip("Cannot mock sysfs path in current implementation")
}

func TestDeviceInfo(t *testing.T) {
	info := &DeviceInfo{
		Path: "/dev/input/event5",
		Name: "ELAN1200:00 Touchpad",
	}

	assert.Equal(t, "/dev/input/event5", info.Path)
	assert.Equal(t, "ELAN1200:00 Touchpad", info.Name)
}

func TestDiscoveryWithLogger(t *testing.T) {
	// Test with mock functions
	t.Run("Test functions with logger", func(t *testing.T) {
		// These would normally interact with the filesystem
		// For unit tests, we verify the functions exist and accept the right parameters
		assert.NotPanics(t, func() {
			IsTouchpadPresent()
		})

		assert.NotPanics(t, func() {
			IsKeyboardPresent()
		})
	})
}