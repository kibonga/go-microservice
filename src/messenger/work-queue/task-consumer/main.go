package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

import (
	"common"
)

func main() {
	config := common.AmqpConfig{}

	// Create Connection
	connectToRabbitMq(&config)
	defer config.Conn.Close()

	// Open Channel
	openChannel(&config)
	defer config.Channel.Close()

	// Create Queue
	createQueue(&config, common.WorkerQueue)

	stopCh := make(chan bool)
	// Receive messages
	messages := receiveMessages(&config)
	go consumeMessages(messages, stopCh)

	time.Sleep(time.Second * 15)
	stopCh <- true
}

func connectToRabbitMq(config *common.AmqpConfig) {
	conn, err := amqp.Dial(common.Addr)
	common.FailOnError(err, "failed to connect to rabbitmq")
	config.Conn = conn
}

func openChannel(config *common.AmqpConfig) {
	ch, err := config.Conn.Channel()
	common.FailOnError(err, "failed to open channel")
	config.Channel = ch
}

func createQueue(config *common.AmqpConfig, name string) {
	q, err := config.Channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil)
	common.FailOnError(err, "failed to create queue")
	config.Queue = q
}

func receiveMessages(config *common.AmqpConfig) <-chan amqp.Delivery {
	messages, err := config.Channel.Consume(
		config.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	common.FailOnError(err, "failed to consume message")
	return messages
}

func consumeMessages(messages <-chan amqp.Delivery, stopCh chan bool) {
	for m := range messages {
		select {
		case <-stopCh:
			fmt.Println("stopping TASK-CONSUMER...")
			break
		default:
			log.Printf("received message: %s", m.Body)
			time.Sleep(time.Millisecond * 600)
			log.Println("message consumed")
		}
	}
}
