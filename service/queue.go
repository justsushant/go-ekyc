package service

import (
	"encoding/json"
	"log"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskQueue interface {
	PushJobOnQueue(payload types.QueuePayload) error
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

	return &RabbitMqQueue{
		queue: q,
		ch:    *ch,
	}
}

func (t *RabbitMqQueue) PushJobOnQueue(payload types.QueuePayload) error {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error while marshalling JSON: ", err)
	}

	err = t.ch.Publish(
		"",
		t.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonBytes,
		},
	)
	if err != nil {
		log.Println("Error while passing message to queue: ", err)
		return err
	}

	return nil
}
