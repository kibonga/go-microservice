package main

import (
	"bytes"
	"common"
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

	q, err := config.DeclareQueue(common.WorkerQueue)
	common.FailOnError(err, "failed to declare queue")
	config.Queue = q

	msgs, err := config.RecvMsgs()
	common.FailOnError(err, "failed to register a consumer")

	var forever chan struct{}

	go func() {
		for m := range msgs {
			log.Printf("received a message: %s", m.Body)
			dotCount := bytes.Count(m.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("message is done")
		}
	}()

	log.Printf("waiting for messages. to exit pres ctrl+c")
	<-forever
}
