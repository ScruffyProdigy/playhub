package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func scanSessionRow(row interface{ Scan(dest ...any) error }) (*Session, error) {
	var session Session
	var endedAt sql.NullTime
	var modeID, modeQueueID sql.NullString
	if err := row.Scan(
		&session.ID, &session.GameID, &modeID, &modeQueueID,
		&session.Status, &session.StartedAt, &endedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if modeID.Valid {
		id, err := uuid.Parse(modeID.String)
		if err != nil {
			return nil, err
		}
		session.ModeID = &id
	}
	if modeQueueID.Valid {
		id, err := uuid.Parse(modeQueueID.String)
		if err != nil {
			return nil, err
		}
		session.ModeQueueID = &id
	}
	if endedAt.Valid {
		t := endedAt.Time
		session.EndedAt = &t
	}
	return &session, nil
}

const sessionColumns = `id, game_id, mode_id, mode_queue_id, status, started_at, ended_at`

func (s *Store) GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+sessionColumns+`
		FROM game_sessions
		WHERE id = $1
	`, id)
	return scanSessionRow(row)
}

func (s *Store) ListActiveSessionsByGame(ctx context.Context, gameID uuid.UUID, limit int) ([]Session, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT `+sessionColumns+`
		FROM game_sessions
		WHERE game_id = $1 AND status = 'active'
		ORDER BY started_at DESC
		LIMIT $2
	`, gameID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		session, err := scanSessionRow(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, *session)
	}
	return sessions, rows.Err()
}

func (s *Store) CreateSession(ctx context.Context, gameID uuid.UUID) (*Session, error) {
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO game_sessions (game_id, status)
		VALUES ($1, 'active')
		RETURNING `+sessionColumns+`
	`, gameID)
	return scanSessionRow(row)
}

func (s *Store) AddSessionParticipant(ctx context.Context, sessionID, userID uuid.UUID, role string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO game_session_participants (session_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_id, user_id) DO NOTHING
	`, sessionID, userID, role)
	return err
}

// GetMatchedSessionForUserAndModeQueue returns the active session for a matched mode-queue row.
// When multiple active sessions exist (e.g. reportMatchResult failed on an earlier match),
// the newest session for this mode queue wins.
func (s *Store) GetMatchedSessionForUserAndModeQueue(ctx context.Context, modeQueueID, userID uuid.UUID) (*Session, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT gs.id, gs.game_id, gs.mode_id, gs.mode_queue_id, gs.status, gs.started_at, gs.ended_at
		FROM game_queues q
		INNER JOIN game_session_participants p
			ON p.user_id = q.user_id AND p.left_at IS NULL
		INNER JOIN game_sessions gs
			ON gs.id = p.session_id
			AND gs.game_id = q.game_id
			AND gs.mode_queue_id = q.mode_queue_id
			AND gs.status = 'active'
		WHERE q.mode_queue_id = $1
		  AND q.user_id = $2
		  AND q.status = 'matched'
		ORDER BY gs.started_at DESC, q.matched_at DESC NULLS LAST, q.joined_at DESC
		LIMIT 1
	`, modeQueueID, userID)
	return scanSessionRow(row)
}

func (s *Store) ListSessionParticipants(ctx context.Context, sessionID uuid.UUID) ([]User, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+strings.ReplaceAll(userColumns, "id,", "u.id,")+`
		FROM game_session_participants p
		JOIN users u ON u.id = p.user_id
		WHERE p.session_id = $1 AND p.left_at IS NULL
		ORDER BY p.joined_at ASC
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}
