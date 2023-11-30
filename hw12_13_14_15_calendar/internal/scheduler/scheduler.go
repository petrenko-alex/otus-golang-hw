package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	period       time.Duration
	queueName    string
	storage      storage.Storage
	queueManager queue.RabbitManager

	ctx    context.Context
	logger logger.Logger
}

func New(
	ctx context.Context,
	logg logger.Logger,
	storage storage.Storage,
	manager queue.RabbitManager,
	period time.Duration,
	queueName string,
) Scheduler {
	return Scheduler{
		period:       period,
		queueName:    queueName,
		queueManager: manager,
		storage:      storage,
		ctx:          ctx,
		logger:       logg,
	}
}

func (s *Scheduler) Run() {
	for {
		s.logger.Info("Looking for events to remind...")

		select {
		case <-time.After(s.period):
			events := s.getEvents()
			if events != nil {
				s.sendEvents(events)
			}
		case <-s.ctx.Done():
			s.logger.Info("Scheduler stopped.")

			return
		}
	}
}

func (s *Scheduler) getEvents() *entity.Events {
	events, getErr := s.storage.GetForRemind()
	if getErr != nil {
		s.logger.Error("Error getting events for reminder: " + getErr.Error())

		return nil
	}

	if len(*events) == 0 {
		s.logger.Info("No events to remind about.")

		return nil
	}

	return events
}

func (s *Scheduler) sendEvents(events *entity.Events) {
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

			continue
		}

		s.logger.Info(fmt.Sprintf("Event \"%s\" sent", event.Title))
	}
}
