package events

import (
    "sync"
    "github.com/rs/zerolog"
)

// SystemEventBus implements a broadcast pattern for system events.
type SystemEventBus struct {
    mu          sync.RWMutex
    subscribers []chan SystemEvent
    logger      zerolog.Logger
}

func NewSystemEventBus(logger zerolog.Logger) *SystemEventBus {
    return &SystemEventBus{
        subscribers: make([]chan SystemEvent, 0),
        logger:      logger,
    }
}

func (b *SystemEventBus) Publish(event SystemEvent) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    b.logger.Debug().Str("event", event.String()).Int("subscribers", len(b.subscribers)).Msg("Publishing SystemEvent")
    for i, sub := range b.subscribers {
        select {
        case sub <- event:
        default:
            b.logger.Warn().Str("event", event.String()).Int("subscriber", i).Msg("SystemEventBus subscriber buffer full, dropping event")
        }
    }
}

func (b *SystemEventBus) Subscribe() <-chan SystemEvent {
    ch := make(chan SystemEvent, 100)
    b.mu.Lock()
    b.subscribers = append(b.subscribers, ch)
    b.mu.Unlock()
    b.logger.Debug().Int("total_subscribers", len(b.subscribers)).Msg("New subscriber")
    return ch
}

func (b *SystemEventBus) Close() {
    b.mu.Lock()
    defer b.mu.Unlock()
    for _, sub := range b.subscribers {
        close(sub)
    }
    b.subscribers = nil
}
