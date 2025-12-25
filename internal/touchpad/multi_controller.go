package touchpad

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

// MultiController manages multiple touchpad controllers.
// The Zenbook Duo has two screens, each with its own touchpad.
// This controller ensures all touchpads are grabbed/ungrabbed together.
type MultiController struct {
	controllers []*Controller
	mu          sync.Mutex
	logger      zerolog.Logger
}

// NewMultiController creates a new multi-touchpad controller.
func NewMultiController(devices []*DeviceInfo, logger zerolog.Logger) *MultiController {
	controllers := make([]*Controller, 0, len(devices))
	for _, dev := range devices {
		controllers = append(controllers, NewController(dev.Path, logger))
	}
	return &MultiController{
		controllers: controllers,
		logger:      logger.With().Str("component", "multi_touchpad_ctrl").Logger(),
	}
}

// Open opens all touchpad devices for control.
func (m *MultiController) Open() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ctrl := range m.controllers {
		if err := ctrl.Open(); err != nil {
			// Close any already opened devices
			for _, c := range m.controllers {
				c.Close()
			}
			return fmt.Errorf("failed to open touchpad: %w", err)
		}
	}

	m.logger.Info().Int("count", len(m.controllers)).Msg("All touchpad controllers opened")
	return nil
}

// Close closes all touchpad devices.
func (m *MultiController) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for _, ctrl := range m.controllers {
		if err := ctrl.Close(); err != nil {
			lastErr = err
			m.logger.Warn().Err(err).Msg("Failed to close touchpad controller")
		}
	}

	m.logger.Info().Msg("All touchpad controllers closed")
	return lastErr
}

// Disable disables all touchpads by grabbing them.
func (m *MultiController) Disable() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ctrl := range m.controllers {
		if err := ctrl.Disable(); err != nil {
			return fmt.Errorf("failed to disable touchpad: %w", err)
		}
	}
	return nil
}

// Enable enables all touchpads by releasing the grab.
func (m *MultiController) Enable() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ctrl := range m.controllers {
		if err := ctrl.Enable(); err != nil {
			return fmt.Errorf("failed to enable touchpad: %w", err)
		}
	}
	return nil
}

// IsDisabled returns whether ALL touchpads are currently disabled.
func (m *MultiController) IsDisabled() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ctrl := range m.controllers {
		if !ctrl.IsDisabled() {
			return false
		}
	}
	return len(m.controllers) > 0
}

// Stop stops the controller and releases all touchpads.
func (m *MultiController) Stop() error {
	return m.Close()
}

// DeviceCount returns the number of touchpad devices being controlled.
func (m *MultiController) DeviceCount() int {
	return len(m.controllers)
}
