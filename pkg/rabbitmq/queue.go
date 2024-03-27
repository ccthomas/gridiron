package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func DeclareQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return q, err
}

func QueueBind(ch *amqp.Channel, q amqp.Queue, key string, exchange string) error {
	err := ch.QueueBind(
		q.Name,   // queue name
		key,      // routing key
		exchange, // exchange
		false,
		nil,
	)

	return err
}
