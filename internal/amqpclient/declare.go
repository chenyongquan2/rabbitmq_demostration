// 用于创建 Exchange / Queue / Binding（幂等声明）
package amqpclient

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeclareTopology declares exchange, queue and binding (idempotent).
func DeclareTopology(ch *amqp.Channel, exchange, exchangeType, queue, routingKey string) {
	err := ch.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[ERROR] Failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[ERROR] Failed to declare queue: %v", err)
	}

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		log.Fatalf("[ERROR] Failed to bind queue: %v", err)
	}

	log.Printf("[INFO] Exchange(%s) ↔ Queue(%s) Binding(%s) initialized", exchange, q.Name, routingKey)
}
