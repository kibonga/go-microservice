package common

import amqp "github.com/rabbitmq/amqp091-go"

func (config *AmqpConfig) RecvMsgs() (<-chan amqp.Delivery, error) {
	msgs, err := config.Channel.Consume(
		config.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "failed to deliver messages...")
	return msgs, nil
}
