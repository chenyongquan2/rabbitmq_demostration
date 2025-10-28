package amqpclient

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeclareQueueAndBind: 声明 Exchange、Queue，并建立绑定关系
// 注意：此函数一般由 Consumer 调用
func DeclareQueueAndBind(ch *amqp.Channel, exchange, exchangeType, queue, bindingKey string) {
	// 声明 Exchange（幂等操作：重复声明不会报错）
	if err := ch.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil); err != nil {
		log.Fatalf("[ERROR] 声明 Exchange 失败: %v", err)
	}

	// 声明 Queue（幂等）
	q, err := ch.QueueDeclare(
		queue,
		true,  // durable：持久化队列
		false, // auto-delete：不自动删除
		false, // exclusive：非独占（多个消费者共享）
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("[ERROR] 声明 Queue 失败: %v", err)
	}

	// 建立 Binding：Exchange 与 Queue 之间的路由规则
	// bindingKey 表示该 Queue 希望接收哪些 routing_key 的消息
	if err := ch.QueueBind(q.Name, bindingKey, exchange, false, nil); err != nil {
		log.Fatalf("[ERROR] 建立绑定关系失败: %v", err)
	}

	log.Printf("[INFO] Exchange(%s) ↔ Queue(%s) 绑定成功，binding_key=%s", exchange, q.Name, bindingKey)
}
