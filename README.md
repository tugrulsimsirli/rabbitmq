
# RabbitMQ Service for Go

This package provides an easy-to-use implementation for integrating RabbitMQ into your Go projects. It simplifies the process of setting up channels, publishing messages, and consuming messages, allowing developers to focus on business logic rather than boilerplate RabbitMQ configuration.

## Features

- **Channel Creation**: Easily create and manage RabbitMQ channels for your services.
- **Service Management**: Start new RabbitMQ services with minimal configuration.
- **Message Publishing**: Simple methods for publishing messages to RabbitMQ queues.
- **Message Consuming**: Easily consume messages from RabbitMQ queues with configurable handlers.

## Installation

To install this package, simply run:

```bash
go get github.com/tugrulsimsirli/rabbitmq
```

## Usage

Here is an example of how to use the RabbitMQ service in your project:

```go
package main

import (
    "github.com/tugrulsimsirli/rabbitmq"
    "log"
)

func main() {
    // Start a new RabbitMQ service
    service, err := rabbitmq.NewRabbitMQService("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to create RabbitMQ service: %v", err)
    }
    
    // Create a new channel
    channel, err := service.CreateChannel()
    if err != nil {
        log.Fatalf("Failed to create RabbitMQ channel: %v", err)
    }

    // Publish a message
    err = service.Publish("Hello, RabbitMQ!")
    if err != nil {
        log.Fatalf("Failed to publish message: %v", err)
    }

    // Consume messages
    err = service.Consume("queue_name", func(message string) {
        log.Printf("Received message: %s", message)
    })
    if err != nil {
        log.Fatalf("Failed to consume messages: %v", err)
    }

    // Close the channel
    err := service.CloseChannel()
    if err != nil {
        log.Fatalf("Failed to create RabbitMQ channel: %v", err)
    }

    // Close the RabbitMQ service
    service.Close()
}
```

## Configuration

- **RabbitMQ URL**: You can provide the RabbitMQ URL when starting the service. For example: `amqp://guest:guest@localhost:5672/`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Feel free to submit issues, fork the repository, and send pull requests. Contributions are welcome!
