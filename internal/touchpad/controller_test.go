package touchpad

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestController_BasicOperations(t *testing.T) {
	logger := zerolog.Nop()
	controller := NewController("/dev/input/event5", logger)

	// Test initial state
	assert.Equal(t, "/dev/input/event5", controller.DevicePath())
	assert.False(t, controller.IsOpen())
	assert.False(t, controller.IsDisabled())

	// Test operations on closed device
	err := controller.Disable()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "touchpad device not open")

	err = controller.Enable()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "touchpad device not open")

	// Test Stop on unopened device
	err = controller.Stop()
	assert.NoError(t, err)
}

func TestController_OpenClose(t *testing.T) {
	logger := zerolog.Nop()
	controller := NewController("/dev/input/event99", logger) // Non-existent device

	// Test opening non-existent device
	err := controller.Open()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open touchpad device")

	// Test closing unopened device
	err = controller.Close()
	assert.NoError(t, err)
}

func TestController_StateManagement(t *testing.T) {
	logger := zerolog.Nop()
	controller := NewController("/dev/input/event5", logger)

	// Test state transitions
	assert.False(t, controller.IsDisabled())

	// Try to disable without opening
	err := controller.Disable()
	assert.Error(t, err)
	assert.False(t, controller.IsDisabled())
}

func TestController_DeviceInfo(t *testing.T) {
	info := &DeviceInfo{
		Path: "/dev/input/event5",
		Name: "ELAN1200:00 Touchpad",
	}

	assert.Equal(t, "/dev/input/event5", info.Path)
	assert.Equal(t, "ELAN1200:00 Touchpad", info.Name)
}