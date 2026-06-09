package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func listGameModeSeats(ctx context.Context, q sqlQueryRowContext, modeID uuid.UUID) ([]GameModeSeat, error) {
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
