package store

import (
	"context"

	"github.com/google/uuid"
)

// RollbackMatchedSession cancels a failed provision: removes the session and returns
// non-banned matched players to waiting; banned players are removed from the queue.
func (s *Store) RollbackMatchedSession(ctx context.Context, sessionID uuid.UUID, bannedUserIDs []uuid.UUID) error {
	banned := make(map[uuid.UUID]struct{}, len(bannedUserIDs))
	for _, id := range bannedUserIDs {
		banned[id] = struct{}{}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var gameID uuid.UUID
	if err := tx.QueryRowContext(ctx, `SELECT game_id FROM game_sessions WHERE id = $1`, sessionID).Scan(&gameID); err != nil {
		return err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT user_id FROM game_session_participants WHERE session_id = $1
	`, sessionID)
	if err != nil {
		return err
	}
	var userIDs []uuid.UUID
	for rows.Next() {
		var uid uuid.UUID
		if err := rows.Scan(&uid); err != nil {
			rows.Close()
			return err
		}
		userIDs = append(userIDs, uid)
	}
	rows.Close()

	if _, err := tx.ExecContext(ctx, `DELETE FROM game_session_participants WHERE session_id = $1`, sessionID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM game_sessions WHERE id = $1`, sessionID); err != nil {
		return err
	}

	for _, userID := range userIDs {
		if _, ok := banned[userID]; ok {
			if _, err := tx.ExecContext(ctx, `
				UPDATE game_queues
				SET status = 'cancelled'
				WHERE game_id = $1 AND user_id = $2 AND status = 'matched'
			`, gameID, userID); err != nil {
				return err
			}
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE game_queues
			SET status = 'waiting', matched_at = NULL
			WHERE game_id = $1 AND user_id = $2 AND status = 'matched'
		`, gameID, userID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
