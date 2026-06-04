package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// UserActiveQueue is the user's current waiting or matched queue membership, if any.
type UserActiveQueue struct {
	GameID      uuid.UUID
	GameName    string
	ModeQueueID uuid.UUID
	Waiting     bool
	Matched     bool
	QueuedCount int
	SessionID   *uuid.UUID
}

// GetUserActiveQueue returns the user's active queue (waiting preferred over matched).
func (s *Store) GetUserActiveQueue(ctx context.Context, userID uuid.UUID) (*UserActiveQueue, error) {
	if waiting, err := s.getUserWaitingQueueAny(ctx, userID); err != nil {
		return nil, err
	} else if waiting != nil {
		count, err := s.CountWaitingInModeQueue(ctx, waiting.ModeQueueID)
		if err != nil {
			return nil, err
		}
		game, err := s.GetGameByID(ctx, waiting.GameID)
		if err != nil {
			return nil, err
		}
		return &UserActiveQueue{
			GameID:      waiting.GameID,
			GameName:    game.Name,
			ModeQueueID: waiting.ModeQueueID,
			Waiting:     true,
			QueuedCount: count,
		}, nil
	}

	session, modeQueueID, err := s.getUserMatchedModeQueueAny(ctx, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	game, err := s.GetGameByID(ctx, session.GameID)
	if err != nil {
		return nil, err
	}
	return &UserActiveQueue{
		GameID:      session.GameID,
		GameName:    game.Name,
		ModeQueueID: modeQueueID,
		Matched:     true,
		SessionID:   &session.ID,
	}, nil
}

func (s *Store) getUserWaitingQueueAny(ctx context.Context, userID uuid.UUID) (*struct {
	GameID      uuid.UUID
	ModeQueueID uuid.UUID
}, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT game_id, mode_queue_id
		FROM game_queues
		WHERE user_id = $1 AND status = 'waiting' AND mode_queue_id IS NOT NULL
		ORDER BY joined_at DESC
		LIMIT 1
	`, userID)
	var gameID uuid.UUID
	var modeQueueID uuid.UUID
	if err := row.Scan(&gameID, &modeQueueID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &struct {
		GameID      uuid.UUID
		ModeQueueID uuid.UUID
	}{GameID: gameID, ModeQueueID: modeQueueID}, nil
}

func (s *Store) getUserMatchedModeQueueAny(ctx context.Context, userID uuid.UUID) (*Session, uuid.UUID, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT mode_queue_id
		FROM game_queues
		WHERE user_id = $1 AND status = 'matched' AND mode_queue_id IS NOT NULL
		ORDER BY matched_at DESC NULLS LAST, joined_at DESC
		LIMIT 1
	`, userID)
	var modeQueueID uuid.UUID
	if err := row.Scan(&modeQueueID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, uuid.Nil, nil
		}
		return nil, uuid.Nil, err
	}
	session, err := s.GetMatchedSessionForUserAndModeQueue(ctx, modeQueueID, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, uuid.Nil, nil
		}
		return nil, uuid.Nil, err
	}
	return session, modeQueueID, nil
}

func getUserWaitingQueueSwitchInfoTx(ctx context.Context, tx *sql.Tx, userID, exceptModeQueueID uuid.UUID) (*SwitchedFromQueue, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT gq.mode_queue_id, gq.game_id, g.name
		FROM game_queues gq
		INNER JOIN games g ON g.id = gq.game_id
		WHERE gq.user_id = $1 AND gq.status = 'waiting' AND gq.mode_queue_id IS NOT NULL
		  AND gq.mode_queue_id <> $2
		ORDER BY gq.joined_at DESC
		LIMIT 1
	`, userID, exceptModeQueueID)
	var info SwitchedFromQueue
	if err := row.Scan(&info.ModeQueueID, &info.GameID, &info.GameName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func getUserMatchedModeQueueIDTx(ctx context.Context, tx *sql.Tx, userID, exceptModeQueueID uuid.UUID) (*uuid.UUID, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT mode_queue_id
		FROM game_queues
		WHERE user_id = $1 AND status = 'matched' AND mode_queue_id IS NOT NULL
		  AND mode_queue_id <> $2
		ORDER BY matched_at DESC NULLS LAST, joined_at DESC
		LIMIT 1
	`, userID, exceptModeQueueID)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &id, nil
}
