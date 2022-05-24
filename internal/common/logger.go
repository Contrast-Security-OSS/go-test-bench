package common

import (
	"log"
)

// Logger contains selected methods of testing.TB
type Logger interface {
	Logf(f string, va ...interface{})
	Errorf(f string, va ...interface{})
	Fatalf(f string, va ...interface{})
}

// LogWrapper wraps log.Logger to work with the above logger
// interface. Errorf() and Fatalf() are equivalent.
type LogWrapper struct {
	l *log.Logger
}

// NewLogWrapper creates a wrapper around log.Logger conforming to Logger.
func NewLogWrapper(l *log.Logger) Logger { return &LogWrapper{l: l} }

// Logf implements Logger.
func (w *LogWrapper) Logf(f string, va ...interface{}) { w.l.Printf(f, va...) }

// Errorf implements Logger.
func (w *LogWrapper) Errorf(f string, va ...interface{}) { w.l.Fatalf(f, va...) }

// Fatalf implements Logger.
func (w *LogWrapper) Fatalf(f string, va ...interface{}) { w.l.Fatalf(f, va...) }
