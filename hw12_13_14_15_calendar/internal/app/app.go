package app

import (
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
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
