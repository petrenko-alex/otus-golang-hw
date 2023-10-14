package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yml", "Path to configuration file")
}

// TODO: separate service to get from context so not to hardcode context value keys

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	file, fileErr := os.Open(configFile)
	if fileErr != nil {
		log.Fatal("Error opening config file.")
	}

	cfg, configErr := config.NewConfig(file)
	if configErr != nil {
		log.Fatal("Error parsing config file.")
	}

	logg := createLogger(cfg)

	appStorage, storageErr := storage.GetStorage(cfg.Storage)
	if storageErr != nil {
		log.Fatal("Error getting storage.")
	}

	server := internalhttp.NewServer(
		internalhttp.ServerOptions{
			Host:         cfg.Server.Host,
			Port:         cfg.Server.Port,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
		logg,
		app.New(logg, appStorage),
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		logg.Info("calendar stopped.")
		wg.Done()
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	wg.Wait()
}

func createLogger(cfg *config.Config) app.Logger {
	levelMap := map[string]logger.Level{
		"debug":   logger.Debug,
		"info":    logger.Info,
		"warning": logger.Warning,
		"error":   logger.Error,
	}

	return logger.New(levelMap[cfg.Logger.Level], os.Stdout)
}
