package main

import (
	"log"

	"rabbitmq-demo/internal/amqpclient"
)

const (
	amqpURL      = "amqp://guest:guest@localhost:5672/"
	exchangeName = "task_exchange"
	exchangeType = "direct"
	queueName    = "task_queue"
	routingKey   = "task.create"
)

func main() {
	conn, ch := amqpclient.Connect(amqpURL)
	defer conn.Close()
	defer ch.Close()

	amqpclient.DeclareTopology(ch, exchangeName, exchangeType, queueName, routingKey)

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("[INFO] Consumer started. Waiting for messages...")

	for d := range msgs {
		log.Printf("[RECV] Message: %s", string(d.Body))
		// Simulate business processing...
		d.Ack(false)
	}
}
