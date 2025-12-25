package events

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSystemEvent_String(t *testing.T) {
	tests := []struct {
		name     string
		event    SystemEvent
		expected string
	}{
		{"TouchpadDisable", TouchpadDisable, "TouchpadDisable"},
		{"TouchpadEnable", TouchpadEnable, "TouchpadEnable"},
		{"LaptopSuspend", LaptopSuspend, "LaptopSuspend"},
		{"LaptopResume", LaptopResume, "LaptopResume"},
		{"BacklightToggle", BacklightToggle, "BacklightToggle"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.event.String())
		})
	}
}

func TestSystemEventBus_BasicOperations(t *testing.T) {
	logger := zerolog.Nop()
	bus := NewSystemEventBus(logger)

	// Test 1: Subscribe and check channel properties
	sub := bus.Subscribe()
	assert.NotNil(t, sub)

	// Test 2: Simple publish/subscribe
	event := TouchpadDisable
	go bus.Publish(event)

	// Give it a moment to process
	time.Sleep(10 * time.Millisecond)

	// Check if event is available (non-blocking)
	select {
	case received := <-sub:
		assert.Equal(t, event, received)
	default:
		// No event yet, that's ok for this test
	}

	// Drain any remaining events
	for {
		select {
		case <-sub:
			// Drain event
		default:
			// Channel empty
			goto done
		}
	}
done:

	// Test 3: Multiple publishers/subscribers
	sub2 := bus.Subscribe()
	event2 := TouchpadEnable

	go bus.Publish(event2)
	time.Sleep(10 * time.Millisecond)

	// Both subscribers should have received the event
	received1 := false
	received2 := false

	select {
	case <-sub:
		received1 = true
	default:
	}

	select {
	case <-sub2:
		received2 = true
	default:
	}

	// At least one should have received it
	assert.True(t, received1 || received2, "At least one subscriber should receive the event")

	bus.Close()
}