package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/queue"
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

	msgChan, consumeErr := queueManager.Consume(q.Name)
	if consumeErr != nil {
		logg.Error("Error registering consumer: " + consumeErr.Error())

		return 1
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for message := range msgChan {
			eventMsg := entity.EventMsg{}

			unmarshalErr := json.Unmarshal(message.Body, &eventMsg)
			if unmarshalErr != nil {
				logg.Error("Error reading msg from RabbitMQ: " + unmarshalErr.Error())

				continue
			}

			logg.Info(fmt.Sprintf(
				"Sending reminder about \"%s\" event to #%d user. Event time: %s.",
				eventMsg.Title,
				eventMsg.UserId,
				eventMsg.DateTime.Format(time.RFC822),
			))
		}

		wg.Done()
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
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
