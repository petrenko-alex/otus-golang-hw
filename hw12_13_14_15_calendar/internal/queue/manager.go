package queue

import (
	"context"
	"errors"
	"fmt"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrNotConnected = errors.New("need to connect first")

type RabbitManager struct {
	connection *amqp.Connection
	channel    *amqp.Channel

	cfg    config.Config
	logger logger.Logger
	ctx    context.Context
}

func NewRabbitManager(ctx context.Context, cfg config.Config, logger logger.Logger) RabbitManager {
	return RabbitManager{
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
	}
}

func (q *RabbitManager) Connect() error {
	var dialErr, chanErr error

	q.connection, dialErr = amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		q.cfg.RabbitMQServer.Login,
		q.cfg.RabbitMQServer.Password,
		q.cfg.RabbitMQServer.Host,
		q.cfg.RabbitMQServer.Port,
	))
	if dialErr != nil {
		return dialErr
	}

	q.channel, chanErr = q.connection.Channel()
	if chanErr != nil {
		return chanErr
	}

	return nil
}

func (q *RabbitManager) Close() error {
	var err error

	if q.connection != nil {
		err = q.connection.Close()
	}

	if q.channel != nil {
		err = q.channel.Close()
	}

	return err
}

func (q *RabbitManager) CreateQueue(name string) (amqp.Queue, error) {
	if q.channel == nil {
		return amqp.Queue{}, ErrNotConnected
	}

	queue, queueErr := q.channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if queueErr != nil {
		return amqp.Queue{}, queueErr
	}

	return queue, nil
}

func (q *RabbitManager) Produce(queueName string, jsonMsg []byte) error {
	if q.channel == nil {
		return ErrNotConnected
	}

	return q.channel.PublishWithContext(
		q.ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMsg,
		},
	)
}

func (q *RabbitManager) Consume(queueName string) (<-chan amqp.Delivery, error) {
	if q.channel == nil {
		return nil, ErrNotConnected
	}

	messages, consumeErr := q.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if consumeErr != nil {
		return nil, consumeErr
	}

	return messages, nil
}
