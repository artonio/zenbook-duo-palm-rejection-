package consumer

import (
    "context"
    "sync"
    "time"

    "github.com/rs/zerolog"

    "github.com/artonio/zenbook-duo-palm-rejection/internal/events"
    "github.com/artonio/zenbook-duo-palm-rejection/internal/touchpad"
)

// TypingDetectionConsumer disables the touchpad while typing to prevent accidental cursor movement (palm rejection).
type TypingDetectionConsumer struct {
    ctx             context.Context
    cancel          context.CancelFunc
    keyboardMonitor *touchpad.KeyboardMonitor
    touchpadCtrl    touchpad.TouchpadController
    systemEventBus  *events.SystemEventBus
    cooldown        time.Duration
    logger          zerolog.Logger

    mu           sync.Mutex
    lastKeyPress time.Time
    timer        *time.Timer
    isDisabled   bool
}

// NewTypingDetectionConsumer creates a new typing detection consumer.
func NewTypingDetectionConsumer(
    keyboardMonitor *touchpad.KeyboardMonitor,
    touchpadCtrl touchpad.TouchpadController,
    systemEventBus *events.SystemEventBus,
    cooldown time.Duration,
    logger zerolog.Logger,
) *TypingDetectionConsumer {
    return &TypingDetectionConsumer{
        keyboardMonitor: keyboardMonitor,
        touchpadCtrl:    touchpadCtrl,
        systemEventBus:  systemEventBus,
        cooldown:        cooldown,
        logger:          logger.With().Str("component", "typing_detection").Logger(),
    }
}

// Start starts the typing detection consumer.
func (c *TypingDetectionConsumer) Start(ctx context.Context) error {
    c.ctx, c.cancel = context.WithCancel(ctx)

    // Subscribe to system events for suspend/resume
    go c.systemEventLoop()

    c.logger.Info().
        Dur("cooldown", c.cooldown).
        Msg("Typing detection consumer started")

    return nil
}

// Stop stops the typing detection consumer.
func (c *TypingDetectionConsumer) Stop() error {
    if c.cancel != nil {
        c.cancel()
    }

    c.mu.Lock()
    defer c.mu.Unlock()

    // Stop the cooldown timer
    if c.timer != nil {
        c.timer.Stop()
        c.timer = nil
    }

    // Ensure touchpad is enabled on shutdown
    if c.isDisabled {
        if err := c.touchpadCtrl.Enable(); err != nil {
            c.logger.Warn().Err(err).Msg("Failed to enable touchpad during shutdown")
        }
        c.isDisabled = false
    }

    c.logger.Info().Msg("Typing detection consumer stopped")
    return nil
}

// OnKeyPress is called when a key is pressed on the keyboard.
func (c *TypingDetectionConsumer) OnKeyPress() {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.lastKeyPress = time.Now()

    // Disable touchpad if not already disabled
    if !c.isDisabled {
        if err := c.touchpadCtrl.Disable(); err != nil {
            c.logger.Error().Err(err).Msg("Failed to disable touchpad")
            return
        }
        c.isDisabled = true
        c.logger.Debug().Msg("Touchpad disabled (typing detected)")
    }

    // Reset or start the cooldown timer
    if c.timer != nil {
        c.timer.Stop()
    }
    c.timer = time.AfterFunc(c.cooldown, c.onCooldownExpired)
}

// onCooldownExpired is called when the cooldown timer expires.
func (c *TypingDetectionConsumer) onCooldownExpired() {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Check if we should re-enable (no recent keypresses)
    if time.Since(c.lastKeyPress) >= c.cooldown && c.isDisabled {
        if err := c.touchpadCtrl.Enable(); err != nil {
            c.logger.Error().Err(err).Msg("Failed to enable touchpad after cooldown")
            return
        }
        c.isDisabled = false
        c.logger.Debug().Msg("Touchpad enabled (cooldown expired)")
    }
}

// systemEventLoop handles system events (suspend, resume).
func (c *TypingDetectionConsumer) systemEventLoop() {
    sub := c.systemEventBus.Subscribe()

    for {
        select {
        case <-c.ctx.Done():
            return
        case event := <-sub:
            c.handleSystemEvent(event)
        }
    }
}

// handleSystemEvent processes system events.
func (c *TypingDetectionConsumer) handleSystemEvent(event events.SystemEvent) {
    switch event {
    case events.LaptopSuspend:
        c.mu.Lock()
        // Ensure touchpad is enabled before suspend
        if c.isDisabled {
            if err := c.touchpadCtrl.Enable(); err != nil {
                c.logger.Warn().Err(err).Msg("Failed to enable touchpad for suspend")
            }
            c.isDisabled = false
        }
        // Stop the timer
        if c.timer != nil {
            c.timer.Stop()
            c.timer = nil
        }
        c.mu.Unlock()
        c.logger.Debug().Msg("Touchpad enabled for suspend")

    case events.LaptopResume:
        // Nothing special needed - touchpad should be enabled
        c.logger.Debug().Msg("Laptop resumed, touchpad control ready")

    case events.TouchpadDisable:
        c.mu.Lock()
        if !c.isDisabled {
            if err := c.touchpadCtrl.Disable(); err != nil {
                c.logger.Error().Err(err).Msg("Failed to disable touchpad via pipe command")
            } else {
                c.isDisabled = true
                c.logger.Info().Msg("Touchpad disabled via pipe command")
            }
        }
        // Stop any cooldown timer since this is a manual action
        if c.timer != nil {
            c.timer.Stop()
            c.timer = nil
        }
        c.mu.Unlock()

    case events.TouchpadEnable:
        c.mu.Lock()
        if c.isDisabled {
            if err := c.touchpadCtrl.Enable(); err != nil {
                c.logger.Error().Err(err).Msg("Failed to enable touchpad via pipe command")
            } else {
                c.isDisabled = false
                c.logger.Info().Msg("Touchpad enabled via pipe command")
            }
        }
        // Stop any cooldown timer
        if c.timer != nil {
            c.timer.Stop()
            c.timer = nil
        }
        c.mu.Unlock()

    case events.TouchpadToggle:
        c.mu.Lock()
        if c.isDisabled {
            if err := c.touchpadCtrl.Enable(); err != nil {
                c.logger.Error().Err(err).Msg("Failed to enable touchpad via pipe command")
            } else {
                c.isDisabled = false
                c.logger.Info().Msg("Touchpad enabled via pipe command (toggle)")
            }
        } else {
            if err := c.touchpadCtrl.Disable(); err != nil {
                c.logger.Error().Err(err).Msg("Failed to disable touchpad via pipe command")
            } else {
                c.isDisabled = true
                c.logger.Info().Msg("Touchpad disabled via pipe command (toggle)")
            }
        }
        // Stop any cooldown timer since this is a manual action
        if c.timer != nil {
            c.timer.Stop()
            c.timer = nil
        }
        c.mu.Unlock()
    }
}

// IsDisabled returns whether the touchpad is currently disabled due to typing.
func (c *TypingDetectionConsumer) IsDisabled() bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.isDisabled
}
