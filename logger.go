package xdp

import "fmt"

// Logger xdp logger
type Logger interface {
	Debug(string, ...interface{})

	Info(string, ...interface{})

	Error(string, ...interface{})
}

type fmtLogger struct {
}

func (l *fmtLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

func (l *fmtLogger) Info(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

func (l *fmtLogger) Error(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}
