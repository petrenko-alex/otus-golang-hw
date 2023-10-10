package memorystorage

import (
	"sync"

	"github.com/google/uuid"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu   sync.RWMutex //nolint:unused
	data map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		data: make(map[string]storage.Event, 0),
	}
}

func (s *Storage) ReadOne(id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, has := s.data[id]
	if !has {
		return nil, storage.ErrEventNotFound
	}

	return &event, nil
}

func (s *Storage) ReadAll() (*storage.Events, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make(storage.Events, 0, len(s.data))

	for _, event := range s.data {
		events = append(events, event)
	}

	return &events, nil
}

func (s *Storage) Create(event storage.Event) (string, error) {
	event.ID = uuid.New().String()

	s.mu.Lock()
	s.data[event.ID] = event
	s.mu.Unlock()

	return event.ID, nil
}

func (s *Storage) Update(event storage.Event) error {
	_, err := s.ReadOne(event.ID)
	if err != nil {
		return err // todo: wrap?
	}

	s.mu.Lock()
	s.data[event.ID] = event
	s.mu.Unlock()

	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()

	return nil
}
