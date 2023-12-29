package main

import (
	"common"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	config := common.AmqpConfig{}

	// Connect to RabbitMq
	connectToRabbitMq(&config)
	defer config.Conn.Close()

	// Open up a Channel
	openChannel(&config)
	defer config.Channel.Close()

	// Set fair dispatch for channel
	setFairChannelDispatch(&config)

	// Create Queue
	createQueue(&config, common.WorkerQueue)

	// Consume messages
	messages := bufferMessages(&config)

	stopCh := make(chan bool)
	// Consume messages
	go consumeMessages(messages, stopCh)

	time.Sleep(time.Second * 15)
	stopCh <- true

}

func connectToRabbitMq(config *common.AmqpConfig) {
	conn, err := amqp.Dial(common.Addr)
	common.FailOnError(err, "failed to connect to rabbitmq")
	config.Conn = conn
	log.Println("connected to rabbitmq")
}

func openChannel(config *common.AmqpConfig) {
	ch, err := config.Conn.Channel()
	common.FailOnError(err, "failed to open channel")
	config.Channel = ch
}

func createQueue(config *common.AmqpConfig, name string) {
	q, err := config.Channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil)
	common.FailOnError(err, "failed to create queue")
	config.Queue = q
}

func bufferMessages(config *common.AmqpConfig) <-chan amqp.Delivery {
	messages, err := config.Channel.Consume(
		config.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	common.FailOnError(err, "failed to consume messages")
	return messages
}

func consumeMessages(messages <-chan amqp.Delivery, stopCh chan bool) {
	for m := range messages {
		select {
		case <-stopCh:
			fmt.Println("stopping consumer...")
			break
		default:
			log.Printf("received message: %s", m.Body)
			time.Sleep(time.Millisecond * 400)
			log.Println("message consumed")
			m.Ack(false)
		}
	}
}

func setFairChannelDispatch(config *common.AmqpConfig) {
	err := config.Channel.Qos(
		1,
		0,
		false)
	common.FailOnError(err, "failed to set QoS")
}
