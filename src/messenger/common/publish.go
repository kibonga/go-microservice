package common

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (config *AmqpConfig) PlainPublish(ctx context.Context, body string) error {
	err := config.Channel.PublishWithContext(
		ctx,
		"",
		config.Queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	FailOnError(err, "failed to publish with context")
	return err
}
