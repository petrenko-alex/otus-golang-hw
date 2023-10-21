package grpc

import (
	"context"
	proto "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AppHandler struct {
	proto.UnimplementedEventServiceServer

	app    Application
	logger Logger
}

func NewAppHandler(app Application, logger Logger) *AppHandler {
	return &AppHandler{
		app:    app,
		logger: logger,
	}
}

func (h AppHandler) GetDayEvents(ctx context.Context, date *proto.StartDate) (*proto.Events, error) {
	events, err := h.app.GetDayEvents(date.GetStartDate().AsTime())
	if err != nil {
		h.logger.Error(err.Error())

		return nil, err
	}

	return h.entities2Proto(events), nil

}

// todo: move to another package, fix naming
func (h AppHandler) entities2Proto(entities *entity.Events) *proto.Events {
	events := make([]*proto.Event, 0, len(*entities))

	for _, event := range *entities {
		events = append(
			events,
			h.entity2Proto(event),
		)
	}

	return &proto.Events{Events: events}
}

func (h AppHandler) entity2Proto(entity entity.Event) *proto.Event {
	event := proto.Event{
		Id:          entity.ID,
		Title:       entity.Title,
		DateTime:    timestamppb.New(entity.DateTime),
		Description: entity.Description,
		Duration:    entity.Duration,
		RemindTime:  entity.RemindTime,
	}

	return &event
}
