package logging

import (
    "os"
    "strings"
    "time"

    "github.com/rs/zerolog"
)

// SetupLogger creates and configures a logger based on the log level.
func SetupLogger(level string) zerolog.Logger {
    var logLevel zerolog.Level
    switch strings.ToLower(level) {
    case "trace":
        logLevel = zerolog.TraceLevel
    case "debug":
        logLevel = zerolog.DebugLevel
    case "info":
        logLevel = zerolog.InfoLevel
    case "warn", "warning":
        logLevel = zerolog.WarnLevel
    case "error":
        logLevel = zerolog.ErrorLevel
    case "fatal":
        logLevel = zerolog.FatalLevel
    default:
        logLevel = zerolog.InfoLevel
    }

    output := zerolog.ConsoleWriter{
        Out:        os.Stderr,
        TimeFormat: time.RFC3339,
    }

    return zerolog.New(output).
        Level(logLevel).
        With().
        Timestamp().
        Logger()
}

// SetupLoggerJSON creates a logger that outputs JSON.
func SetupLoggerJSON(level string) zerolog.Logger {
    var logLevel zerolog.Level
    switch strings.ToLower(level) {
    case "trace":
        logLevel = zerolog.TraceLevel
    case "debug":
        logLevel = zerolog.DebugLevel
    case "info":
        logLevel = zerolog.InfoLevel
    case "warn", "warning":
        logLevel = zerolog.WarnLevel
    case "error":
        logLevel = zerolog.ErrorLevel
    case "fatal":
        logLevel = zerolog.FatalLevel
    default:
        logLevel = zerolog.InfoLevel
    }

    return zerolog.New(os.Stderr).
        Level(logLevel).
        With().
        Timestamp().
        Logger()
}

// GetLogLevelFromEnv reads the LOG_LEVEL env var.
func GetLogLevelFromEnv() string {
    level := os.Getenv("LOG_LEVEL")
    if level == "" {
        return "info"
    }
    return level
}
