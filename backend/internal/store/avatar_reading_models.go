package store

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	ReadingStatusGeneratingQuestions = "generating_questions"
	ReadingStatusAwaitingAnswers     = "awaiting_answers"
	ReadingStatusProcessing          = "processing"
	ReadingStatusReady               = "ready"
	ReadingStatusCompleted           = "completed"
	ReadingStatusFailed              = "failed"
)

const SourceSpiritAnimal = "spirit_animal"

// AvatarReading is a spirit-animal flow session and its persisted LLM artifacts.
type AvatarReading struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Status              string
	Draw                []int32
	QuestionsJSON       json.RawMessage
	UserAnswers         []string
	PersonalityJSON     json.RawMessage
	TotemsJSON          json.RawMessage
	RankingJSON         json.RawMessage
	SelectedTotemName   *string
	ArtDirectionVersion int
	ErrorMessage        *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	PhaseStartedAt      *time.Time
	CompletedAt         *time.Time
}

// AvatarRender is a generated mascot image for a reading.
type AvatarRender struct {
	ID                  uuid.UUID
	ReadingID           uuid.UUID
	TotemName           string
	ArtDirectionVersion int
	ImageURL            string
	ImagePrompt         string
	CreatedAt           time.Time
}
