package app

import (
	"context"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface {
	Create(storage.Event) (string, error)
	ReadOne(string) (storage.Event, error)
	ReadAll() storage.Events
	Update(storage.Event) error
	Delete(string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
