package store

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID               uuid.UUID
	Name             string
	Description      *string
	IconURL          string
	HeroURL          string
	CatalogHeroURL   *string
	ShortDescription *string
	HowToPlay        *string
	TutorialURL      *string
	Screenshots      []string
	Tags             []string
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
	ID             uuid.UUID
	ModeID         uuid.UUID
	Name           string
	PlayersToStart int
	Status         string
	IsDefault      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RegisterGameParams struct {
	Slug        string
	Name        string
	Description *string
	IconURL     string
	HeroURL     string
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
	Game          *Game
	Changed       bool
	Kicked        []KickedWaiter
	WebhookSecret string // set only on register
}
