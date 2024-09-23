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

func (r *RabbitMQService) Consume(queueName string) (<-chan amqp.Delivery, error) {
	channel, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		queueName, // kuyruk adı
		"",        // consumer tag (boş bırakabilirsiniz)
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
	r.Channel.Close()
	r.Connection.Close()
}
