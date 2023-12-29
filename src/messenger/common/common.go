package common

import amqp "github.com/rabbitmq/amqp091-go"

const (
	Addr         = "amqp://guest:guest@localhost:5672"
	WorkerQueue  = "worker_queue_task"
	DefaultQueue = "default_queue"
)

type AmqpConfig struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}
