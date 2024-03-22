package system

import "go.uber.org/zap"

type HealthMessage struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type SystemHandlers struct {
	Logger *zap.Logger
}
