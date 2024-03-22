package logger

// This test does not mock zap.
// Instead it verifies our Zap Logger works properly

import (
	"testing"
)

func TestGetLoggerNotNil(t *testing.T) {
	// When
	logger := Get()

	// Then
	if logger == nil {
		t.Error("Expected logger to be initialized, got nil")
	}
}
