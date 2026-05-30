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
	var description, slug, playURL, apiBaseURL, gameMode sql.NullString
	if err := row.Scan(
		&g.ID, &g.Name, &description, &slug, &playURL, &apiBaseURL, &gameMode,
		&g.Status, &g.MinPlayers, &g.MaxPlayers, &g.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if description.Valid {
		g.Description = &description.String
	}
	if slug.Valid {
		g.Slug = &slug.String
	}
	if playURL.Valid {
		g.PlayURL = &playURL.String
	}
	if apiBaseURL.Valid {
		g.APIBaseURL = &apiBaseURL.String
	}
	if gameMode.Valid {
		g.GameMode = &gameMode.String
	}
	return &g, nil
}

const gameColumns = `id, name, description, slug, play_url, api_base_url, game_mode, status, min_players, max_players, created_at`

func (s *Store) ListGames(ctx context.Context, limit, offset int) ([]Game, error) {
	return s.listGames(ctx, limit, offset, "")
}

func (s *Store) ListDemoGames(ctx context.Context, limit, offset int) ([]Game, error) {
	return s.listGames(ctx, limit, offset, "demo")
}

func (s *Store) listGames(ctx context.Context, limit, offset int, category string) ([]Game, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT ` + gameColumns + `
		FROM games
		WHERE status = 'active'`
	args := []any{}

	if category != "" {
		query += ` AND category = $1`
		args = append(args, category)
	}

	limitArg := len(args) + 1
	offsetArg := len(args) + 2
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, limitArg, offsetArg)
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		g, err := scanGame(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, *g)
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
