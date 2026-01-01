package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageBroker defines the wrapper for RabbitMQ operations
type MessageBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewMessageBroker creates a new RabbitMQ message broker
func NewMessageBroker(url string) (*MessageBroker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &MessageBroker{
		conn:    conn,
		channel: ch,
	}, nil
}

// Publish publishes a message to an exchange with a routing key
func (b *MessageBroker) Publish(ctx context.Context, exchange, routingKey string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return b.channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
}

// Subscribe registers a consumer for a specific queue
func (b *MessageBroker) Subscribe(queueName string, handler func([]byte) error) error {
	msgs, err := b.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				// In a real app, we'd use a logger here
				fmt.Printf("Error handling message from %s: %v\n", queueName, err)
			}
		}
	}()

	return nil
}

// DeclareQueue ensures a queue exists
func (b *MessageBroker) DeclareQueue(name string) (amqp.Queue, error) {
	return b.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareExchange ensures an exchange exists
func (b *MessageBroker) DeclareExchange(name, kind string) error {
	return b.channel.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// BindQueue binds a queue to an exchange
func (b *MessageBroker) BindQueue(queueName, routingKey, exchangeName string) error {
	return b.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
}

// Close gracefully shuts down the broker connection
func (b *MessageBroker) Close() {
	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
}
