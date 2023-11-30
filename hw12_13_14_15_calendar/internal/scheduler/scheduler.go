package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	queueName    string
	storage      storage.Storage
	queueManager queue.RabbitManager

	ctx    context.Context
	logger logger.Logger
}

func New(
	queueName string,
	ctx context.Context,
	logg logger.Logger,
	storage storage.Storage,
	manager queue.RabbitManager,
) Scheduler {
	return Scheduler{
		queueName:    queueName,
		queueManager: manager,
		storage:      storage,
		ctx:          ctx,
		logger:       logg,
	}
}

func (s *Scheduler) SendEvents() {
	events, getErr := s.storage.GetForRemind()
	if getErr != nil {
		s.logger.Error("Error getting events for reminder: " + getErr.Error())
	}

	if len(*events) == 0 {
		s.logger.Info("No events to remind about.")
	}

	for _, event := range *events {
		eventMsg := event.ToMsg()
		jsonMsg, marshalErr := json.Marshal(eventMsg)
		if marshalErr != nil {
			s.logger.Error("Error sending msg to RabbitMQ: " + marshalErr.Error())

			continue
		}

		produceErr := s.queueManager.Produce(s.queueName, jsonMsg)
		if produceErr != nil {
			s.logger.Error("Error sending msg to RabbitMQ: " + produceErr.Error())

			//return 1
		}

		s.logger.Info(fmt.Sprintf("Event \"%s\" sent", event.Title))
	}
}
