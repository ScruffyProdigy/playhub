package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const (
	QueueStatusWaiting   = "waiting"
	QueueStatusMatched   = "matched"
	QueueStatusCancelled = "cancelled"
)

func getGameForMatchmaking(ctx context.Context, q sqlQueryRowContext, gameID uuid.UUID) (*Game, error) {
	row := q.QueryRowContext(ctx, `
		SELECT `+gameColumns+`
		FROM games
		WHERE id = $1 AND status = 'active'
	`, gameID)
	return scanGame(row)
}

func pickDistinctWaitingEntries(entries []QueueEntry, n int) []QueueEntry {
	if n <= 0 || len(entries) == 0 {
		return nil
	}
	seen := make(map[uuid.UUID]struct{}, n)
	out := make([]QueueEntry, 0, n)
	for _, entry := range entries {
		if _, ok := seen[entry.UserID]; ok {
			continue
		}
		seen[entry.UserID] = struct{}{}
		out = append(out, entry)
		if len(out) == n {
			break
		}
	}
	return out
}

func markQueueEntryMatchedTx(ctx context.Context, tx *sql.Tx, entryID uuid.UUID) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE game_queues
		SET status = 'matched', matched_at = NOW()
		WHERE id = $1 AND status = 'waiting'
	`, entryID)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result, ErrNotFound)
}

func addSessionParticipantTx(ctx context.Context, tx *sql.Tx, sessionID, userID uuid.UUID, role string) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO game_session_participants (session_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_id, user_id) DO NOTHING
	`, sessionID, userID, role)
	return err
}
