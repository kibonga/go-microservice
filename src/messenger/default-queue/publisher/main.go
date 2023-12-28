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
	var config common.AmqpConfig

	conn, err := amqp.Dial(common.Addr)
	common.FailOnError(err, "failed to connect to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	common.FailOnError(err, "failed to open a channel")
	defer ch.Close()
	config.Channel = ch

	q, err := config.DeclareQueue(common.DefaultQueue)
	common.FailOnError(err, "failed to declare queue")
	config.Queue = q

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count := 1

	stopCh := make(chan struct{})
	func() {
		for {
			select {
			case <-stopCh:
				fmt.Println("stopping publishing messages")
				break
			default:
				body := fmt.Sprintf("This is message = %d", count)
				err = config.PlainPublish(ctx, body)
				common.FailOnError(err, "failed to publish a message")
				count++
				log.Printf("message sent: %s", body)
				time.Sleep(time.Second)
			}
		}
	}()

	//scanner := bufio.NewScanner(os.Stdin)
	//
	//fmt.Println("Enter 'stop' to close connection...")
	//for scanner.Scan() {
	//	text := scanner.Text()
	//	if text == "stop" {
	//		close(stopCh)
	//		break
	//	}
	//}
}
