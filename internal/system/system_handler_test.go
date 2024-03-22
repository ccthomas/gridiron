package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestHealthHandler(t *testing.T) {
	// Given
	logger := zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel))

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	handlers := &SystemHandlers{
		Logger: logger,
	}

	// When
	startTime := time.Now().UTC()
	handlers.HealthHandler(rec, req)
	endTime := time.Now().UTC()

	// Then
	assert.Equal(t, http.StatusOK, rec.Code, "http status is not 200.")
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"), "content type header is not application/json.")

	var healthMsg HealthMessage
	err := json.Unmarshal(rec.Body.Bytes(), &healthMsg)
	assert.NoError(t, err, "There was an error parsing health message.")
	assert.Equal(t, "Gridiron Service is Healthy", healthMsg.Message, "health message is incorrect.")

	parsedTime, err := time.Parse(time.RFC3339Nano, healthMsg.Timestamp)
	startTime = startTime.Truncate(time.Second)
	endTime = endTime.Truncate(time.Second)

	assert.NoError(t, err)
	assert.True(t, parsedTime.Equal(startTime) || parsedTime.After(startTime), "timestamp is not after or equal to the start of the test.")
	assert.True(t, parsedTime.Equal(endTime) || parsedTime.Before(endTime), "timestamp is not before or equal to the end of the test.")
}
