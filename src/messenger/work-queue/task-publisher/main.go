package main

import (
	"common"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	config := common.AmqpConfig{}

	// Create connection
	connectToRabbitMq(&config)
	defer config.Conn.Close()

	// Open Channel
	openChannel(&config)

	// Create Queue
	createQueue(&config, common.WorkerQueue)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stopCh := make(chan bool)
	msg := "hello from TASK-publisher!"
	go publishMsg(&config, ctx, stopCh, msg)

	time.Sleep(20 * time.Second)
	stopCh <- true
}

func connectToRabbitMq(config *common.AmqpConfig) {
	conn, err := amqp.Dial(common.Addr)
	common.FailOnError(err, "failed to connect to rabbitmq")
	config.Conn = conn
}

func openChannel(config *common.AmqpConfig) {
	ch, err := config.Conn.Channel()
	common.FailOnError(err, "failed to open a channel")
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

func publishMsg(config *common.AmqpConfig, ctx context.Context, stopCh chan bool, msg string) {
	var count int64 = 1
	for {
		select {
		case <-stopCh:
			fmt.Println("stopping task-publisher...")
			break
		default:
			config.Channel.PublishWithContext(
				ctx,
				"",
				config.Queue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("%s - %d", msg, count)),
				},
			)
			log.Println("message sent")
			count++
			time.Sleep(time.Millisecond * 700)
		}
	}
}
