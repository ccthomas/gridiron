package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ccthomas/gridiron/api"
	"github.com/ccthomas/gridiron/internal/useracc"
	"github.com/ccthomas/gridiron/pkg/database"
	gridironLogger "github.com/ccthomas/gridiron/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger := gridironLogger.Get()
	defer logger.Sync()

	logger.Info("Starting Gridiron...")

	logger.Debug("Connect to database.")
	dm := database.DatabaseManager{
		Logger: logger,
	}

	db := dm.ConnectPostgres()
	defer db.Close()

	userRepo := &useracc.UserAccountRepositoryImpl{
		DB:     db,
		Logger: logger,
	}

	logger.Debug("Construct handlers.")
	handler := api.NewHandlers(db, logger, userRepo)

	logger.Debug("Construct router.")
	r := mux.NewRouter()

	logger.Debug("Route paths to handlers.")
	handler.RouteApis(r)

	logger.Debug("Handle router with http.")
	http.Handle("/", r)

	serverPort := os.Getenv("SERVER_PORT")

	logger.Info("Server starting...", zap.String("port", serverPort))
	err := http.ListenAndServe(fmt.Sprintf(":%v", serverPort), nil)
	if err != nil {
		logger.Fatal("Failed to listen and serve app.", zap.Error(err))
	}
}
