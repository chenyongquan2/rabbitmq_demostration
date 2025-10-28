// // 用于创建 Exchange / Queue / Binding（幂等声明）
package amqpclient

// import (
// 	"log"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

//cient usage belike:
// // 示例示范了完整拓扑声明的幂等做法，目的是保证：
// // 无论消费者是否上线，Producer 也能保证消息能被正确路由，不丢失。
// // 这在教学、调试、或系统初期部署阶段确实方便。
// // 但实际在生产环境中：
// // ⚠️ 真正的最佳实践是由 消费者 创建 Queue 与 Binding，
// // 而不是 Producer。 这样即便 Queue 名称或数量变化，也不会需要修改 Producer 代码。
// amqpclient.DeclareTopology(ch, exchangeName, exchangeType, queueName, routingKey)

// // DeclareTopology declares exchange, queue and binding (idempotent).
// func DeclareTopology(ch *amqp.Channel, exchange, exchangeType, queue, routingKey string) {
// 	err := ch.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)
// 	if err != nil {
// 		log.Fatalf("[ERROR] Failed to declare exchange: %v", err)
// 	}

// 	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
// 	if err != nil {
// 		log.Fatalf("[ERROR] Failed to declare queue: %v", err)
// 	}

// 	// 在 direct 类型中：
// 	// routing_key 必须与 binding_key 完全相等才匹配。
// 	// 所以，当我们构造最简单的一对一示例时，直接使用同一个字符串 “task.create” 作为两边的匹配键最直观：
// 	// Producer 发送时用 routing_key="task.create"
// 	// Consumer 绑定时用 binding_key="task.create"
// 	// 因此，为了避免声明两个变量（值相同），我们直接用同一个变量名。

// 	err = ch.QueueBind(
// 		q.Name,     // Queue name
// 		routingKey, // binding_key
// 		exchange,   // Exchange name
// 		false,      // no-wait
// 		nil,        // args
// 	)
// 	if err != nil {
// 		log.Fatalf("[ERROR] Failed to bind queue: %v", err)
// 	}

// 	log.Printf("[INFO] Exchange(%s) ↔ Queue(%s) Binding(%s) initialized", exchange, q.Name, routingKey)
// }
