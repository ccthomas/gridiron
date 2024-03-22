package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Get() *zap.Logger {
	if Logger == nil {
		Logger, _ = zap.Config{
			Encoding:    "json",
			Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
			OutputPaths: []string{"stdout"},
		}.Build()
	}

	return Logger
}
