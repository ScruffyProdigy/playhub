package store

import (
	"context"

	"github.com/google/uuid"
)

// DemoQuickMatchGameIDStr is the seeded RPS demo game (migrations 000003 / 000004).
// Integration tests must not leave api_base_url pointing at ephemeral httptest ports.
const DemoQuickMatchGameIDStr = "a1000000-0000-4000-8000-000000000001"

// DemoQuickMatchGameID is the parsed form of DemoQuickMatchGameIDStr.
var DemoQuickMatchGameID = uuid.MustParse(DemoQuickMatchGameIDStr)

const (
	DemoGamePlayURL    = "http://localhost:5174"
	DemoGameAPIBaseURL = "http://localhost:3001"
)

// RestoreDemoGameHandoffURLs resets the seeded demo game to local RPS defaults.
func (s *Store) RestoreDemoGameHandoffURLs(ctx context.Context) error {
	return s.SetGameHandoffURLsForTest(ctx, DemoQuickMatchGameID, DemoGamePlayURL, DemoGameAPIBaseURL)
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
