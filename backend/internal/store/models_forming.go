package store

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	PartyStatusWaiting   = "waiting"
	PartyStatusPlaced    = "placed"
	PartyStatusMatched   = "matched"
	PartyStatusCancelled = "cancelled"
)

type Party struct {
	ID           uuid.UUID
	LeaderUserID *uuid.UUID
	ModeQueueID  uuid.UUID
	Status       string
	PartyTree    json.RawMessage
	CreatedAt    time.Time
}

const (
	FormingMatchStatusFilling = "filling"
	FormingMatchStatusFired   = "fired"
)

type FormingMatch struct {
	ID          uuid.UUID
	ModeQueueID uuid.UUID
	ModeID      uuid.UUID
	GameID      uuid.UUID
	Status      string
	TargetSpec  json.RawMessage
	CreatedAt   time.Time
	FiredAt     *time.Time
}

type FormingAssignment struct {
	ID             uuid.UUID
	FormingMatchID uuid.UUID
	SeatKey        string
	UserID         *uuid.UUID
	PartyID        *uuid.UUID
	QueuePath      string
	AffinityKey    *string
	Source         string
	TableID        *uuid.UUID
	AssignedAt     time.Time
}
