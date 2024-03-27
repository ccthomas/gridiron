package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func OpenChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}
