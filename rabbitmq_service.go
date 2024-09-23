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
func (r *RabbitMQService) Publish(routingKey string, message string) error {
	// Ensure the direct exchange is declared
	err := r.Channel.ExchangeDeclare(
		"direct_exchange", // Exchange name
		"direct",          // Exchange type (direct)
		true,              // Durable (survives broker restart)
		false,             // Auto-deleted
		false,             // Internal
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		return err
	}

	// Publish the message to the exchange with the routing key
	err = r.Channel.Publish(
		"direct_exchange", // Exchange
		routingKey,        // Routing key
		false,             // Mandatory
		false,             // Immediate
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
	// Start consuming messages from the queue
	msgs, err := r.Channel.Consume(
		queueName, // Queue
		"",        // Consumer name
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// ConsumeWithRoute listens for messages from the RabbitMQ queue with a specific routing key
func (r *RabbitMQService) ConsumeWithRoute(queueName string, routingKey string) (<-chan amqp.Delivery, error) {
	// Declare a direct exchange if not already declared
	err := r.Channel.ExchangeDeclare(
		"direct_exchange", // Exchange name
		"direct",          // Exchange type (direct)
		true,              // Durable
		false,             // Auto-deleted
		false,             // Internal
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		return nil, err
	}

	// Bind the queue to the exchange with the specific routing key
	err = r.Channel.QueueBind(
		queueName,         // Queue name
		routingKey,        // Routing key
		"direct_exchange", // Exchange name
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		return nil, err
	}

	// Start consuming messages from the queue
	msgs, err := r.Channel.Consume(
		queueName, // Queue name
		"",        // Consumer name
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
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
