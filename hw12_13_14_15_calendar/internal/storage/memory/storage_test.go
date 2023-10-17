package memorystorage_test

import (
	"sort"
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
		storageEvents, _ := memStorage.GetAll() // TODO: убрать это? читать напрямую из sql. Иначе, тестируем код с помощью того же кода, который тестируем.

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
		event1, readErr := memStorage.GetById(id)
		require.NoError(t, readErr)

		event1.Title = newTitle
		updateErr := memStorage.Update(*event1)
		require.NoError(t, updateErr)

		// assert
		event2, readErr := memStorage.GetById(id)
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
		events, _ := memStorage.GetAll()
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
		_, readErr := memStorage.GetById(id)

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

		events, _ := memStorage.GetAll()
		require.Len(t, *events, n)
	})

	t.Run("read for day", func(t *testing.T) {
		initialDate := time.Date(2023, 10, 16, 13, 10, 0, 0, time.UTC)
		strg := memorystorage.NewWithEvents(map[string]entity.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate, UserId: 1},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(time.Hour * 2), UserId: 1},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(-time.Hour * 2), UserId: 1},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24), UserId: 1},
			"5": {ID: "5", Title: "5", DateTime: initialDate.Add(-time.Hour * 24), UserId: 1},
		})

		events, err := strg.GetForPeriod(
			time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC),
		)

		require.NoError(t, err)
		require.Len(t, *events, 3)
		require.Equal(t, []string{"1", "2", "3"}, getKeys(t, events))
	})

	t.Run("read for week", func(t *testing.T) {
		initialDate := time.Date(2023, 10, 16, 13, 10, 0, 0, time.UTC)
		strg := memorystorage.NewWithEvents(map[string]entity.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate, UserId: 1},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(time.Hour * 24 * 10), UserId: 1},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(-time.Hour * 24 * 10), UserId: 1},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24 * 2), UserId: 1},
		})

		events, err := strg.GetForPeriod(
			time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC),
		)

		require.NoError(t, err)
		require.Len(t, *events, 2)
		require.Equal(t, []string{"1", "4"}, getKeys(t, events))
	})

	t.Run("read for month", func(t *testing.T) {
		initialDate := time.Date(2023, 10, 16, 13, 10, 0, 0, time.UTC)
		strg := memorystorage.NewWithEvents(map[string]entity.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate, UserId: 1},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(time.Hour * 24 * 30), UserId: 1},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(time.Hour * 24 * 10), UserId: 1},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24 * 2), UserId: 1},
			"5": {ID: "5", Title: "5", DateTime: initialDate.Add(-time.Hour * 24 * 30), UserId: 1},
		})

		events, err := strg.GetForPeriod(
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
		)

		require.NoError(t, err)
		require.Len(t, *events, 3)
		require.Equal(t, []string{"1", "3", "4"}, getKeys(t, events))
	})
}

func getKeys(t *testing.T, events *entity.Events) []string {
	t.Helper()
	keys := make([]string, 0, len(*events))

	for _, event := range *events {
		keys = append(keys, event.ID)
	}

	sort.Strings(keys)

	return keys
}
