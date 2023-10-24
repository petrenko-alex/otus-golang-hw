package grpc

import (
	"context"
	proto "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AppHandler struct {
	proto.UnimplementedEventServiceServer

	app    Application
	logger logger.Logger
}

func NewAppHandler(app Application, logger logger.Logger) *AppHandler {
	return &AppHandler{
		app:    app,
		logger: logger,
	}
}

func (h AppHandler) CreateEvent(ctx context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	id, err := h.app.CreateEvent(h.proto2entity(request.GetEventData()))
	if err != nil {
		h.logger.Error(err.Error())

		return nil, err
	}

	return &proto.CreateResponse{EventId: &(proto.EventId{Id: id})}, nil
}

func (h AppHandler) UpdateEvent(ctx context.Context, request *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	err := h.app.UpdateEvent(
		request.GetEventId().GetId(),
		h.proto2entity(request.GetEventData()),
	)
	if err != nil {
		h.logger.Error(err.Error())

		return nil, err
	}

	return &proto.UpdateResponse{}, nil
}

func (h AppHandler) DeleteEvent(ctx context.Context, request *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	err := h.app.DeleteEvent(request.GetEventId().GetId())
	if err != nil {
		h.logger.Error(err.Error())

		return nil, err
	}

	return &proto.DeleteResponse{}, nil
}

func (h AppHandler) GetWeekEvents(ctx context.Context, date *proto.StartDate) (*proto.Events, error) {
	events, err := h.app.GetWeekEvents(date.GetStartDate().AsTime())
	if err != nil {
		h.logger.Error(err.Error())

		return nil, err
	}

	return h.entities2Proto(events), nil
}

func (h AppHandler) GetMonthEvents(ctx context.Context, date *proto.StartDate) (*proto.Events, error) {
	events, err := h.app.GetMonthEvents(date.GetStartDate().AsTime())
	if err != nil {
		h.logger.Error(err.Error())

		return nil, err
	}

	return h.entities2Proto(events), nil
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
		EventId: &proto.EventId{Id: entityEvent.ID},
		EventData: &proto.EventData{
			Title:       entityEvent.Title,
			DateTime:    timestamppb.New(entityEvent.DateTime),
			Description: entityEvent.Description,
			Duration:    entityEvent.Duration,
			RemindTime:  entityEvent.RemindTime,
		},
	})
}

func (h AppHandler) proto2entity(protoEvent *proto.EventData) entity.Event {
	return entity.Event{
		Title:       protoEvent.Title,
		Description: protoEvent.Description,
		DateTime:    protoEvent.DateTime.AsTime(),
		Duration:    protoEvent.Duration,
		RemindTime:  protoEvent.RemindTime,
		UserID:      int(protoEvent.UserId),
	}
}
