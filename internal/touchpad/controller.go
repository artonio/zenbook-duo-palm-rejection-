package touchpad

import (
	"fmt"
	"sync"

	evdev "github.com/holoplot/go-evdev"
	"github.com/rs/zerolog"
)

// TouchpadController is the interface for touchpad control.
// Both Controller (single device) and MultiController implement this.
type TouchpadController interface {
	Disable() error
	Enable() error
	IsDisabled() bool
	Stop() error
}

// Controller manages touchpad enable/disable state using evdev GRAB.
// When grabbed, the touchpad device is exclusively owned by this process,
// preventing events from reaching other applications.
type Controller struct {
	devicePath string
	device     *evdev.InputDevice
	grabbed    bool
	mu         sync.Mutex
	logger     zerolog.Logger
}

// NewController creates a new touchpad controller.
// The device is not opened until Open() is called.
func NewController(devicePath string, logger zerolog.Logger) *Controller {
	return &Controller{
		devicePath: devicePath,
		logger:     logger.With().Str("component", "touchpad_ctrl").Str("device", devicePath).Logger(),
	}
}

// Open opens the touchpad device for control.
// Must be called before Disable/Enable.
func (c *Controller) Open() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.device != nil {
		return nil // Already open
	}

	dev, err := evdev.Open(c.devicePath)
	if err != nil {
		return fmt.Errorf("failed to open touchpad device %s: %w", c.devicePath, err)
	}

	name, err := dev.Name()
	if err != nil {
		dev.Close()
		return fmt.Errorf("failed to get device name: %w", err)
	}

	c.device = dev
	c.logger.Info().Str("name", name).Msg("Touchpad controller opened")
	return nil
}

// Close closes the touchpad device.
// If the touchpad is currently disabled, it will be re-enabled first.
func (c *Controller) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.device == nil {
		return nil
	}

	// Ensure touchpad is enabled before closing
	if c.grabbed {
		if err := c.device.Ungrab(); err != nil {
			c.logger.Warn().Err(err).Msg("Failed to ungrab touchpad during close")
		}
		c.grabbed = false
	}

	err := c.device.Close()
	c.device = nil
	c.logger.Info().Msg("Touchpad controller closed")
	return err
}

// Disable disables the touchpad by grabbing it for exclusive access.
// When grabbed, touchpad events don't reach other applications.
func (c *Controller) Disable() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.device == nil {
		return fmt.Errorf("touchpad device not open")
	}

	if c.grabbed {
		return nil // Already disabled
	}

	if err := c.device.Grab(); err != nil {
		return fmt.Errorf("failed to grab touchpad: %w", err)
	}

	c.grabbed = true
	c.logger.Debug().Msg("Touchpad disabled (grabbed)")
	return nil
}

// Enable re-enables the touchpad by releasing the exclusive grab.
func (c *Controller) Enable() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.device == nil {
		return fmt.Errorf("touchpad device not open")
	}

	if !c.grabbed {
		return nil // Already enabled
	}

	if err := c.device.Ungrab(); err != nil {
		return fmt.Errorf("failed to ungrab touchpad: %w", err)
	}

	c.grabbed = false
	c.logger.Debug().Msg("Touchpad enabled (ungrabbed)")
	return nil
}

// IsDisabled returns whether the touchpad is currently disabled.
func (c *Controller) IsDisabled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.grabbed
}

// IsOpen returns whether the device is currently open.
func (c *Controller) IsOpen() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.device != nil
}

// DevicePath returns the touchpad device path.
func (c *Controller) DevicePath() string {
	return c.devicePath
}

// Stop stops the controller and releases the touchpad.
// Implements the component interface for graceful shutdown.
func (c *Controller) Stop() error {
	return c.Close()
}
