package main

import (
	"context"
	"errors"
	"flag"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/server"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()

		return 0
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	file, fileErr := os.Open(configFile)
	if fileErr != nil {
		log.Println("Error opening config file.")

		return 1
	}

	cfg, configErr := config.New(file)
	if configErr != nil {
		log.Println("Error parsing config file.")

		return 1
	}
	ctx = cfg.WithContext(ctx)

	logg := createLogger(cfg)

	appStorage, storageErr := storage.Get(cfg.Storage)
	if storageErr != nil {
		logg.Error("Error getting storage: " + storageErr.Error())

		return 1
	}

	storageInitErr := appStorage.Connect(ctx)
	if storageInitErr != nil {
		logg.Error("Error init storage: " + storageInitErr.Error())

		return 1
	}

	srv := server.NewServer(
		server.Options{
			GRPC: server.GRPCOptions{
				Host:           cfg.GRPCServer.Host,
				Port:           cfg.GRPCServer.Port,
				ConnectTimeout: cfg.GRPCServer.ConnectTimeout,
			},
			HTTP: server.HTTPOptions{
				Host:         cfg.HTTPServer.Host,
				Port:         cfg.HTTPServer.Port,
				ReadTimeout:  cfg.HTTPServer.ReadTimeout,
				WriteTimeout: cfg.HTTPServer.WriteTimeout,
			},
		},
		logg,
		app.New(logg, appStorage),
	)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		logg.Info("Starting GRPC server...")
		err := srv.Start(ctx)
		if err != nil {
			logg.Error("Failed to start GRPC server: " + err.Error())
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logg.Info("Starting HTTP server...")
		err := srv.InitAndStartHTTPProxy(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error("Failed to start HTTP server: " + err.Error())
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("Stopping GRPC server...")
		if err := srv.Stop(ctx); err != nil {
			logg.Error("Failed to stop GRPC server: " + err.Error())
		}

		logg.Info("Stopping HTTP server...")
		if err := srv.StopHTTPProxy(ctx); err != nil {
			logg.Error("Failed to stop HTTP server: " + err.Error())
		}

		logg.Info("Calendar stopped.")
		wg.Done()
	}()

	wg.Wait()

	return 0
}

func createLogger(cfg *config.Config) logger.Logger {
	levelMap := map[string]logger.Level{
		"debug":   logger.Debug,
		"info":    logger.Info,
		"warning": logger.Warning,
		"error":   logger.Error,
	}

	return logger.New(levelMap[cfg.Logger.Level], os.Stdout)
}
