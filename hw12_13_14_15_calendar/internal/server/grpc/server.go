package grpc

import (
	"context"
	proto "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	server *grpc.Server
	logger Logger

	host, port string
}

type ServerOptions struct {
	Host, Port     string
	ConnectTimeout time.Duration
}

type Logger interface {
	app.Logger
}

type Application interface {
	CreateEvent(event entity.Event) (string, error)
	UpdateEvent(id string, event entity.Event) error
	DeleteEvent(id string) error
	GetDayEvents(day time.Time) (*entity.Events, error)
	GetWeekEvents(weekStart time.Time) (*entity.Events, error)
	GetMonthEvents(monthStart time.Time) (*entity.Events, error)
}

func NewServer(options ServerOptions, logger Logger, app Application) *Server {
	logInterceptor := NewLogHandler(logger)

	server := grpc.NewServer(
		grpc.ConnectionTimeout(options.ConnectTimeout),
		grpc.ChainUnaryInterceptor(
			logInterceptor.GetInterceptor(),
		),
	)
	proto.RegisterEventServiceServer(server, NewAppHandler(app, logger))

	return &Server{
		server: server,
		logger: logger,
		host:   options.Host,
		port:   options.Port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var lc net.ListenConfig

	listener, err := lc.Listen(ctx, "tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	err = s.server.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.GracefulStop()

	return nil
}
