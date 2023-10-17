package memorystorage

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
)

type Storage struct {
	mu   sync.RWMutex //nolint:unused
	data map[string]entity.Event
}

func New() *Storage {
	return &Storage{
		data: make(map[string]entity.Event),
	}
}

func NewWithEvents(events map[string]entity.Event) *Storage {
	return &Storage{data: events}
}

func (s *Storage) GetById(id string) (*entity.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, has := s.data[id]
	if !has {
		return nil, entity.ErrEventNotFound
	}

	return &event, nil
}

func (s *Storage) GetAll() (*entity.Events, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make(entity.Events, 0, len(s.data))

	for _, event := range s.data {
		events = append(events, event)
	}

	return &events, nil
}

func (s *Storage) Create(event entity.Event) (string, error) {
	event.ID = uuid.New().String()

	s.mu.Lock()
	s.data[event.ID] = event
	s.mu.Unlock()

	return event.ID, nil
}

func (s *Storage) Update(event entity.Event) error {
	_, err := s.GetById(event.ID)
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

func (s *Storage) GetForTime(t time.Time) (*entity.Event, error) {
	for _, event := range s.data {
		if event.DateTime == t {
			return &event, nil
		}
	}

	return nil, entity.ErrEventNotFound
}

func (s *Storage) GetForPeriod(periodStart time.Time, periodEnd time.Time) (*entity.Events, error) {
	periodEvents := make(entity.Events, 0)

	for _, event := range s.data {
		if event.DateTime.After(periodStart) && event.DateTime.Before(periodEnd) {
			periodEvents = append(periodEvents, event)
		}
	}

	return &periodEvents, nil
}
