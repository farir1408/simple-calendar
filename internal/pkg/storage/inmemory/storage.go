package inmemory

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/farir1408/simple-calendar/internal/pkg/types"
)

const (
	defaultCapacity = 10
)

// Storage ...
type Storage struct {
	sequenceID int64
	mu         sync.Mutex
	events     map[uint64]types.Event
}

// NewInMemoryStorage ...
func NewInMemoryStorage(capacity int) *Storage {
	if capacity == 0 {
		capacity = defaultCapacity
	}
	return &Storage{
		sequenceID: 0,
		events:     make(map[uint64]types.Event, capacity),
	}
}

// SaveEvent ...
func (s *Storage) SaveEvent(ctx context.Context, event types.Event) (uint64, error) {
	id := s.generateID()
	event.ID = id
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[id] = event
	return id, nil
}

// GetEvent ...
func (s *Storage) GetEvent(ctx context.Context, id uint64) (types.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	event, ok := s.events[id]
	if !ok {
		return types.Event{}, errors.New("event not found")
	}
	return event, nil
}

// UpdateEvent ...
func (s *Storage) UpdateEvent(ctx context.Context, event types.Event) error {
	_, ok := s.events[event.ID]
	if !ok {
		return errors.New("event not found")
	}
	s.events[event.ID] = event
	return nil
}

// DeleteEvent ...
func (s *Storage) DeleteEvent(ctx context.Context, id uint64) error {
	_, ok := s.events[id]
	if !ok {
		return errors.New("event not found")
	}

	delete(s.events, id)
	return nil
}

// Close ...
func (s *Storage) Close() error {
	return nil
}

func (s *Storage) generateID() uint64 {
	id := atomic.AddInt64(&s.sequenceID, 1)
	return uint64(id)
}
