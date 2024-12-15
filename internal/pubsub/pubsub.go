package pubsub

import (
	"context"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Used in DeclareAndBind
const (
	DurableQueue int = iota
	TransientQueue
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// Marshal the value to JSON bytes
	jsonVal, err := json.Marshal(val)
	if err != nil {
		returnError := "An error occured while JSON - marshalling the provided value: " + err.Error()
		return errors.New(returnError)
	}

	// ch.PublishWithContext config
	ctx := context.Background()
	mandatory := false
	immediate := false
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonVal,
	}

	// PublishWithContext
	err = ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
	if err != nil {
		returnError := "An error occured when publishing: " + err.Error()
		return errors.New(returnError)
	}

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	// Create a new channel
	rbtChannel, err := conn.Channel()
	if err != nil {
		returnError := "An error occured while creating a new RabbitMQ Channel: " + err.Error()
		return nil, amqp.Queue{}, errors.New(returnError)
	}

	// Declare a new Queue
	durable := false
	autoDelete := false
	exclusive := false

	// DurableQueue (0), TransientQueue (1)
	switch simpleQueueType {
	case DurableQueue:
		durable = true
		autoDelete = false
		exclusive = false
	case TransientQueue:
		durable = false
		autoDelete = true
		exclusive = true
	}

	queue, err := rbtChannel.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		returnError := "An error occured while declaring a Queue: " + err.Error()
		return nil, amqp.Queue{}, errors.New(returnError)
	}

	// Bind the Queue to the Exchange
	err = rbtChannel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		returnError := "An error occured while binding a Queue: " + err.Error()
		return nil, amqp.Queue{}, errors.New(returnError)
	}

	return rbtChannel, queue, nil
}
