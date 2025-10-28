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
	//exchangeName = "task_exchange_test_ae"
	//exchangeName = "task_exchange-not-exist-before"
	exchangeType = "direct"      // direct 交换器：根据 routing_key 精确匹配路由
	routingKey   = "task.create" // 消息的“标签”或“类别”
)

func main() {
	conn, ch := amqpclient.Connect(amqpURL)
	defer conn.Close()
	defer ch.Close()

	aeEnable := false
	var args amqp.Table

	aeName := "unrouted_exchange" // 定义备用交换器名称

	if aeEnable {
		args = amqp.Table{
			"alternate-exchange": aeName, // 👈 指定 AE
		}
	}

	// 声明主 Exchange。当 aeEnable=true 时会带 AE 参数
	// 声明一个 Exchange（交换器），这是消息路由的“中转站”
	// 第二个参数是类型：direct、fanout、topic、headers 等
	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable：持久化，Broker 重启后仍存在
		false, // auto-delete：不自动删除
		false, // internal：false 表示应用可直接使用
		false, // no-wait：阻塞等待服务器确认
		//nil,   // arguments：额外参数（如 Alternate Exchange）
		args,
	)
	if err != nil {
		log.Fatalf("[ERROR] 声明 Exchange 失败: %v", err)
	}

	// 如果启用了 AE，则声明 AE 自身及其绑定队列
	if aeEnable {
		err = ch.ExchangeDeclare(
			aeName,
			"fanout", // 一般用 fanout，广播所有未路由消息
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("[ERROR] 声明 AE Exchange 失败: %v", err)
		}
		_, err = ch.QueueDeclare(
			"unrouted_queue",
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,
		)
		if err != nil {
			log.Fatalf("[ERROR] 声明 AE Queue 失败: %v", err)
		}
		err = ch.QueueBind(
			"unrouted_queue",
			"", // fanout 模式忽略 routing_key
			aeName,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("[ERROR] 绑定 AE Queue 失败: %v", err)
		}
		log.Printf("[INFO] AE 已启用: %s → unrouted_queue", aeName)
	}

	// 启用发布确认模式：Producer 可以确认消息是否被 Broker 接收
	// （类似 MySQL 的事务提交确认）
	if err := ch.Confirm(false); err != nil {
		log.Fatalf("Failed to enable confirm mode: %v", err)
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// 新增：注册 mandatory 返回通道，用于 mandatory 返回通知
	returns := ch.NotifyReturn(make(chan amqp.Return, 1))

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= 5; i++ {
		body := fmt.Sprintf("Task #%d", i)

		// 发布消息：routing_key 决定了消息“分类”
		// mandatory=true：若无法投递到任何队列，将返回给生产者（防止静默丢失）
		err = ch.PublishWithContext(
			ctx,
			exchangeName, // 要发送的 Exchange
			//"fake.routingkey", //Todo:模拟 没有任何 Queue 使用 binding_key=fake.routingkey 绑定在此 Exchange
			routingKey, // routing_key：消息的“路由标签”
			true,       // mandatory，消息投递必须成功，否则退回给发送者”
			false,      // immediate（已废弃）
			amqp.Publishing{
				ContentType:  "text/plain",    // 内容类型
				DeliveryMode: amqp.Persistent, // 2=持久化消息，Broker 会写入磁盘
				Body:         []byte(body),    // 消息体
				Timestamp:    time.Now(),
			},
		)
		if err != nil {
			log.Printf("[WARN] Failed to publish message: %v", err)
			continue
		}

		// >>> 新增：检查是否有 mandatory 返回（路由失败的消息）
		select {
		case ret := <-returns:
			log.Printf("[WARN] Returned message: code=%d, reason=%s, exchange=%s, key=%s, body=%s",
				ret.ReplyCode, ret.ReplyText, ret.Exchange, ret.RoutingKey, string(ret.Body))
		default:
			// no return event
		}

		// 从确认通道读确认结果，确认 Broker 已经收到消息
		if confirmed := <-confirms; confirmed.Ack {
			log.Printf("[OK] Message confirmed: %s", body)
		} else {
			log.Printf("[WARN] Message not confirmed: %s", body)
		}
	}

	log.Println("[DONE] All messages published.")
}
