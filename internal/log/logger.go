package log

import (
	"log"
	"os"
)

// New creates a configured standard logger.
func New(environment string) *log.Logger {
	logger := log.New(os.Stdout, "[goworld] ", log.LstdFlags|log.Lshortfile)
	_ = environment
	return logger
}
