package xlog

import (
	"fmt"
	"log"
	"strings"
)

type stdLog struct {
	addon string
}

func NewStd() Logger {
	return &stdLog{}
}

func (l *stdLog) With(fields ...interface{}) Logger {
	var sf []string
	for i := 0; i+1 < len(fields); i += 2 {
		s := fmt.Sprintf("{%v:%v}", fields[i], fields[i+1])
		sf = append(sf, s)
	}

	return &stdLog{addon: strings.Join(sf, " ")}
}

func (l *stdLog) log(format string, args ...interface{}) {
	if l.addon != "" {
		format += " %s"
		args = append(args, l.addon)
	}

	log.Printf(format, args...)
}

func (l *stdLog) Debugf(format string, args ...interface{}) {
	l.log(format, args...)
}

func (l *stdLog) Infof(format string, args ...interface{}) {
	l.log(format, args...)
}

func (l *stdLog) Warnf(format string, args ...interface{}) {
	l.log(format, args...)
}

func (l *stdLog) Errorf(format string, args ...interface{}) {
	l.log(format, args...)
}
