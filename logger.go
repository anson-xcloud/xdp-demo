package xdp

import "fmt"

type Logger interface {
	// Infof(string, ...interface{})

	Errorf(string, ...interface{})
}

type fmtLogger struct {
}

func (l *fmtLogger) Errorf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}
