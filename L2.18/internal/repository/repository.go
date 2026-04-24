package repository

import (
	"calendar/internal/domain"
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]domain.Event
}

func NewRepository() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]domain.Event),
	}
}

func (s *Storage) Save(e domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[e.ID] = e

	// В данном случае error в функции не нужен, но оставил
	// чтобы можно было заменить реализацию хранилиза на бд
	// и не пришлось менять код в сервисе и хендлере
	return nil
}

func (s *Storage) Update(e domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[e.ID]; !exists {
		return domain.ErrEventNotFound
	}

	s.events[e.ID] = e
	return nil
}

func (s *Storage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return domain.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) List() []domain.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]domain.Event, 0, len(s.events))
	for _, e := range s.events {
		res = append(res, e)
	}

	return res
}
