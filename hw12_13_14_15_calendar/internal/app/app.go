package app

import (
	"context"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage storage.Storage
	logger  Logger
}

type Logger interface {
	Debug(string)
	Info(string)
	Warning(string)
	Error(string)
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
