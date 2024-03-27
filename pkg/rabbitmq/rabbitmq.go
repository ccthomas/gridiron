package rabbitmq

import (
	"fmt"
	"os"

	"github.com/ccthomas/gridiron/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMqBody struct {
	DataVersion string `json:"data_version"`
	Data        any    `json:"data"`
}

type RabbitMqRouter struct {
	Channel   *amqp.Channel
	Exchanges map[string]bool
}

func ConnectRabbitMQ() *amqp.Connection {
	logger.Get().Debug("Connecting to rabbit mq")

	mqHost := os.Getenv("RABBITMQ_HOST")
	mqUser := os.Getenv("RABBITMQ_USER")
	mqPassword := os.Getenv("RABBITMQ_PASSWORD")

	// Construct the connection string
	connStr := fmt.Sprintf("amqp://%s:%s@%s:5672/",
		mqUser, mqPassword, mqHost)

	logger.Get().Debug("Open a connection to rabbit mq.")
	conn, err := amqp.Dial(connStr)
	if err != nil {
		logger.Get().Fatal("Error opening rabbit mq connection:", zap.Error(err))
	}

	logger.Get().Debug("Connection to postgres database was successful.")
	return conn
}
