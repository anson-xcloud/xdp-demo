package xlog

// Default 默认的logger
var Default = NewZap()

// Logger logger需要实现的接口
type Logger interface {
	With(fields ...interface{}) Logger

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Debugf shorcut for default.Debuf
func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}

// Infof shorcut for default.Infof
func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

// Warnf shorcut for default.Warnf
func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

// Errorf shorcut for default.Errorf
func Errorf(format string, args ...interface{}) {
	Default.Errorf(format, args...)
}
