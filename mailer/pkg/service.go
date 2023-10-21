package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GDGVIT/opengraph-thumbnail-backend/pkg/logger"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	Logger        logger.Logger
	Mailer        *MailInstance
	Worker        *WorkerService
	MessageBroker MessageBroker
}

type WorkerService struct {
	QueueName    string
	NumOfWorkers int
}

type MessageBroker interface {
	GetChannel() *amqp.Channel
}

func NewService(logger logger.Logger, mailer *MailInstance, worker *WorkerService, messageBroker MessageBroker) *Service {
	return &Service{
		Logger:        logger,
		Mailer:        mailer,
		Worker:        worker,
		MessageBroker: messageBroker,
	}
}

func (svc *Service) StartConsumer(ctx context.Context) {
	ch := svc.MessageBroker.GetChannel()
	q, err := ch.QueueDeclare(
		svc.Worker.QueueName, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // nowait
		nil,                  // arguments
	)
	if err != nil {
		svc.Logger.Fatal(errors.Wrap(err, "failed to declare a queue"))
		return
	}

	ch.Qos(
		svc.Worker.NumOfWorkers, // Set the maximum number of unacknowledged messages to the number of workers.
		0,
		false,
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		svc.Logger.Fatal(errors.Wrap(err, "failed to register a consumer"))
		return
	}

	svc.Logger.Info(" [*] Waiting for messages. To exit press CTRL+C")

	for i := 0; i < svc.Worker.NumOfWorkers; i++ {
		go func(workerID int) {
			for delivery := range msgs {
				var msg Message
				if err := json.Unmarshal(delivery.Body, &msg); err != nil {
					svc.Logger.Error(errors.Wrap(err, "failed to unmarshal message"))
					continue
				}
				svc.Logger.Info(fmt.Sprintf("Worker %d received a message: %s", workerID, msg))
				if err := svc.Mailer.SendEmail(msg); err != nil {
					svc.Logger.Error(errors.Wrap(err, "failed to send email"))
					continue
				}
				svc.Logger.Info(fmt.Sprintf("Worker %d finished processing message: %s", workerID, msg))
				delivery.Ack(false) // Acknowledge the message
			}
		}(i)

	}

	<-ctx.Done()
}
