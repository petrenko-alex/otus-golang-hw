package memorystorage_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	event := storage.Event{
		Title:       "some event",
		DateTime:    time.Now(),
		Description: "this is some event",
		Duration:    time.Minute * 60,
		RemindTime:  time.Minute * 15,

		UserId: 1,
	}

	t.Run("create", func(t *testing.T) {
		memStorage := memorystorage.New()

		_, err := memStorage.Create(event)
		storageEvents := memStorage.ReadAll()

		require.NoError(t, err)
		require.Len(t, storageEvents, 1)
	})

	t.Run("update", func(t *testing.T) {
		newTitle := "new title"
		memStorage := memorystorage.New()

		// create to init storage
		id, createErr := memStorage.Create(event)
		require.NoError(t, createErr)

		// read & update
		event, readErr := memStorage.ReadOne(id)
		require.NoError(t, readErr)

		event.Title = newTitle
		updateErr := memStorage.Update(event)
		require.NoError(t, updateErr)

		// assert
		event, readErr = memStorage.ReadOne(id)
		require.NoError(t, readErr)
		require.Equal(t, newTitle, event.Title)
	})

	t.Run("update unknown", func(t *testing.T) {
		memStorage := memorystorage.New()

		// create to init storage
		_, createErr := memStorage.Create(event)
		require.NoError(t, createErr)

		// generate random ID & try to update
		event.ID = uuid.New().String()
		updateErr := memStorage.Update(event)

		// assert
		require.ErrorIs(t, updateErr, storage.ErrEventNotFound)

	})

	t.Run("delete", func(t *testing.T) {
		memStorage := memorystorage.New()

		// create to init storage
		id, createErr := memStorage.Create(event)
		require.NoError(t, createErr)

		// delete
		deleteErr := memStorage.Delete(id)

		// assert
		require.NoError(t, deleteErr)
		require.Len(t, memStorage.ReadAll(), 0)
	})

	t.Run("read unknown", func(t *testing.T) {
		memStorage := memorystorage.New()

		// create to init storage
		_, createErr := memStorage.Create(event)
		require.NoError(t, createErr)

		// generate random ID & try to read
		id := uuid.New().String()
		_, readErr := memStorage.ReadOne(id)

		// assert
		require.ErrorIs(t, readErr, storage.ErrEventNotFound)
	})

	t.Run("read all", func(t *testing.T) {
		memStorage := memorystorage.New()
		n := 3

		// create to init storage
		for i := 0; i < n; i++ {
			_, createErr := memStorage.Create(event)
			require.NoError(t, createErr)
		}

		require.Len(t, memStorage.ReadAll(), n)

	})

}
