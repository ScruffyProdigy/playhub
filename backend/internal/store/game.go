package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func scanGame(row interface{ Scan(dest ...any) error }) (*Game, error) {
	var g Game
	var description sql.NullString
	if err := row.Scan(&g.ID, &g.Name, &description, &g.Status, &g.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if description.Valid {
		g.Description = &description.String
	}
	return &g, nil
}

const gameColumns = `id, name, description, status, created_at`

func (s *Store) ListGames(ctx context.Context, limit, offset int) ([]Game, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT `+gameColumns+`
		FROM games
		WHERE status = 'active'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var g Game
		var description sql.NullString
		if err := rows.Scan(&g.ID, &g.Name, &description, &g.Status, &g.CreatedAt); err != nil {
			return nil, err
		}
		if description.Valid {
			g.Description = &description.String
		}
		games = append(games, g)
	}
	return games, rows.Err()
}

func (s *Store) GetGameByID(ctx context.Context, id uuid.UUID) (*Game, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+gameColumns+`
		FROM games
		WHERE id = $1
	`, id)
	return scanGame(row)
}

func (s *Store) CreateGame(ctx context.Context, name string) (*Game, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("store: game name is required")
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO games (name)
		VALUES ($1)
		RETURNING `+gameColumns+`
	`, name)
	return scanGame(row)
}
