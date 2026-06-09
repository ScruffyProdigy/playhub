package store

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID
	GameID      uuid.UUID
	ModeID      *uuid.UUID
	ModeQueueID *uuid.UUID
	Status      string
	StartedAt   time.Time
	EndedAt     *time.Time
}
