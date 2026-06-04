package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ParticipantIsActive reports whether the user is a seated, not-left participant.
func (s *Store) ParticipantIsActive(ctx context.Context, sessionID, userID uuid.UUID) error {
	var one int
	err := s.db.QueryRowContext(ctx, `
		SELECT 1
		FROM game_session_participants
		WHERE session_id = $1 AND user_id = $2 AND left_at IS NULL
	`, sessionID, userID).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// GetParticipantReturnContext loads per-player return routing for a session.
func (s *Store) GetParticipantReturnContext(ctx context.Context, sessionID, userID uuid.UUID) (ReturnContext, error) {
	var raw []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT return_context
		FROM game_session_participants
		WHERE session_id = $1 AND user_id = $2 AND left_at IS NULL
	`, sessionID, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ReturnContext{}, ErrNotFound
		}
		return ReturnContext{}, err
	}
	return decodeReturnContext(raw)
}

// MarkParticipantFinished records that a player is done with the match for themselves.
func (s *Store) MarkParticipantFinished(ctx context.Context, sessionID, userID uuid.UUID, finishedAt time.Time) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE game_session_participants
		SET finished_at = $3
		WHERE session_id = $1 AND user_id = $2 AND left_at IS NULL AND finished_at IS NULL
	`, sessionID, userID, finishedAt)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result, ErrNotFound)
}

// AcknowledgePlayerReturn clears the user's matched queue row when they land on the
// return hub after a game. Complements reportMatchResult from the game server.
func (s *Store) AcknowledgePlayerReturn(ctx context.Context, sessionID, userID uuid.UUID, returnedAt time.Time) error {
	session, err := s.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if err := s.MarkParticipantFinished(ctx, sessionID, userID, returnedAt); err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	if session.ModeQueueID != nil {
		if err := s.ReleaseUserMatchedQueue(ctx, *session.ModeQueueID, userID); err != nil {
			return err
		}
	}
	return nil
}

// ReleaseUserMatchedQueue clears the user's matched row for a mode queue so they can queue again.
func (s *Store) ReleaseUserMatchedQueue(ctx context.Context, modeQueueID, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE game_queues
		SET status = 'cancelled'
		WHERE mode_queue_id = $1 AND user_id = $2 AND status = 'matched'
	`, modeQueueID, userID)
	return err
}

// completePriorActiveSessionsForModeQueueTx ends older active sessions in the same mode
// queue so re-queues never hand off a stale match id to the game server.
func completePriorActiveSessionsForModeQueueTx(ctx context.Context, tx *sql.Tx, modeQueueID, exceptSessionID uuid.UUID, endedAt time.Time) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE game_sessions
		SET status = 'completed', ended_at = $3
		WHERE mode_queue_id = $1 AND status = 'active' AND id <> $2
	`, modeQueueID, exceptSessionID, endedAt)
	return err
}

// CompleteSession marks a session ended and releases matched queue rows for all seated players.
func (s *Store) CompleteSession(ctx context.Context, sessionID uuid.UUID, endedAt time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var modeQueueID sql.NullString
	err = tx.QueryRowContext(ctx, `
		UPDATE game_sessions
		SET status = 'completed', ended_at = $2
		WHERE id = $1 AND status = 'active'
		RETURNING mode_queue_id
	`, sessionID, endedAt).Scan(&modeQueueID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if modeQueueID.Valid {
		mqID, parseErr := uuid.Parse(modeQueueID.String)
		if parseErr != nil {
			return parseErr
		}
		rows, qErr := tx.QueryContext(ctx, `
			SELECT user_id FROM game_session_participants
			WHERE session_id = $1 AND left_at IS NULL
		`, sessionID)
		if qErr != nil {
			return qErr
		}
		var participantIDs []uuid.UUID
		for rows.Next() {
			var uid uuid.UUID
			if scanErr := rows.Scan(&uid); scanErr != nil {
				_ = rows.Close()
				return scanErr
			}
			participantIDs = append(participantIDs, uid)
		}
		if err := rows.Close(); err != nil {
			return err
		}
		if err := rows.Err(); err != nil {
			return err
		}
		for _, uid := range participantIDs {
			if _, execErr := tx.ExecContext(ctx, `
				UPDATE game_queues
				SET status = 'cancelled'
				WHERE mode_queue_id = $1 AND user_id = $2 AND status = 'matched'
			`, mqID, uid); execErr != nil {
				return execErr
			}
		}
	}

	return tx.Commit()
}

// CountActiveParticipants returns seated players who have not finished individually.
func (s *Store) CountActiveParticipants(ctx context.Context, sessionID uuid.UUID) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM game_session_participants
		WHERE session_id = $1 AND left_at IS NULL AND finished_at IS NULL
	`, sessionID).Scan(&n)
	return n, err
}
