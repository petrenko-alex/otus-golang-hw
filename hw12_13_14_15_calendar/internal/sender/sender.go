package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	ctx    context.Context
	logger logger.Logger
}

func New(ctx context.Context, logger logger.Logger) Sender {
	return Sender{
		ctx:    ctx,
		logger: logger,
	}
}

func (s *Sender) Run(channel <-chan amqp.Delivery) {
	for {
		s.logger.Info("Waiting for events...")

		select {
		case <-s.ctx.Done():
			s.logger.Info("Sender stopped.")

			return
		case msg := <-channel:
			eventMsg := entity.EventMsg{}

			unmarshalErr := json.Unmarshal(msg.Body, &eventMsg)
			if unmarshalErr != nil {
				s.logger.Error("Error reading msg from channel: " + unmarshalErr.Error())

				continue
			}

			s.logger.Info(fmt.Sprintf(
				"Sending reminder about \"%s\" event to #%d user. Event time: %s.",
				eventMsg.Title,
				eventMsg.UserID,
				eventMsg.DateTime.Format(time.RFC822),
			))
		}
	}
}
