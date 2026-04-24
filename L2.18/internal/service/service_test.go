package service

import (
	"calendar/internal/domain"
	"calendar/internal/repository"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCalendarService(t *testing.T) {
	repo := repository.NewRepository()
	svc := NewCalendarService(repo)

	userID := 1
	today := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC)

	t.Run("Create and Get Day", func(t *testing.T) {
		event := domain.Event{
			ID:     uuid.New(),
			UserID: userID,
			Date:   today,
			Text:   "Test Event",
		}

		// Тестируем создание
		err := svc.CreateEvent(event)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Тестируем выборку за день
		events, err := svc.GetEventsForDay(userID, today)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(events) != 1 {
			t.Errorf("expected 1 event, got %d", len(events))
		}
	})

	t.Run("Filter by user", func(t *testing.T) {
		// Событие другого пользователя
		eventOther := domain.Event{
			ID:     uuid.New(),
			UserID: 999,
			Date:   today,
			Text:   "Other User Event",
		}
		_ = svc.CreateEvent(eventOther)

		events, _ := svc.GetEventsForDay(userID, today)
		if len(events) != 1 {
			t.Errorf("expected only 1 event for user %d, but filtered incorrectly", userID)
		}
	})

	t.Run("Week range check", func(t *testing.T) {
		// Событие через 5 дней (в пределах недели)
		eventInWeek := domain.Event{
			ID:     uuid.New(),
			UserID: userID,
			Date:   today.AddDate(0, 0, 5),
			Text:   "In week",
		}
		// Событие через 10 дней (вне недели)
		eventOutWeek := domain.Event{
			ID:     uuid.New(),
			UserID: userID,
			Date:   today.AddDate(0, 0, 10),
			Text:   "Out of week",
		}

		_ = svc.CreateEvent(eventInWeek)
		_ = svc.CreateEvent(eventOutWeek)

		events, _ := svc.GetEventsForWeek(userID, today)
		// Должно быть 2: первое (сегодня) и второе (через 5 дней)
		if len(events) != 2 {
			t.Errorf("expected 2 events in week range, got %d", len(events))
		}
	})
}
