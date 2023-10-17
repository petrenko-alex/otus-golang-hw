package app_test

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func createApp(t *testing.T) *app.App {
	t.Helper()

	strg, err := storage.GetStorage(string(storage.Memory))
	if err != nil {
		log.Fatal(err)
	}

	storageInitErr := strg.Connect(context.Background())
	if storageInitErr != nil {
		log.Fatal(storageInitErr)
	}

	return app.New(
		logger.New(logger.Debug, io.Discard),
		strg,
	)
}

func TestApp_CreateEvent(t *testing.T) {
	application := createApp(t)
	dateTime := time.Date(2023, 7, 8, 12, 0, 0, 0, time.UTC)

	id1, err1 := application.CreateEvent(
		entity.Event{
			Title:       "event 1",
			Description: "this is event 1",
			DateTime:    dateTime,
			Duration:    "02:00:00",
			RemindTime:  "00:15:00",
			UserId:      1,
		},
	)

	require.NoError(t, err1)
	require.NotEmpty(t, id1)

	_, err2 := application.CreateEvent(
		entity.Event{
			Title:       "event 2",
			Description: "this is event 2",
			DateTime:    dateTime,
			Duration:    "03:00:00",
			RemindTime:  "00:30:00",
			UserId:      1,
		},
	)

	require.ErrorIs(t, err2, app.ErrDateBusy)
}

func TestApp_UpdateEvent(t *testing.T) {
	application := createApp(t)
	dateTime := time.Date(2023, 7, 8, 12, 0, 0, 0, time.UTC)
	event := entity.Event{
		Title:       "event 1",
		Description: "this is event 1",
		DateTime:    dateTime,
		Duration:    "02:00:00",
		RemindTime:  "00:15:00",
		UserId:      1,
	}

	// update unknown
	updateErr := application.UpdateEvent("random id", event)
	require.ErrorIs(t, updateErr, app.ErrNotFound)

	// fill storage
	id1, err1 := application.CreateEvent(event)
	require.NoError(t, err1)
	event2 := event
	event2.DateTime = dateTime.Add(time.Hour * 24)
	_, err2 := application.CreateEvent(event2)
	require.NoError(t, err2)

	// update to busy date
	event.DateTime = event2.DateTime
	updateErr = application.UpdateEvent(id1, event)
	require.ErrorIs(t, updateErr, app.ErrDateBusy)

	// successful update
	event.ID = id1
	event.DateTime = time.Now()
	updateErr = application.UpdateEvent(id1, event)
	require.NoError(t, updateErr)

	// update active
	event.Title = "event 2"
	updateErr = application.UpdateEvent(id1, event)
	require.ErrorIs(t, updateErr, app.ErrEventIsActive)
}

func TestApp_DeleteEvent(t *testing.T) {
	application := createApp(t)
	dateTime := time.Date(2023, 7, 8, 12, 0, 0, 0, time.UTC)

	id1, err1 := application.CreateEvent(
		entity.Event{
			Title:       "event 1",
			Description: "this is event 1",
			DateTime:    dateTime,
			Duration:    "02:00:00",
			RemindTime:  "00:15:00",
			UserId:      1,
		},
	)
	require.NoError(t, err1)

	id2, err2 := application.CreateEvent(
		entity.Event{
			Title:       "event 2",
			Description: "this is event 2",
			DateTime:    time.Now(),
			Duration:    "03:00:00",
			RemindTime:  "00:30:00",
			UserId:      1,
		},
	)
	require.NoError(t, err2)

	deleteErr1 := application.DeleteEvent(id1)
	require.NoError(t, deleteErr1)

	deleteErr2 := application.DeleteEvent(id2)
	require.ErrorIs(t, deleteErr2, app.ErrEventIsActive)
}
