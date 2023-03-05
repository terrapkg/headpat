package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func conn() {
	println("Connecting to fedora-messaging")
	conn, err := amqp.Dial("amqps://rabbitmq.fedoraproject.org/%2Fpublic_pubsub")
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel();
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	msgs, err := ch.Consume(
		os.Getenv("QUEUE_UUID"),
		"",
		true, // autoAck
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	log.Println("Connected to fedora-messaging")
	<-forever

}
