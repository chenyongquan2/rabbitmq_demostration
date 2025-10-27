// 负责连接和通道封装
package amqpclient

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var retryInterval = 5 * time.Second

// Connect creates a connection and channel to RabbitMQ with automatic retry.
func Connect(url string) (*amqp.Connection, *amqp.Channel) {
	for {
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Printf("[WARN] Failed to connect to RabbitMQ: %v. Retrying in %v...", err, retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("[WARN] Failed to open channel: %v. Retrying in %v...", err, retryInterval)
			conn.Close()
			time.Sleep(retryInterval)
			continue
		}

		log.Printf("[INFO] Connected to RabbitMQ: %s", url)
		return conn, ch
	}
}
