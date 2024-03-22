package api

import (
	"go.uber.org/zap"
)

// Handlers struct to hold dependencies for API handlers
type Handlers struct {
	Logger *zap.Logger
}

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(logger *zap.Logger) *Handlers {
	logger.Debug("Constructing new handlers")
	return &Handlers{
		Logger: logger,
	}
}
