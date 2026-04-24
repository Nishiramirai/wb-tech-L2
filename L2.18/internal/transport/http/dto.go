package http

import (
	"calendar/internal/domain"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type eventBase struct {
	UserID    int    `json:"user_id"`
	Date      string `json:"date"`
	EventText string `json:"event"`
}

func (eb *eventBase) validate() (time.Time, error) {
	var errs []string

	if eb.UserID <= 0 {
		errs = append(errs, "user_id is required and must be > 0")
	}
	if eb.EventText == "" {
		errs = append(errs, "event text cannot be empty")
	}

	var parsedDate time.Time
	if eb.Date == "" {
		errs = append(errs, "date is required")
	} else {
		var err error
		parsedDate, err = time.Parse("2006-01-02", eb.Date)
		if err != nil {
			errs = append(errs, "invalid date format, expected YYYY-MM-DD")
		}
	}

	if len(errs) > 0 {
		return time.Time{}, fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}

	return parsedDate, nil
}

// Create event
type createEventRequest struct {
	eventBase
}

func (req *createEventRequest) ToDomain() (domain.Event, error) {
	parsedDate, err := req.validate()
	if err != nil {
		return domain.Event{}, err
	}

	return domain.Event{
		ID:     uuid.New(),
		UserID: req.UserID,
		Date:   parsedDate,
		Text:   req.EventText,
	}, nil
}

// Update Event
type updateEventRequest struct {
	ID string `json:"id"`
	eventBase
}

func (req *updateEventRequest) ToDomain() (domain.Event, error) {
	uid, err := uuid.Parse(req.ID)
	if err != nil {
		return domain.Event{}, errors.New("invalid event id format")
	}

	parsedDate, err := req.validate()
	if err != nil {
		return domain.Event{}, err
	}

	return domain.Event{
		ID:     uid,
		UserID: req.UserID,
		Date:   parsedDate,
		Text:   req.EventText,
	}, nil
}

// Delete Event
type deleteEventRequest struct {
	ID string `json:"id"`
}
