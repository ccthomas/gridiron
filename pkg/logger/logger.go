package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Get() *zap.Logger {
	logger, _ := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
	}.Build()

	return logger
}
