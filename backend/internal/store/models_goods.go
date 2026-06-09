package store

import (
	"time"

	"github.com/google/uuid"
)

type DigitalGood struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Category    *string
	GameID      *uuid.UUID
	CreatedAt   time.Time
}

type InventoryItem struct {
	Good       DigitalGood
	Quantity   int
	AcquiredAt time.Time
}
