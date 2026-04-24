package http

import (
	"calendar/internal/domain"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type CalendarService interface {
	CreateEvent(e domain.Event) error
	UpdateEvent(e domain.Event) error
	DeleteEvent(id uuid.UUID) error

	GetEventsForDay(userID int, date time.Time) ([]domain.Event, error)
	GetEventsForWeek(userID int, date time.Time) ([]domain.Event, error)
	GetEventsForMonth(userID int, date time.Time) ([]domain.Event, error)
}

type Handler struct {
	service CalendarService
	logger  *slog.Logger
}

func NewHandler(service CalendarService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
