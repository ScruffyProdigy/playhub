package store

import (
	"context"

	"github.com/google/uuid"
)

// DemoDefaultQueueIDStr is the seeded default queue for the local RPS demo game.
const DemoDefaultQueueIDStr = "a3000000-0000-4000-8000-000000000001"

// DemoDefaultQueueID is the parsed form of DemoDefaultQueueIDStr.
var DemoDefaultQueueID = uuid.MustParse(DemoDefaultQueueIDStr)

// DemoRPSGameIDStr is the seeded RPS demo game (migrations 000003 / 000004 / 000009).
// Integration tests must not leave api_base_url pointing at ephemeral httptest ports.
const DemoRPSGameIDStr = "a1000000-0000-4000-8000-000000000001"

// DemoRPSGameID is the parsed form of DemoRPSGameIDStr.
var DemoRPSGameID = uuid.MustParse(DemoRPSGameIDStr)

const (
	DemoGamePlayURL    = "http://localhost:5174"
	DemoGameAPIBaseURL = "http://localhost:3001"
)

// RestoreDemoGameHandoffURLs resets the seeded demo game to local RPS defaults.
func (s *Store) RestoreDemoGameHandoffURLs(ctx context.Context) error {
	return s.SetGameHandoffURLsForTest(ctx, DemoRPSGameID, DemoGamePlayURL, DemoGameAPIBaseURL)
}

// SetGameHandoffURLsForTest updates play/api URLs on a game row (integration tests only).
func (s *Store) SetGameHandoffURLsForTest(ctx context.Context, gameID uuid.UUID, playURL, apiBaseURL string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE games SET play_url = $2, api_base_url = $3 WHERE id = $1
	`, gameID, playURL, apiBaseURL)
	return err
}

// ClearWaitingQueueForGame removes waiting queue rows for a game (integration tests).
func (s *Store) ClearWaitingQueueForGame(ctx context.Context, gameID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM game_queues
		WHERE game_id = $1 AND status = 'waiting'
	`, gameID)
	return err
}

// ClearWaitingModeQueue removes waiting rows for a mode queue (integration tests).
func (s *Store) ClearWaitingModeQueue(ctx context.Context, modeQueueID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM game_queues
		WHERE mode_queue_id = $1 AND status = 'waiting'
	`, modeQueueID)
	return err
}
