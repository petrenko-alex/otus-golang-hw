package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	AppHandler

	server http.Server
	logger Logger
}

type ServerOptions struct {
	Host, Port                string
	ReadTimeout, WriteTimeout time.Duration
}

type Logger interface {
	app.Logger
}

type Application interface { // TODO
}

func NewServer(options ServerOptions, logger Logger, app Application) *Server {
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

func (s *Server) Start(ctx context.Context) error {
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
