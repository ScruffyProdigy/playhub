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

func (s *Store) expireStaleMatchedModeQueue(ctx context.Context, modeQueueID, userID uuid.UUID) error {
	cutoff := time.Now().Add(-staleMatchMaxAge())
	_, err := s.db.ExecContext(ctx, `
		UPDATE game_queues
		SET status = 'cancelled'
		WHERE mode_queue_id = $1
		  AND user_id = $2
		  AND status = 'matched'
		  AND (matched_at IS NULL OR matched_at < $3)
	`, modeQueueID, userID, cutoff)
	return err
}

func (s *Store) cancelUserMatchedModeQueue(ctx context.Context, modeQueueID, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE game_queues
		SET status = 'cancelled'
		WHERE mode_queue_id = $1 AND user_id = $2 AND status = 'matched'
	`, modeQueueID, userID)
	return err
}
