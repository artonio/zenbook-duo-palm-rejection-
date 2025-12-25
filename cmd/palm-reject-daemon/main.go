// Package main is the entry point for the palm‑rejection daemon.
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/spf13/cobra"

    "github.com/artonio/zenbook-duo-palm-rejection/internal/consumer"
    "github.com/artonio/zenbook-duo-palm-rejection/internal/events"
    "github.com/artonio/zenbook-duo-palm-rejection/internal/pipe"
    "github.com/artonio/zenbook-duo-palm-rejection/internal/touchpad"
    "github.com/artonio/zenbook-duo-palm-rejection/pkg/logging"
)

const version = "0.1.0"

func main() {
    rootCmd := &cobra.Command{
        Use:     "palm-reject-daemon",
        Short:   "Daemon that disables the touchpad while typing",
        Version: version,
    }

    var timeout time.Duration

    runCmd := &cobra.Command{
        Use:   "run",
        Short: "Run the daemon",
        RunE:  runDaemon,
    }

    runCmd.Flags().DurationVar(&timeout, "timeout", 0, "Auto-stop after duration (e.g., 10s, 1m) for safe testing")

    rootCmd.AddCommand(runCmd)

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func runDaemon(cmd *cobra.Command, _ []string) error {
    // Get timeout flag
    timeout, _ := cmd.Flags().GetDuration("timeout")
    // Logging
    logLevel := logging.GetLogLevelFromEnv()
    logger := logging.SetupLogger(logLevel)

    logger.Info().
        Str("version", version).
        Msg("Starting palm‑rejection daemon")

    if timeout > 0 {
        logger.Warn().Dur("timeout", timeout).Msg("Running with timeout - will auto-stop")
    }

    // Event bus
    systemEventBus := events.NewSystemEventBus(logger)

    // Context for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Components to stop on shutdown
    type component interface {
        Stop() error
    }
    var components []component

    // Pipe receiver
    pipeReceiver := pipe.NewReceiver(
        pipe.DefaultPipePath,
        systemEventBus,
        logger,
    )
    if err := pipeReceiver.Start(ctx); err != nil {
        logger.Warn().Err(err).Msg("pipe receiver failed to start")
    } else {
        components = append(components, pipeReceiver)
    }

    // Touchpad discovery
    devs, err := touchpad.FindAllTouchpadDevices(logger)
    if err != nil {
        logger.Warn().Err(err).Msg("no touchpad devices found; daemon will exit")
        return err
    }

    // Keyboard discovery
    keyInfo, err := touchpad.FindKeyboardDevice(logger)
    if err != nil {
        logger.Warn().Err(err).Msg("keyboard device not found; daemon will exit")
        return err
    }

    // Touchpad controller
    touchpadCtrl := touchpad.NewMultiController(devs, logger)
    if err := touchpadCtrl.Open(); err != nil {
        logger.Error().Err(err).Msg("failed to open touchpads")
        return err
    }
    components = append(components, touchpadCtrl)

    // Typing detection consumer
    cooldown := 300 * time.Millisecond
    typingConsumer := consumer.NewTypingDetectionConsumer(
        nil,            // will be set after monitor is created
        touchpadCtrl,
        systemEventBus,
        cooldown,
        logger,
    )

    // Keyboard monitor
    keyboardMonitor := touchpad.NewKeyboardMonitor(
        keyInfo.Path,
        typingConsumer.OnKeyPress,
        logger,
    )
    if err := keyboardMonitor.Start(ctx); err != nil {
        logger.Error().Err(err).Msg("keyboard monitor failed to start")
        return err
    }
    components = append(components, keyboardMonitor)

    // Start consumer
    if err := typingConsumer.Start(ctx); err != nil {
        logger.Error().Err(err).Msg("typing consumer failed to start")
        return err
    }
    components = append(components, typingConsumer)

    logger.Info().
        Strs("touchpads", getPaths(devs)).
        Str("keyboard", keyInfo.Path).
        Dur("cooldown", cooldown).
        Int("touchpad_count", len(devs)).
        Msg("Palm rejection active")

    // Wait for shutdown
    var sig os.Signal
    if timeout > 0 {
        select {
        case sig = <-sigChan:
            logger.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
        case <-time.After(timeout):
            logger.Info().Dur("timeout", timeout).Msg("Timeout reached, shutting down")
        }
    } else {
        sig = <-sigChan
        logger.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
    }

    // Stop all components
    for _, c := range components {
        if err := c.Stop(); err != nil {
            logger.Warn().Err(err).Msg("failed to stop component")
        }
    }

    systemEventBus.Close()

    logger.Info().Msg("daemon stopped")
    return nil
}

func getPaths(devs []*touchpad.DeviceInfo) []string {
    var paths []string
    for _, d := range devs {
        paths = append(paths, d.Path)
    }
    return paths
}
