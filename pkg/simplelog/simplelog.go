package simplelog

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	levelNames    = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	defaultLogger = NewLogger(INFO)
)

// NewLogger creates a new Logger instance with the given log level.
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

// SetLevel sets the log level for the logger.
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// logMessage logs a message if the log level is sufficient.
func (l *Logger) logMessage(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	l.logger.Printf("[%s] %s", levelNames[level], formatMessage(format, args...))
}

// formatMessage formats a log message with the given arguments.
func formatMessage(format string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

// Debug logs a debug-level message.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logMessage(DEBUG, format, args...)
}

// Info logs an info-level message.
func (l *Logger) Info(format string, args ...interface{}) {
	l.logMessage(INFO, format, args...)
}

// Warn logs a warning-level message.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.logMessage(WARN, format, args...)
}

// Error logs an error-level message.
func (l *Logger) Error(format string, args ...interface{}) {
	l.logMessage(ERROR, format, args...)
}

// Fatal logs a fatal-level message and exits the program.
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logMessage(FATAL, format, args...)
	os.Exit(1)
}

// Debug logs a debug-level message using the default logger.
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info-level message using the default logger.
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning-level message using the default logger.
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error-level message using the default logger.
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal logs a fatal-level message using the default logger and exits the program.
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// SetLogLevel sets the log level for the default logger.
func SetLogLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}
