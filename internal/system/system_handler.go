package system

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(logger *zap.Logger) *SystemHandlers {
	logger.Debug("Constructing new system handlers")
	return &SystemHandlers{
		Logger: logger,
	}
}

func (h *SystemHandlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Health Handler hit.")
	h.Logger.Debug("Request for Health Handler.", zap.Any("request", r))

	// Create a HealthMessage instance
	message := HealthMessage{
		Message:   "Gridiron Service is Healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	h.Logger.Debug("JSON encode message.")
	response, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.Logger.Debug("Write the response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	h.Logger.Debug("Health Handler completed.")
}
