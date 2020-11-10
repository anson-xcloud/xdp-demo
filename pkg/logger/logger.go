package logger

import (
	"fmt"
	"time"
)

var Default Logger

func init() {
	Default = new(fmtLogger)
}

// Logger xdp logger
type Logger interface {
	Debug(string, ...interface{})

	Info(string, ...interface{})

	Warn(string, ...interface{})

	Error(string, ...interface{})
}

func Info(format string, args ...interface{}) {
	Default.Info(format, args...)
}

func Error(format string, args ...interface{}) {
	Default.Error(format, args...)
}

type fmtLogger struct {
}

func (l *fmtLogger) output(lv, msg string) {
	tf := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Printf("%s %s: %s\n", tf, lv, msg)
}

func (l *fmtLogger) Debug(format string, args ...interface{}) {
	l.output("DEBUG", fmt.Sprintf(format, args...))
}

func (l *fmtLogger) Info(format string, args ...interface{}) {
	l.output(" INFO", fmt.Sprintf(format, args...))
}

func (l *fmtLogger) Warn(format string, args ...interface{}) {
	l.output(" WARN", fmt.Sprintf(format, args...))
}

func (l *fmtLogger) Error(msg string, args ...interface{}) {
	l.output("ERROR", fmt.Sprintf(msg, args...))
}
