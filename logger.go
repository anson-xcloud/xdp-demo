package xdp

import (
	"fmt"
	"time"
)

// Logger xdp logger
type Logger interface {
	Debug(string, ...interface{})

	Info(string, ...interface{})

	Error(string, ...interface{})
}

type fmtLogger struct {
}

func (l *fmtLogger) output(lv, msg string) {
	tf := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Printf("%s %s: %s\n", tf, lv, msg)
}

func (l *fmtLogger) Debug(msg string, args ...interface{}) {
	l.output("DEBUG", fmt.Sprintf(msg, args...))
}

func (l *fmtLogger) Info(msg string, args ...interface{}) {
	l.output(" INFO", fmt.Sprintf(msg, args...))
}

func (l *fmtLogger) Error(msg string, args ...interface{}) {
	l.output("ERROR", fmt.Sprintf(msg, args...))
}
