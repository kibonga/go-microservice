package common

import amqp "github.com/rabbitmq/amqp091-go"

const (
	Addr = "amqp://guest:guest@localhost:5672"
)

type Config struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}
