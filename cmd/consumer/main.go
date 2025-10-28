package main

import (
	"log"
	"rabbitmq-demo/internal/amqpclient"
)

const (
	amqpURL      = "amqp://guest:guest@localhost:5672/"
	exchangeName = "task_exchange"
	exchangeType = "direct"
	queueName    = "task_queue"  // 队列名：消息实体缓存处
	bindingKey   = "task.create" // 绑定键：与 routing_key 匹配
)

func main() {
	conn, ch := amqpclient.Connect(amqpURL)
	defer conn.Close()
	defer ch.Close()

	// Consumer 负责声明 Queue + Binding
	amqpclient.DeclareQueueAndBind(ch, exchangeName, exchangeType, queueName, bindingKey)

	// 注册消费者，autoAck=false 表示需要手动 ACK 确认
	msgs, err := ch.Consume(
		queueName,
		"",
		false, // autoAck
		false, // exclusive：多个消费者可共享队列
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("[INFO] Consumer started. Waiting for messages...")

	for d := range msgs {
		log.Printf("[RECV] Message: %s", string(d.Body))

		// 模拟业务逻辑处理...
		// 手动回复确认：告诉 Broker 消息已处理完，可从队列删除
		d.Ack(false)
	}
}
