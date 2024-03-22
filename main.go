package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ccthomas/gridiron/api"
	gridironLogger "github.com/ccthomas/gridiron/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger := gridironLogger.Get()
	defer logger.Sync()

	logger.Info("Starting Gridiron...")

	handler := api.NewHandlers(logger)

	logger.Debug("Construct router.")
	r := mux.NewRouter()

	logger.Debug("Route paths to handlers.")
	r.HandleFunc("/health", handler.HealthHandler).Methods("GET")

	logger.Debug("Handle router with http.")
	http.Handle("/", r)

	serverPort := os.Getenv("SERVER_PORT")

	logger.Info("Server starting...", zap.String("port", serverPort))
	err := http.ListenAndServe(fmt.Sprintf(":%v", serverPort), nil)
	if err != nil {
		logger.Fatal("Failed to listen and serve app.", zap.Error(err))
	}
}
