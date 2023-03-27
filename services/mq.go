package services

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewMqProvider() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}

	return conn, ch
}
