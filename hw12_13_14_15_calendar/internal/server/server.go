package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	proto "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/server/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server          *grpc.Server
	httpProxyServer *http.Server

	logger logger.Logger

	options Options
}

type Application interface {
	CreateEvent(event entity.Event) (string, error)
	UpdateEvent(id string, event entity.Event) error
	DeleteEvent(id string) error
	GetDayEvents(day time.Time) (*entity.Events, error)
	GetWeekEvents(weekStart time.Time) (*entity.Events, error)
	GetMonthEvents(monthStart time.Time) (*entity.Events, error)
}

func NewServer(options Options, logger logger.Logger, app Application) *Server {
	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(options.GRPC.ConnectTimeout),
		grpc.ChainUnaryInterceptor(
			log.NewInterceptor(logger).GetInterceptor(),
		),
	)
	proto.RegisterEventServiceServer(grpcServer, NewService(app, logger))

	return &Server{
		server:  grpcServer,
		logger:  logger,
		options: options,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var lc net.ListenConfig

	listener, err := lc.Listen(
		ctx,
		"tcp",
		net.JoinHostPort(s.options.GRPC.Host, s.options.GRPC.Port),
	)
	if err != nil {
		return err
	}

	err = s.server.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()

	return nil
}

func (s *Server) InitAndStartHttpProxy(ctx context.Context) error {
	conn, err := grpc.DialContext(
		ctx,
		net.JoinHostPort(s.options.GRPC.Host, s.options.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	mux := runtime.NewServeMux()
	err = proto.RegisterEventServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}

	s.httpProxyServer = &http.Server{
		Addr:    net.JoinHostPort(s.options.HTTP.Host, s.options.HTTP.Port),
		Handler: log.NewHandler(s.logger, mux),
	}

	err = s.httpProxyServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) StopHttpProxy(ctx context.Context) error {
	if s.httpProxyServer == nil {
		return nil
	}

	err := s.httpProxyServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
