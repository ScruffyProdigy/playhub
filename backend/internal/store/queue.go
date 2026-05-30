package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func scanQueueEntry(row interface{ Scan(dest ...any) error }) (*QueueEntry, error) {
	var entry QueueEntry
	if err := row.Scan(&entry.ID, &entry.GameID, &entry.UserID, &entry.Status, &entry.JoinedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &entry, nil
}

const queueColumns = `id, game_id, user_id, status, joined_at`

func (s *Store) GetWaitingQueueEntry(ctx context.Context, gameID, userID uuid.UUID) (*QueueEntry, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+queueColumns+`
		FROM game_queues
		WHERE game_id = $1 AND user_id = $2 AND status = 'waiting'
		ORDER BY joined_at DESC
		LIMIT 1
	`, gameID, userID)
	return scanQueueEntry(row)
}

func (s *Store) EnqueueGame(ctx context.Context, gameID, userID uuid.UUID) (*QueueEntry, error) {
	if entry, err := s.GetWaitingQueueEntry(ctx, gameID, userID); err == nil {
		return entry, nil
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO game_queues (game_id, user_id, status)
		VALUES ($1, $2, 'waiting')
		RETURNING `+queueColumns+`
	`, gameID, userID)
	return scanQueueEntry(row)
}

func (s *Store) LeaveQueue(ctx context.Context, gameID, userID uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE game_queues
		SET status = 'cancelled'
		WHERE game_id = $1 AND user_id = $2 AND status = 'waiting'
	`, gameID, userID)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result, ErrNotFound)
}
