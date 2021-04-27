package xlog

import (
	"fmt"
)

type consoleLogger struct {
}

func NewConsole() Logger {
	l := new(consoleLogger)
	return l
}

func (l *consoleLogger) With(fields ...interface{}) Logger {
	return NewConsole()
}

func (l *consoleLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf(format+"\r\n", args...)
}

func (l *consoleLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\r\n", args...)
}

func (l *consoleLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf(format+"\r\n", args...)
}

func (l *consoleLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format+"\r\n", args...)
}
