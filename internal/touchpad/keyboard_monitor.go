package touchpad

import (
	"context"
	"fmt"

	evdev "github.com/holoplot/go-evdev"
	"github.com/rs/zerolog"
)

// KeyboardMonitor monitors a keyboard evdev device for typing activity.
// This detects regular keypresses (a-z, numbers, etc.) - not the special Fn keys
// that come through hidraw.
type KeyboardMonitor struct {
	ctx        context.Context
	cancel     context.CancelFunc
	devicePath string
	device     *evdev.InputDevice
	onKeyPress func() // Callback when any key is pressed
	logger     zerolog.Logger
}

// NewKeyboardMonitor creates a new keyboard monitor.
func NewKeyboardMonitor(devicePath string, onKeyPress func(), logger zerolog.Logger) *KeyboardMonitor {
	return &KeyboardMonitor{
		devicePath: devicePath,
		onKeyPress: onKeyPress,
		logger:     logger.With().Str("component", "kb_monitor").Str("device", devicePath).Logger(),
	}
}

// Start starts the keyboard monitor.
func (m *KeyboardMonitor) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)

	// Open the evdev device
	dev, err := evdev.Open(m.devicePath)
	if err != nil {
		return fmt.Errorf("failed to open keyboard device %s: %w", m.devicePath, err)
	}
	m.device = dev

	name, err := m.device.Name()
	if err != nil {
		m.device.Close()
		return fmt.Errorf("failed to get device name: %w", err)
	}

	m.logger.Info().Str("name", name).Msg("Keyboard monitor started")

	// Start reading in a goroutine
	go m.readLoop()

	return nil
}

// Stop stops the keyboard monitor.
func (m *KeyboardMonitor) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}

	if m.device != nil {
		m.device.Close()
		m.device = nil
	}

	m.logger.Info().Msg("Keyboard monitor stopped")
	return nil
}

// readLoop reads events from the keyboard evdev device.
func (m *KeyboardMonitor) readLoop() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			// Read one event (blocking)
			ev, err := m.device.ReadOne()
			if err != nil {
				if m.ctx.Err() != nil {
					return // Context cancelled
				}
				m.logger.Error().Err(err).Msg("Keyboard read error")
				return
			}

			// We're looking for EV_KEY events with value=1 (key press)
			// value=0 is key release, value=2 is key repeat
			if ev.Type == evdev.EV_KEY && ev.Value == 1 {
				m.logger.Debug().
					Uint16("code", uint16(ev.Code)).
					Msg("Key press detected")

				if m.onKeyPress != nil {
					m.onKeyPress()
				}
			}
		}
	}
}

// DevicePath returns the keyboard device path.
func (m *KeyboardMonitor) DevicePath() string {
	return m.devicePath
}
