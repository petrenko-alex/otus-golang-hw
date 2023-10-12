package memorystorage_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	memorystorage "github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	event := entity.Event{
		Title:       "some event",
		DateTime:    time.Now(),
		Description: "this is some event",
		Duration:    "60",
		RemindTime:  "15",

		UserId: 1,
	}

	t.Run("create", func(t *testing.T) {
		memStorage := memorystorage.New()

		_, createErr := memStorage.Create(event)
		storageEvents, _ := memStorage.ReadAll() // TODO: убрать это? читать сыро. Иначе, тестируем код с помощью того же кода, который тестируем.

		require.NoError(t, createErr)
		require.Len(t, *storageEvents, 1)
	})

	t.Run("update", func(t *testing.T) {
		newTitle := "new title"
		memStorage := memorystorage.New()

		// create to init storage
		id, createErr := memStorage.Create(event)
		require.NoError(t, createErr)

		// read & update
		event1, readErr := memStorage.ReadOne(id)
		require.NoError(t, readErr)

		event1.Title = newTitle
		updateErr := memStorage.Update(*event1)
		require.NoError(t, updateErr)

		// assert
		event2, readErr := memStorage.ReadOne(id)
		require.NoError(t, readErr)
		require.Equal(t, newTitle, event2.Title)
	})

	t.Run("update unknown", func(t *testing.T) {
		memStorage := memorystorage.New()

		// create to init storage
		_, createErr := memStorage.Create(event) // фикстуры?
		require.NoError(t, createErr)

		// generate random ID & try to update
		event.ID = uuid.New().String()
		updateErr := memStorage.Update(event)

		// assert
		require.ErrorIs(t, updateErr, entity.ErrEventNotFound)
	})

	t.Run("delete", func(t *testing.T) {
		memStorage := memorystorage.New()

		// create to init storage
		id, createErr := memStorage.Create(event)
		require.NoError(t, createErr)

		// delete
		deleteErr := memStorage.Delete(id)

		// assert
		events, _ := memStorage.ReadAll()
		require.NoError(t, deleteErr)
		require.Len(t, *events, 0)
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
		require.ErrorIs(t, readErr, entity.ErrEventNotFound)
	})

	t.Run("read all", func(t *testing.T) {
		memStorage := memorystorage.New()
		n := 3

		// create to init storage
		for i := 0; i < n; i++ {
			_, createErr := memStorage.Create(event)
			require.NoError(t, createErr)
		}

		events, _ := memStorage.ReadAll()
		require.Len(t, *events, n)
	})
}
