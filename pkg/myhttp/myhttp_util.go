package myhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ccthomas/gridiron/pkg/logger"
	"go.uber.org/zap"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	logger.Get().Debug("Write Error.")

	logger.Get().Debug("Write api error.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(&ApiError{
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	if err != nil {
		logger.Get().Error("Error occurred trying to construct api error.", zap.String("Message", message))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
