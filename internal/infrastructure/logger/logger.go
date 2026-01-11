package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

type StandardLogger struct {
	w io.Writer
}

func New(w io.Writer) *StandardLogger {
	return &StandardLogger{w: w}
}

func NewFileLogger(path string) (*StandardLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &StandardLogger{w: f}, nil
}

// StderrLogger returns a logger that writes to Stderr (safe for MCP)
func NewStderrLogger() *StandardLogger {
	return &StandardLogger{w: os.Stderr}
}

func (l *StandardLogger) log(level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(l.w, "[%s] [%s] %s\n", timestamp, level, msg)
}

func (l *StandardLogger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

func (l *StandardLogger) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

func (l *StandardLogger) Debug(format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}

// NoOpLogger for tests or when logging is disabled
type NoOpLogger struct{}

func (l *NoOpLogger) Info(format string, args ...interface{})  {}
func (l *NoOpLogger) Error(format string, args ...interface{}) {}
func (l *NoOpLogger) Debug(format string, args ...interface{}) {}
