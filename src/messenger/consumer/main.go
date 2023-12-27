package main

import (
	"common"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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

	msgs, err := config.Channel.Consume(
		config.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	common.FailOnError(err, "failed to register a consumer")

	var forever chan struct{}

	go func() {
		for m := range msgs {
			log.Printf("received a message: %s", m.Body)
		}
	}()

	log.Printf("waiting for messages. to exit pres ctrl+c")
	<-forever
}
