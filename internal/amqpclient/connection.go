// 负责连接和通道封装
package amqpclient

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var retryInterval = 5 * time.Second

// Connect：负责连接 RabbitMQ 服务器并返回频道（Channel）
// 包含简单的自动重试逻辑，保证服务启动时 RabbitMQ 尚未就绪也能自动等待。
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
