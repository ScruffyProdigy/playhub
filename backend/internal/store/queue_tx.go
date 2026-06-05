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

func addSessionParticipantTx(ctx context.Context, tx *sql.Tx, sessionID, userID uuid.UUID, role string, returnCtx ReturnContext) error {
	raw, err := encodeReturnContext(returnCtx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO game_session_participants (session_id, user_id, role, return_context)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (session_id, user_id) DO UPDATE SET
			role = EXCLUDED.role,
			return_context = EXCLUDED.return_context,
			left_at = NULL,
			finished_at = NULL
	`, sessionID, userID, role, raw)
	return err
}
