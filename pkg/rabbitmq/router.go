package rabbitmq

import (
	"encoding/json"

	"github.com/ccthomas/gridiron/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func NewRouter(ch *amqp.Channel) *RabbitMqRouter {
	logger.Get().Debug("New RabbitMQ Router.")
	return &RabbitMqRouter{
		Channel:   ch,
		Exchanges: map[string]bool{},
	}
}

func (r *RabbitMqRouter) HandleFunc(exchange string, key string, f func(body RabbitMqBody)) {
	go r.handleFuncAsync(exchange, key, f)
}

func (r *RabbitMqRouter) handleFuncAsync(exchange string, key string, f func(body RabbitMqBody)) {
	logger.Get().Debug("Handle Func Async.")
	if _, ok := r.Exchanges[exchange]; !ok {
		logger.Get().Debug("Declaring exchange.")
		err := DeclareExchange(r.Channel, exchange, "fanout")
		if err != nil {
			logger.Get().Fatal("Failed to declare exchange.")
		}

		r.Exchanges[exchange] = true
	}

	logger.Get().Debug("Declaring queue.")
	q, err := DeclareQueue(r.Channel, "")
	if err != nil {
		logger.Get().Fatal("Failed to declare queue.")
	}

	logger.Get().Debug("Binding queue.")
	err = QueueBind(r.Channel, q, key, exchange)
	if err != nil {
		logger.Get().Fatal("Failed to bind queue.")
	}

	logger.Get().Debug("Consume message.")
	msgs, err := Consume(r.Channel, q)
	if err != nil {
		logger.Get().Fatal("Failed to consume from queue.")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logger.Get().Debug("Message received from queue.", zap.String("Exchange", exchange), zap.String("Key", key))

			logger.Get().Debug("Unmarshal rabbit mq body")
			var body RabbitMqBody
			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				logger.Get().Error("Failed to parse message", zap.Error(err))
				return
			}

			logger.Get().Debug("Invoking handler")
			f(body)
		}
	}()

	logger.Get().Debug("Waiting to process message.", zap.String("Exchange", exchange), zap.String("Key", key))
	<-forever
}

func (r *RabbitMqRouter) PublishMessage(exchange string, key string, bodies []RabbitMqBody) error {
	logger.Get().Debug("Publish Message.")
	if _, ok := r.Exchanges[exchange]; !ok {
		logger.Get().Debug("Declaring exchange.")
		err := DeclareExchange(r.Channel, exchange, "fanout")
		if err != nil {
			logger.Get().Fatal("Failed to declare exchange.")
		}

		r.Exchanges[exchange] = true
	}

	for _, body := range bodies {
		err := Publish(r.Channel, exchange, key, body)
		if err != nil {
			return err
		}
	}

	return nil
}
