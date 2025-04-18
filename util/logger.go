package util

import (
	"log"
	"os"
)

// Logger is a simple structured logger
type Logger struct {
	*log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{Logger: log.New(os.Stdout, "[gk-migrations] ", log.LstdFlags)}
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	l.Println("INFO: " + msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.Println("ERROR: " + msg)
}
