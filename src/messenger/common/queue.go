package common

import amqp "github.com/rabbitmq/amqp091-go"

func (config *AmqpConfig) DeclareQueue(name string) (amqp.Queue, error) {
	return config.Channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil)
}
