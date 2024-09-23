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

// NewRabbitMQService creates a new instance of RabbitMQService
func NewRabbitMQService(rabbitmqURL, queueName string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
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

// Publish sends a message to the RabbitMQ queue
func (r *RabbitMQService) Publish(routeKey string, message string) error {
	err := r.Channel.Publish(
		"",       // exchange
		routeKey, // routing key (queue name)
		false,    // mandatory
		false,    // immediate
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

// Consume listens for messages from the RabbitMQ queue with a specific routing key
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
	r.Channel.Close()
	r.Connection.Close()
}
