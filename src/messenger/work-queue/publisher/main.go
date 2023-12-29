package main

import (
	"common"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

const (
	postErrMsg = " from default-worker publisher"
)

func main() {
	var config common.AmqpConfig

	// Connect to rabbitmq
	connectToRabbitMq(&config)
	defer config.Conn.Close()

	// Create channel
	createChannel(&config)
	defer config.Channel.Close()

	// Create queue
	createQueue(&config, common.WorkerQueue)

	// Create context
	var ctx context.Context
	ctx = createContext()

	stopCh := make(chan struct{})
	msg := "hello world! from publisher"
	// Init publish goroutine
	go publishMsg(&config, ctx, stopCh, msg)

	time.Sleep(10 * time.Second)
	stopCh <- struct{}{}
}

func connectToRabbitMq(config *common.AmqpConfig) {
	var connCount int64
	backoff := 1 * time.Second

	for {
		conn, err := amqp.Dial(common.Addr)
		if err != nil {
			fmt.Println("rabbitmq not ready yet...")
			connCount++
		} else {
			fmt.Println("connected to rabbitmq")
			config.Conn = conn
			break
		}

		if connCount > 5 {
			common.FailOnError(err, "conn attempts exceeded limit")
		}

		backoff = time.Duration(math.Pow(float64(connCount), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
	}
}

func createChannel(config *common.AmqpConfig) {
	ch, err := config.Conn.Channel()
	common.FailOnError(err, "failed to create channel")
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

func createContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return ctx
}

func publishMsg(config *common.AmqpConfig, ctx context.Context, stopCh chan struct{}, msg string) {
	var count = 1
	for {
		select {
		case <-stopCh:
			fmt.Println("stopping channel...")
			break
		default:
			err := config.PlainPublish(ctx, fmt.Sprintf("%s - %d", msg, count))
			common.FailOnError(err, "failed to publish message"+postErrMsg)
			log.Printf("message sent")
			count++
			time.Sleep(time.Second)
		}
	}
}
