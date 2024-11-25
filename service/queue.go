package service

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskQueue interface {
	PushJobOnQueue(payload []byte) error
	PullJobFromQueue() (<-chan amqp.Delivery, error)
}

type RabbitMqQueue struct {
	ch    amqp.Channel
	queue amqp.Queue
}

func NewTaskQueue(dsn, name string) *RabbitMqQueue {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to open a queue: %v", err)
	}

	log.Println("Rabitmq client connected")

	return &RabbitMqQueue{
		queue: q,
		ch:    *ch,
	}
}

func (t *RabbitMqQueue) PushJobOnQueue(payload []byte) error {
	err := t.ch.Publish(
		"",
		t.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         payload,
		},
	)
	if err != nil {
		log.Println("Error while passing message to queue: ", err)
		return err
	}

	return nil
}

func (t *RabbitMqQueue) PullJobFromQueue() (<-chan amqp.Delivery, error) {
	err := t.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	msgs, err := t.ch.Consume(
		t.queue.Name, // queue
		"",           // consumer
		// false,        // auto-ack
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
