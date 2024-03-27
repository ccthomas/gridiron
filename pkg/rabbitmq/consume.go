package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func Consume(ch *amqp.Channel, q amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	return msgs, err
}
