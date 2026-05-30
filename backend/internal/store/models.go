package store

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DisplayName string
	CreatedAt   time.Time
}

type CreateUserParams struct {
	Email       string
	DisplayName string
}

type MagicLink struct {
	ID        uuid.UUID
	UserID    *uuid.UUID
	Email     string
	Token     string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type CreateMagicLinkParams struct {
	Email     string
	UserID    *uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type Game struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Status      string
	CreatedAt   time.Time
}

type QueueEntry struct {
	ID       uuid.UUID
	GameID   uuid.UUID
	UserID   uuid.UUID
	Status   string
	JoinedAt time.Time
}

type Session struct {
	ID        uuid.UUID
	GameID    uuid.UUID
	Status    string
	StartedAt time.Time
	EndedAt   *time.Time
}

type DigitalGood struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Category    *string
	GameID      *uuid.UUID
	CreatedAt   time.Time
}

type InventoryItem struct {
	Good        DigitalGood
	Quantity    int
	AcquiredAt  time.Time
}
