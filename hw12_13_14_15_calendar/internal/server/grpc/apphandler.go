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

func (h AppHandler) entities2Proto(entityEvents *entity.Events) *proto.Events {
	protoEvents := make([]*proto.Event, 0, len(*entityEvents))

	for _, entityEvent := range *entityEvents {
		protoEvents = append(
			protoEvents,
			h.entity2Proto(entityEvent),
		)
	}

	return &proto.Events{Events: protoEvents}
}

func (h AppHandler) entity2Proto(entityEvent entity.Event) *proto.Event {
	return &(proto.Event{
		Id:          entityEvent.ID,
		Title:       entityEvent.Title,
		DateTime:    timestamppb.New(entityEvent.DateTime),
		Description: entityEvent.Description,
		Duration:    entityEvent.Duration,
		RemindTime:  entityEvent.RemindTime,
	})
}
