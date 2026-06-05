package store

import (
	"context"

	"github.com/google/uuid"
)

// DemoDefaultQueueIDStr is the seeded default mode queue for demo game 001.
const DemoDefaultQueueIDStr = "a3000000-0000-4000-8000-000000000001"

// DemoDefaultQueueID is the parsed form of DemoDefaultQueueIDStr.
var DemoDefaultQueueID = uuid.MustParse(DemoDefaultQueueIDStr)

// DemoPrimaryGameIDStr is seeded catalog game 001 (migrations 000003 / 000009).
// Production uses this row for Word Hunt; integration tests repoint handoff URLs to localhost.
const DemoPrimaryGameIDStr = "a1000000-0000-4000-8000-000000000001"

// DemoPrimaryGameID is the parsed form of DemoPrimaryGameIDStr.
var DemoPrimaryGameID = uuid.MustParse(DemoPrimaryGameIDStr)

const (
	DemoGamePlayURL    = "http://localhost:5174"
	DemoGameAPIBaseURL = "http://localhost:3001"
)

// RestorePrimaryGameHandoffURLs resets game 001 to local dev handoff URLs after tests.
func (s *Store) RestorePrimaryGameHandoffURLs(ctx context.Context) error {
	return s.SetGameHandoffURLsForTest(ctx, DemoPrimaryGameID, DemoGamePlayURL, DemoGameAPIBaseURL)
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
