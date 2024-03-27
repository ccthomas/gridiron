package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func DeclareExchange(ch *amqp.Channel, name string, kind string) error {
	err := ch.ExchangeDeclare(
		name,  // name
		kind,  // kind
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)

	return err
}
