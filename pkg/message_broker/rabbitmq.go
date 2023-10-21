package messagebroker

import (
	"context"
	"log"

	"github.com/GDGVIT/opengraph-thumbnail-backend/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQHelper implements the MessageBroker interface for RabbitMQ.
type RabbitMQHelper struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	numOfWorkers int
	logger       logger.Logger
}

// NewRabbitMQHelper initializes a new RabbitMQHelper instance.
func NewRabbitMQHelper(rabbitMQURL string, numOfWorkers int, logger logger.Logger) (*RabbitMQHelper, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &RabbitMQHelper{conn: conn, ch: ch, numOfWorkers: numOfWorkers, logger: logger}, nil
}

// Publish publishes a message to the RabbitMQ message broker.
func (rh *RabbitMQHelper) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	err := rh.ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	log.Printf("Published message to exchange '%s', routing key '%s'", exchange, routingKey)
	return nil
}

// Get Channel
func (rh *RabbitMQHelper) GetChannel() *amqp.Channel {
	return rh.ch
}

// Close closes the connection to RabbitMQ.
func (rh *RabbitMQHelper) Close() error {
	if rh.ch != nil {
		_ = rh.ch.Close()
	}

	if rh.conn != nil {
		return rh.conn.Close()
	}

	return nil
}
