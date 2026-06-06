package store

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DisplayName string
	AvatarURL   *string
	AvatarKey   *string
	AvatarSource *string
	CreatedAt   time.Time
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

type Game struct {
	ID               uuid.UUID
	Name             string
	Description      *string
	Slug             *string
	PlayURL          *string
	APIBaseURL       *string
	Status           string
	ManifestHash     *string
	ManifestETag     *string
	ManifestSyncedAt *time.Time
	GameVersion      *string
	WebhookSecret    *string
	CreatedAt        time.Time
}

type GameMode struct {
	ID           uuid.UUID
	GameID       uuid.UUID
	ModeKey      string
	DisplayName  string
	MinPlayers   int
	MaxPlayers   int
	SeatTemplate json.RawMessage
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type GameModeSeat struct {
	ID          uuid.UUID
	ModeID      uuid.UUID
	SeatKey     string
	Team        *string
	Role        *string
	AffinityKey *string
	QueuePath   *string
	SortOrder   int
}

type ModeQueue struct {
	ID              uuid.UUID
	ModeID          uuid.UUID
	Name            string
	PlayersToStart  int
	Status          string
	IsDefault       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RegisterGameParams struct {
	Slug        string
	Name        string
	Description *string
	PlayURL     string
	APIBaseURL  string
}

type KickedWaiter struct {
	UserID      uuid.UUID
	GameID      uuid.UUID
	ModeQueueID uuid.UUID
	Message     string
}

type ApplyManifestResult struct {
	Game         *Game
	Changed      bool
	Kicked       []KickedWaiter
	WebhookSecret string // set only on register
}

// UserQueueView is the current queue/match state for a user.
type UserQueueView struct {
	GameID      uuid.UUID
	ModeQueueID uuid.UUID
	InQueue     bool
	Waiting     bool
	Matched     bool
	QueuedCount int
	QueuePath   *string
	SessionID   *uuid.UUID
}

// SwitchedFromQueue records a waiting queue the user left to join another.
type SwitchedFromQueue struct {
	ModeQueueID uuid.UUID
	GameID      uuid.UUID
	GameName    string
}

// SwitchedFromPlayerMessage is shown on joinQueue when the player left another game's wait list.
func SwitchedFromPlayerMessage(gameName string) string {
	return "You left the group for " + gameName + " to look for a group here."
}

// QueueJoinResult is returned when a user joins a game queue.
type QueueJoinResult struct {
	GameID         uuid.UUID
	ModeQueueID    uuid.UUID
	Status         string
	SessionID      *uuid.UUID
	QueuedCount    int
	QueuePath      *string
	NotifyUserIDs  []uuid.UUID
	AlreadyInQueue bool
	SwitchedFrom   *SwitchedFromQueue
}

type QueueEntry struct {
	ID          uuid.UUID
	GameID      uuid.UUID
	ModeQueueID *uuid.UUID
	UserID      uuid.UUID
	Status      string
	QueuePath   *string
	JoinedAt    time.Time
}

type Session struct {
	ID          uuid.UUID
	GameID      uuid.UUID
	ModeID      *uuid.UUID
	ModeQueueID *uuid.UUID
	Status      string
	StartedAt   time.Time
	EndedAt     *time.Time
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
