package store

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	DisplayName  string
	AvatarURL    *string
	AvatarKey    *string
	AvatarSource *string
	IsGuest      bool
	CreatedAt    time.Time
}

type CreateUserParams struct {
	Email       string
	DisplayName string
}

type MagicLink struct {
	ID             uuid.UUID
	UserID         *uuid.UUID
	Email          string
	TokenHash      string
	FailedAttempts int
	CodeHash       string
	ExpiresAt      time.Time
	UsedAt         *time.Time
	CreatedAt      time.Time
}

type CreateMagicLinkParams struct {
	Email     string
	UserID    *uuid.UUID
	TokenHash string
	CodeHash  string
	ExpiresAt time.Time
}
