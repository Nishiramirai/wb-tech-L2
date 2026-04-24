package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID     uuid.UUID
	UserID int
	Date   time.Time
	Text   string
}
