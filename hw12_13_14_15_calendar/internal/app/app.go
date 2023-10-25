package app

import (
	"errors"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage storage.Storage
	logger  Logger
}

var (
	ErrNotFound      = errors.New("event not found")
	ErrDateBusy      = errors.New("time not available")
	ErrEventIsActive = errors.New("can't modify active event")
)

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

// CreateEvent create event if requested time is not busy.
func (a *App) CreateEvent(event entity.Event) (string, error) {
	existingEvent, getErr := a.storage.GetForTime(event.DateTime)
	if existingEvent != nil {
		return "", ErrDateBusy
	}

	if getErr != nil && !errors.Is(getErr, entity.ErrEventNotFound) {
		a.logger.Error(getErr.Error())

		return "", getErr
	}

	id, createErr := a.storage.Create(event)
	if createErr != nil {
		a.logger.Error(createErr.Error())

		return "", createErr
	}

	return id, nil
}

// UpdateEvent updates event if it is not active and requested time is not busy.
func (a *App) UpdateEvent(id string, event entity.Event) error {
	// check has event
	existingEvent, readErr := a.storage.GetByID(id)
	if readErr != nil {
		a.logger.Error(readErr.Error())

		if errors.Is(readErr, entity.ErrEventNotFound) {
			return ErrNotFound
		}
	}

	// check not active
	if existingEvent.DateTime.Round(time.Minute) == time.Now().Round(time.Minute) {
		return ErrEventIsActive
	}

	// check new time not busy
	eventWithRequestedTime, getErr := a.storage.GetForTime(event.DateTime)
	if eventWithRequestedTime != nil {
		return ErrDateBusy
	}

	if getErr != nil && !errors.Is(getErr, entity.ErrEventNotFound) {
		a.logger.Error(getErr.Error())

		return getErr
	}

	// update
	updateErr := a.storage.Update(event)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

// DeleteEvent deletes event if it is not active.
func (a *App) DeleteEvent(id string) error {
	event, readErr := a.storage.GetByID(id)
	if readErr != nil {
		a.logger.Error(readErr.Error())

		return readErr
	}

	if event.DateTime.Round(time.Minute) == time.Now().Round(time.Minute) {
		return ErrEventIsActive
	}

	deleteErr := a.storage.Delete(id)
	if deleteErr != nil {
		a.logger.Error(deleteErr.Error())

		return deleteErr
	}

	return nil
}

// GetDayEvents returns events for passed day. Use UTC time format.
func (a *App) GetDayEvents(day time.Time) (*entity.Events, error) {
	events, err := a.storage.GetForPeriod(
		day.Truncate(time.Hour*24),
		day.Round(time.Hour*24),
	)
	if err != nil {
		a.logger.Error(err.Error())

		return nil, err
	}

	return events, nil
}

// GetWeekEvents returns events for week starts with weekStart. Use UTC time format.
func (a *App) GetWeekEvents(weekStart time.Time) (*entity.Events, error) {
	weekStart = weekStart.Truncate(time.Hour * 24)
	weekEnd := weekStart.Add(time.Hour * 24 * 7)

	events, err := a.storage.GetForPeriod(weekStart, weekEnd)
	if err != nil {
		a.logger.Error(err.Error())

		return nil, err
	}

	return events, nil
}

// GetMonthEvents returns events for month starts with monthStart. Use UTC time format.
func (a *App) GetMonthEvents(monthStart time.Time) (*entity.Events, error) {
	monthStart = monthStart.Truncate(time.Hour * 24)
	monthEnd := time.Date(monthStart.Year(), monthStart.Month()+1, monthStart.Day(), 0, 0, 0, 0, monthStart.Location())

	events, err := a.storage.GetForPeriod(monthStart, monthEnd)
	if err != nil {
		a.logger.Error(err.Error())

		return nil, err
	}

	return events, nil
}
