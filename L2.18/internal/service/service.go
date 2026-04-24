package service

import (
	"calendar/internal/domain"
	"time"

	"github.com/google/uuid"
)

type EventRepository interface {
	Save(e domain.Event) error
	Update(e domain.Event) error
	Delete(id uuid.UUID) error
	List() []domain.Event
}

type CalendarService struct {
	repo EventRepository
}

func NewCalendarService(repo EventRepository) *CalendarService {
	return &CalendarService{
		repo: repo,
	}
}

func (s *CalendarService) CreateEvent(e domain.Event) error {
	return s.repo.Save(e)
}

func (s *CalendarService) UpdateEvent(e domain.Event) error {
	return s.repo.Update(e)
}

func (s *CalendarService) DeleteEvent(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *CalendarService) GetEventsForDay(userID int, date time.Time) ([]domain.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1).Add(-time.Nanosecond)

	return s.getEventsForPeriod(userID, start, end)
}

func (s *CalendarService) GetEventsForWeek(userID int, date time.Time) ([]domain.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)

	return s.getEventsForPeriod(userID, start, end)
}

func (s *CalendarService) GetEventsForMonth(userID int, date time.Time) ([]domain.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return s.getEventsForPeriod(userID, start, end)
}

func (s *CalendarService) getEventsForPeriod(userID int, start, end time.Time) ([]domain.Event, error) {
	allEvents := s.repo.List()
	filtered := make([]domain.Event, 0)

	for _, e := range allEvents {
		if e.UserID == userID {
			if (e.Date.After(start) || e.Date.Equal(start)) && (e.Date.Before(end) || e.Date.Equal(end)) {
				filtered = append(filtered, e)
			}
		}
	}
	return filtered, nil
}
