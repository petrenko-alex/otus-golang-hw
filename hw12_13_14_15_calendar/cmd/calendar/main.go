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
		log.Println("Error getting storage.")

		return 1
	}

	storageInitErr := appStorage.Connect(ctx)
	if storageInitErr != nil {
		log.Println("Error init storage.")

		return 1
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

		return 1
	}

	wg.Wait()

	return 0
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
