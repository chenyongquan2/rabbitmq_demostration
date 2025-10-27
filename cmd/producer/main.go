package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"rabbitmq-demo/internal/amqpclient"

	amqp "github.com/rabbitmq/amqp091-go"
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

	if err := ch.Confirm(false); err != nil {
		log.Fatalf("Failed to enable confirm mode: %v", err)
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= 5; i++ {
		body := fmt.Sprintf("Task #%d", i)
		err := ch.PublishWithContext(ctx, exchangeName, routingKey, true, false, amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         []byte(body),
			Timestamp:    time.Now(),
		})
		if err != nil {
			log.Printf("[WARN] Failed to publish message: %v", err)
			continue
		}

		confirmed := <-confirms
		if confirmed.Ack {
			log.Printf("[OK] Message confirmed: %s", body)
		} else {
			log.Printf("[WARN] Message not confirmed: %s", body)
		}
	}

	log.Println("[DONE] All messages published.")
}
