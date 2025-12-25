package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/artonio/zenbook-duo-palm-rejection/internal/events"
)

// MockTouchpadController is a mock implementation of touchpad.TouchpadController
type MockTouchpadController struct {
	mock.Mock
	disabled bool
}

func (m *MockTouchpadController) Disable() error {
	args := m.Called()
	if args.Error(0) == nil {
		m.disabled = true
	}
	return args.Error(0)
}

func (m *MockTouchpadController) Enable() error {
	args := m.Called()
	if args.Error(0) == nil {
		m.disabled = false
	}
	return args.Error(0)
}

func (m *MockTouchpadController) IsDisabled() bool {
	return m.disabled
}

func (m *MockTouchpadController) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func TestTypingDetectionConsumer_BasicOperations(t *testing.T) {
	// Create mocks
	mockCtrl := new(MockTouchpadController)
	eventBus := events.NewSystemEventBus(zerolog.Nop())
	logger := zerolog.Nop()

	// Create consumer with 100ms cooldown for faster testing
	consumer := NewTypingDetectionConsumer(
		nil, // keyboard monitor will be set after consumer creation
		mockCtrl,
		eventBus,
		100*time.Millisecond,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test 1: Initial state - touchpad should be enabled
	assert.False(t, consumer.isDisabled)

	// Test 2: Single keypress - should disable touchpad
	mockCtrl.On("Disable").Return(nil).Once()
	consumer.OnKeyPress()
	assert.True(t, consumer.isDisabled)

	// Test 3: Multiple keypresses during cooldown - should not re-disable
	consumer.OnKeyPress()
	assert.True(t, consumer.isDisabled) // Still disabled

	// Set up expectation for Stop() which will call Enable()
	mockCtrl.On("Enable").Return(nil).Once()

	err = consumer.Stop()
	assert.NoError(t, err)
}

func TestTypingDetectionConsumer_DisableError(t *testing.T) {
	mockCtrl := new(MockTouchpadController)
	eventBus := events.NewSystemEventBus(zerolog.Nop())
	logger := zerolog.Nop()

	consumer := NewTypingDetectionConsumer(
		nil,
		mockCtrl,
		eventBus,
		300*time.Millisecond,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := consumer.Start(ctx)
	assert.NoError(t, err)

	// Test: Error when disabling touchpad
	mockCtrl.On("Disable").Return(assert.AnError).Once()
	consumer.OnKeyPress()
	assert.False(t, consumer.isDisabled) // Should remain false on error
	mockCtrl.AssertExpectations(t)

	// Set up expectation for Stop() which will call Enable()
	mockCtrl.On("Enable").Return(nil).Once()

	err = consumer.Stop()
	assert.NoError(t, err)
}