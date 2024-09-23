package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

// NewRabbitMQService creates a new instance of RabbitMQService and declares a queue
func NewRabbitMQService(rabbitmqURL, queueName string) (*RabbitMQService, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	// Create a new channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable (will survive broker restart)
		false,     // Auto delete when unused
		false,     // Exclusive (only accessible by this connection)
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{
		Connection: conn,
		Channel:    ch,
		Queue:      q,
	}, nil
}

// Publish sends a message to the RabbitMQ queue with a routing key
func (r *RabbitMQService) Publish(message string) error {
	err := r.Channel.Publish(
		"",           // exchange
		r.Queue.Name, // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", message)
	return nil
}

// Consume listens for messages from the RabbitMQ queue
func (r *RabbitMQService) Consume(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQService) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ channel: %v", err)
	}
	if err := r.Connection.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ connection: %v", err)
	}
}
