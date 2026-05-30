package store

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const defaultStaleMatchMinutes = 30

// staleMatchMaxAge returns how long a matched queue row may persist before auto-expiry.
func staleMatchMaxAge() time.Duration {
	if v := os.Getenv("LOBBY_STALE_MATCH_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Minute
		}
	}
	return defaultStaleMatchMinutes * time.Minute
}

// ExpireStaleMatchedQueue cancels old matched rows so login does not show a stale Launch.
func (s *Store) ExpireStaleMatchedQueue(ctx context.Context, gameID, userID uuid.UUID) error {
	cutoff := time.Now().Add(-staleMatchMaxAge())
	_, err := s.db.ExecContext(ctx, `
		UPDATE game_queues
		SET status = 'cancelled'
		WHERE game_id = $1
		  AND user_id = $2
		  AND status = 'matched'
		  AND (matched_at IS NULL OR matched_at < $3)
	`, gameID, userID, cutoff)
	return err
}

// CancelUserMatchedQueue clears an active matched queue row (user abandons the lobby match).
func (s *Store) CancelUserMatchedQueue(ctx context.Context, gameID, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE game_queues
		SET status = 'cancelled'
		WHERE game_id = $1 AND user_id = $2 AND status = 'matched'
	`, gameID, userID)
	return err
}
