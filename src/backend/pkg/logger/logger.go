package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the global logger instance
var Logger zerolog.Logger

// Init initializes the logger with the given service name and log level
func Init(serviceName, level string) {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.DebugLevel
	}

	// Configure zerolog
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	// Create logger with console writer for development
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}

	Logger = zerolog.New(output).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()
}

// NewLogger creates a new logger with additional context
func NewLogger(component string) zerolog.Logger {
	return Logger.With().Str("component", component).Logger()
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

// Info logs an info message
func Info(msg string) {
	Logger.Info().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

// Error logs an error message
func Error(err error, msg string) {
	Logger.Error().Err(err).Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(err error, msg string) {
	Logger.Fatal().Err(err).Msg(msg)
}
