package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Get() *zap.Logger {
	if Logger == nil {
		Logger, _ = zap.NewDevelopment()
	}

	return Logger
}
