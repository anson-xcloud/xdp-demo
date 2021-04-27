package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	*zap.SugaredLogger
}

func NewZap() Logger {
	cfg := zap.NewDevelopmentConfig()
	// 设置颜色
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, _ := cfg.Build(zap.AddCallerSkip(0), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	return &zapLogger{SugaredLogger: sugar}
}

func (zl *zapLogger) With(fields ...interface{}) Logger {
	sugar := zl.SugaredLogger.With(fields...)
	return &zapLogger{SugaredLogger: sugar}
}
