package system

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/myhttp"
)

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(db *sql.DB) *SystemHandlers {
	logger.Get().Debug("Constructing new system handlers")
	return &SystemHandlers{
		DB: db,
	}
}

func (h *SystemHandlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Health Handler hit.")

	// Create a HealthMessage instance
	message := HealthMessage{
		Message:   "Gridiron Service is Healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	logger.Get().Debug("JSON encode message.")
	response, err := json.Marshal(message)
	if err != nil {
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Write the response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	logger.Get().Debug("Health Handler completed.")
}

func (h *SystemHandlers) DatabaseHealthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Database Health Handler hit.")

	err := h.DB.Ping()
	if err != nil {
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	// Create a HealthMessage instance
	message := HealthMessage{
		Message:   "Gridiron has a healthy connection to the database.",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	logger.Get().Debug("JSON encode message.")
	response, err := json.Marshal(message)
	if err != nil {
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Write the response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	logger.Get().Debug("Database Health Handler completed.")
}
