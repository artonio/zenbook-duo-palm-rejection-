package pipe

import (
    "bufio"
    "context"
    "os"
    "strings"
    "syscall"
    "time"

    "github.com/rs/zerolog"
    "github.com/artonio/zenbook-duo-palm-rejection/internal/events"
)

const DefaultPipePath = "/tmp/zenbook-duo-daemon.pipe"

type Receiver struct {
    ctx            context.Context
    cancel         context.CancelFunc
    path           string
    systemEventBus *events.SystemEventBus
    logger         zerolog.Logger
}

func NewReceiver(path string, bus *events.SystemEventBus, logger zerolog.Logger) *Receiver {
    if path == "" {
        path = DefaultPipePath
    }
    return &Receiver{
        path:           path,
        systemEventBus: bus,
        logger:         logger.With().Str("component", "pipe_receiver").Logger(),
    }
}

func (r *Receiver) Start(ctx context.Context) error {
    r.ctx, r.cancel = context.WithCancel(ctx)
    os.Remove(r.path)
    if err := syscall.Mkfifo(r.path, 0622); err != nil {
        return err
    }
    go r.readLoop()
    r.logger.Info().Str("path", r.path).Msg("Pipe receiver started")
    return nil
}

func (r *Receiver) Stop() error {
    if r.cancel != nil {
        r.cancel()
    }
    os.Remove(r.path)
    return nil
}

func (r *Receiver) readLoop() {
    for {
        select {
        case <-r.ctx.Done():
            return
        default:
            r.readOnce()
        }
    }
}

func (r *Receiver) readOnce() {
    file, err := os.OpenFile(r.path, os.O_RDONLY, os.ModeNamedPipe)
    if err != nil {
        if r.ctx.Err() == nil {
            r.logger.Error().Err(err).Msg("failed to open pipe")
            time.Sleep(time.Second)
        }
        return
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        select {
        case <-r.ctx.Done():
            return
        default:
            r.handleCommand(scanner.Text())
        }
    }
}

func (r *Receiver) handleCommand(cmd string) {
    cmd = strings.TrimSpace(strings.ToLower(cmd))
    if cmd == "" {
        return
    }
    r.logger.Info().Str("command", cmd).Msg("pipe command received")
    switch cmd {
    case "touchpad_disable":
        r.systemEventBus.Publish(events.TouchpadDisable)
    case "touchpad_enable":
        r.systemEventBus.Publish(events.TouchpadEnable)
    case "touchpad_toggle":
        r.systemEventBus.Publish(events.TouchpadToggle)
    default:
        r.logger.Warn().Str("command", cmd).Msg("unknown pipe command")
    }
}
