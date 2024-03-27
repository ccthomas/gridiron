package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ccthomas/gridiron/api"
	"github.com/ccthomas/gridiron/internal/team"
	"github.com/ccthomas/gridiron/internal/tenant"
	"github.com/ccthomas/gridiron/internal/useracc"
	"github.com/ccthomas/gridiron/pkg/database"
	gridironLogger "github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/rabbitmq"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger := gridironLogger.Get()
	defer logger.Sync()

	logger.Info("Starting Gridiron...")

	logger.Debug("Connect to rabbit mq.")
	amqpConnection := rabbitmq.ConnectRabbitMQ()
	defer amqpConnection.Close()

	ch, err := rabbitmq.OpenChannel(amqpConnection)
	if err != nil {
		logger.Fatal("Failed to open channel.")
	}
	defer ch.Close()

	logger.Debug("Construct rabbit router.")
	rmq := rabbitmq.NewRouter(ch)

	logger.Debug("Connect to database.")

	db := database.ConnectPostgres()
	defer db.Close()

	teamRepo := &team.TeamRepositoryImpl{
		DB: db,
	}

	tenantRepo := &tenant.TenantRepositoryImpl{
		DB: db,
	}

	userRepo := &useracc.UserAccountRepositoryImpl{
		DB: db,
	}

	logger.Debug("Construct handlers.")
	handler := api.NewHandlers(db, rmq, teamRepo, tenantRepo, userRepo)

	logger.Debug("Construct router.")
	r := mux.NewRouter()

	logger.Debug("Route paths to handlers.")
	handler.RouteApis(r, rmq)

	logger.Debug("Handle router with http.")
	http.Handle("/", r)

	serverPort := os.Getenv("SERVER_PORT")

	logger.Info("Server starting...", zap.String("port", serverPort))
	err = http.ListenAndServe(fmt.Sprintf(":%v", serverPort), nil)
	if err != nil {
		logger.Fatal("Failed to listen and serve app.", zap.Error(err))
	}
}
