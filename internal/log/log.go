// Package log provides a logger that can be disabled.
package log

import "log"

// Logger is an interface for logging.
type Logger interface {
	Printf(format string, v ...interface{})
}

// New returns a new logger.
func New(debug bool) Logger {
	if debug {
		return &logger{}
	}
	return &noOpLogger{}
}

type logger struct{}

func (l *logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

type noOpLogger struct{}

func (l *noOpLogger) Printf(format string, v ...interface{}) {}
