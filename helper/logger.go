package helper

import (
	"log"
	"os"
)

// Logger is a wrapper for logging with different levels.
type Logger struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// NewLogger initializes a new logger with standard output.
func NewLogger() *Logger {
	return &Logger{
		Info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warning: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
