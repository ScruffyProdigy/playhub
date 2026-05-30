package store

import (
	"context"

	"github.com/google/uuid"
)

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
