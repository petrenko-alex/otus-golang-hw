package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	AppHandler

	server http.Server
	logger logger.Logger
}

type ServerOptions struct {
	Host, Port                string
	ReadTimeout, WriteTimeout time.Duration
}

type Application interface {
	CreateEvent(event entity.Event) (string, error)
	UpdateEvent(id string, event entity.Event) error
	DeleteEvent(id string) error
	GetDayEvents(day time.Time) (*entity.Events, error)
	GetWeekEvents(weekStart time.Time) (*entity.Events, error)
	GetMonthEvents(monthStart time.Time) (*entity.Events, error)
}

func NewServer(options ServerOptions, logger logger.Logger, app Application) *Server {
	return &Server{
		server: http.Server{
			Addr:         net.JoinHostPort(options.Host, options.Port),
			Handler:      NewLogHandler(logger, NewAppHandler(app, logger)),
			ReadTimeout:  options.ReadTimeout,
			WriteTimeout: options.WriteTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
