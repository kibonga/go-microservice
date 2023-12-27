package main

import (
	"common"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	var config common.Config

	conn, err := amqp.Dial(common.Addr)
	common.FailOnError(err, "failed to connect to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	common.FailOnError(err, "failed to open a channel")
	defer ch.Close()
	config.Channel = ch

	q, err := config.DeclareQueue("hello")
	common.FailOnError(err, "failed to declare queue")
	config.Queue = q

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "hello from main messenger"

	err = config.PlainPublish(ctx, body)
	common.FailOnError(err, "failed to publish a message")
	log.Printf("message sent: %s", body)
}
