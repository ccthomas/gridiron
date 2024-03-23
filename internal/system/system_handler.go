package system

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(logger *zap.Logger, db *sql.DB) *SystemHandlers {
	logger.Debug("Constructing new system handlers")
	return &SystemHandlers{
		Logger: logger,
		DB:     db,
	}
}

func (h *SystemHandlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Health Handler hit.")

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

func (h *SystemHandlers) DatabaseHealthHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Database Health Handler hit.")

	err := h.DB.Ping()
	if err != nil {
		http.Error(w, "Failed to ping the database", http.StatusInternalServerError)
	}

	// Create a HealthMessage instance
	message := HealthMessage{
		Message:   "Gridiron has a healthy connection to the database.",
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
	h.Logger.Debug("Database Health Handler completed.")
}
