package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(ch *amqp.Channel, exchange string, key string, body RabbitMqBody) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchange, // Exchange
		key,      // Routing key (queue name)
		false,    // Mandatory
		false,    // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
	return err
}
