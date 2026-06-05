package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

const queueColumns = `id, game_id, user_id, status, joined_at, mode_queue_id, queue_path`

func scanQueueEntry(row interface{ Scan(dest ...any) error }) (*QueueEntry, error) {
	var entry QueueEntry
	var modeQueueID, queuePath sql.NullString
	if err := row.Scan(&entry.ID, &entry.GameID, &entry.UserID, &entry.Status, &entry.JoinedAt, &modeQueueID, &queuePath); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if modeQueueID.Valid {
		id, err := uuid.Parse(modeQueueID.String)
		if err != nil {
			return nil, err
		}
		entry.ModeQueueID = &id
	}
	if queuePath.Valid {
		entry.QueuePath = &queuePath.String
	}
	return &entry, nil
}

func scanQueueEntryRow(rows *sql.Rows) (*QueueEntry, error) {
	var entry QueueEntry
	var modeQueueID, queuePath sql.NullString
	if err := rows.Scan(&entry.ID, &entry.GameID, &entry.UserID, &entry.Status, &entry.JoinedAt, &modeQueueID, &queuePath); err != nil {
		return nil, err
	}
	if modeQueueID.Valid {
		id, err := uuid.Parse(modeQueueID.String)
		if err != nil {
			return nil, err
		}
		entry.ModeQueueID = &id
	}
	if queuePath.Valid {
		entry.QueuePath = &queuePath.String
	}
	return &entry, nil
}

func (s *Store) GetWaitingModeQueueEntry(ctx context.Context, modeQueueID, userID uuid.UUID) (*QueueEntry, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+queueColumns+`
		FROM game_queues
		WHERE mode_queue_id = $1 AND user_id = $2 AND status = 'waiting'
		ORDER BY joined_at DESC
		LIMIT 1
	`, modeQueueID, userID)
	return scanQueueEntry(row)
}

func listGameModeSeatsQuery(ctx context.Context, q sqlQueryRowContext, modeID uuid.UUID) ([]GameModeSeat, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT id, mode_id, seat_key, team, role, affinity_key, queue_path, sort_order
		FROM game_mode_seats
		WHERE mode_id = $1
		ORDER BY sort_order ASC, seat_key ASC
	`, modeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []GameModeSeat
	for rows.Next() {
		var seat GameModeSeat
		var team, role, affinity, queuePath sql.NullString
		if err := rows.Scan(&seat.ID, &seat.ModeID, &seat.SeatKey, &team, &role, &affinity, &queuePath, &seat.SortOrder); err != nil {
			return nil, err
		}
		if team.Valid {
			seat.Team = &team.String
		}
		if role.Valid {
			seat.Role = &role.String
		}
		if affinity.Valid {
			seat.AffinityKey = &affinity.String
		}
		if queuePath.Valid {
			seat.QueuePath = &queuePath.String
		}
		seats = append(seats, seat)
	}
	return seats, rows.Err()
}
