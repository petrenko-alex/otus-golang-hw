package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/scheduler"
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

	// Init RabbitMQ
	queueManager := queue.NewRabbitManager(ctx, *cfg, logg)
	connectErr := queueManager.Connect()
	if connectErr != nil {
		logg.Error("Error connecting to RabbitMQ server: " + connectErr.Error())

		return 1
	}

	q, queueErr := queueManager.CreateQueue(cfg.App.Scheduler.Queue)
	if queueErr != nil {
		logg.Error("Error declaring queue: " + queueErr.Error())

		return 1
	}

	// Init storage
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

	eventScheduler := scheduler.New(
		cfg.App.Scheduler.Period,
		q.Name,
		ctx,
		logg,
		appStorage,
		queueManager,
	)
	eventScheduler.Run()

	closeErr := queueManager.Close()
	if closeErr != nil {
		logg.Error("Error closing RabbitMQ connection: " + closeErr.Error())

		return 1
	}

	logg.Info("Connection closed.")

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
