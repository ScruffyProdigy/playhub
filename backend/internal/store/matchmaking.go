package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

const (
	QueueStatusWaiting   = "waiting"
	QueueStatusMatched   = "matched"
	QueueStatusCancelled = "cancelled"
)

// JoinGameQueue enqueues the user and starts a session when enough distinct players are waiting.
func (s *Store) JoinGameQueue(ctx context.Context, gameID, userID uuid.UUID) (*QueueJoinResult, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	game, err := getGameForMatchmaking(ctx, tx, gameID)
	if err != nil {
		return nil, err
	}
	minPlayers := game.MinPlayers
	if minPlayers < 1 {
		minPlayers = 1
	}

	alreadyWaiting, err := enqueueGameTx(ctx, tx, gameID, userID)
	if err != nil {
		return nil, err
	}

	waiting, err := listWaitingQueueEntriesTx(ctx, tx, gameID)
	if err != nil {
		return nil, err
	}

	matched := pickDistinctWaitingEntries(waiting, minPlayers)
	if len(matched) < minPlayers {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &QueueJoinResult{
			Status:         QueueStatusWaiting,
			QueuedCount:    len(waiting),
			NotifyUserIDs:  []uuid.UUID{userID},
			AlreadyInQueue: alreadyWaiting,
		}, nil
	}

	session, err := createSessionTx(ctx, tx, gameID)
	if err != nil {
		return nil, err
	}

	seatKeys := defaultSeatKeysForCount(len(matched))
	notifyIDs := make([]uuid.UUID, 0, len(matched))
	seenNotify := make(map[uuid.UUID]struct{}, len(matched))
	for i, entry := range matched {
		if err := markQueueEntryMatchedTx(ctx, tx, entry.ID); err != nil {
			return nil, err
		}
		seatKey := seatKeys[i]
		if err := addSessionParticipantTx(ctx, tx, session.ID, entry.UserID, seatKey); err != nil {
			return nil, err
		}
		if _, ok := seenNotify[entry.UserID]; !ok {
			seenNotify[entry.UserID] = struct{}{}
			notifyIDs = append(notifyIDs, entry.UserID)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &QueueJoinResult{
		Status:        QueueStatusMatched,
		SessionID:     &session.ID,
		QueuedCount:   0,
		NotifyUserIDs: notifyIDs,
	}, nil
}

// LeaveGameQueue removes the user from waiting or matched queue state for a game.
func (s *Store) LeaveGameQueue(ctx context.Context, gameID, userID uuid.UUID) (int, error) {
	_ = s.LeaveQueue(ctx, gameID, userID) // ok if not waiting
	_ = s.CancelUserMatchedQueue(ctx, gameID, userID)
	return s.countWaitingQueueEntries(ctx, gameID)
}

func (s *Store) countWaitingQueueEntries(ctx context.Context, gameID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM game_queues
		WHERE game_id = $1 AND status = 'waiting'
	`, gameID).Scan(&count)
	return count, err
}

func getGameForMatchmaking(ctx context.Context, q sqlQueryRow, gameID uuid.UUID) (*Game, error) {
	row := q.QueryRowContext(ctx, `
		SELECT `+gameColumns+`
		FROM games
		WHERE id = $1 AND status = 'active'
	`, gameID)
	return scanGame(row)
}

type sqlQueryRow interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// enqueueGameTx returns (alreadyWaiting, error). alreadyWaiting is true when the user
// was already in the waiting queue (idempotent re-join).
func enqueueGameTx(ctx context.Context, tx *sql.Tx, gameID, userID uuid.UUID) (bool, error) {
	if _, err := getMatchedQueueEntryTx(ctx, tx, gameID, userID); err == nil {
		return false, ErrAlreadyMatched
	} else if !errors.Is(err, ErrNotFound) {
		return false, err
	}

	row := tx.QueryRowContext(ctx, `
		SELECT `+queueColumns+`
		FROM game_queues
		WHERE game_id = $1 AND user_id = $2 AND status = 'waiting'
		ORDER BY joined_at DESC
		LIMIT 1
	`, gameID, userID)
	if _, err := scanQueueEntry(row); err == nil {
		return true, nil
	} else if !errors.Is(err, ErrNotFound) {
		return false, err
	}

	row = tx.QueryRowContext(ctx, `
		INSERT INTO game_queues (game_id, user_id, status)
		VALUES ($1, $2, 'waiting')
		RETURNING `+queueColumns+`
	`, gameID, userID)
	if _, err := scanQueueEntry(row); err != nil {
		if isUniqueViolation(err) {
			row = tx.QueryRowContext(ctx, `
				SELECT `+queueColumns+`
				FROM game_queues
				WHERE game_id = $1 AND user_id = $2 AND status = 'waiting'
				ORDER BY joined_at DESC
				LIMIT 1
			`, gameID, userID)
			if _, err := scanQueueEntry(row); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func getMatchedQueueEntryTx(ctx context.Context, tx *sql.Tx, gameID, userID uuid.UUID) (*QueueEntry, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT `+queueColumns+`
		FROM game_queues
		WHERE game_id = $1 AND user_id = $2 AND status = 'matched'
		ORDER BY matched_at DESC NULLS LAST, joined_at DESC
		LIMIT 1
	`, gameID, userID)
	return scanQueueEntry(row)
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

func listWaitingQueueEntriesTx(ctx context.Context, tx *sql.Tx, gameID uuid.UUID) ([]QueueEntry, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT `+queueColumns+`
		FROM game_queues
		WHERE game_id = $1 AND status = 'waiting'
		ORDER BY joined_at ASC
		FOR UPDATE
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []QueueEntry
	for rows.Next() {
		var entry QueueEntry
		if err := rows.Scan(&entry.ID, &entry.GameID, &entry.UserID, &entry.Status, &entry.JoinedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func createSessionTx(ctx context.Context, tx *sql.Tx, gameID uuid.UUID) (*Session, error) {
	row := tx.QueryRowContext(ctx, `
		INSERT INTO game_sessions (game_id, status)
		VALUES ($1, 'active')
		RETURNING `+sessionColumns+`
	`, gameID)
	return scanSession(row)
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
