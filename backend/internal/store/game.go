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
	var description, slug, playURL, apiBaseURL sql.NullString
	var manifestHash, manifestETag, gameVersion, webhookSecret sql.NullString
	var manifestSyncedAt sql.NullTime
	if err := row.Scan(
		&g.ID, &g.Name, &description, &slug, &playURL, &apiBaseURL,
		&g.Status,
		&manifestHash, &manifestETag, &manifestSyncedAt, &gameVersion, &webhookSecret,
		&g.CreatedAt,
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
	if manifestHash.Valid {
		g.ManifestHash = &manifestHash.String
	}
	if manifestETag.Valid {
		g.ManifestETag = &manifestETag.String
	}
	if manifestSyncedAt.Valid {
		g.ManifestSyncedAt = &manifestSyncedAt.Time
	}
	if gameVersion.Valid {
		g.GameVersion = &gameVersion.String
	}
	if webhookSecret.Valid {
		g.WebhookSecret = &webhookSecret.String
	}
	return &g, nil
}

const gameColumns = `id, name, description, slug, play_url, api_base_url, status,
	manifest_hash, manifest_etag, manifest_synced_at, game_version, webhook_secret, created_at`

func (s *Store) ListGames(ctx context.Context, limit, offset int) ([]Game, error) {
	return s.ListCatalogGames(ctx, limit, offset)
}

// ListCatalogGames returns active games that have at least one active mode queue.
func (s *Store) ListCatalogGames(ctx context.Context, limit, offset int) ([]Game, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT `+gameColumns+`
		FROM games g
		WHERE g.status = 'active'
		  AND EXISTS (
		    SELECT 1
		    FROM game_modes gm
		    INNER JOIN mode_queues mq ON mq.mode_id = gm.id
		    WHERE gm.game_id = g.id
		      AND gm.status = 'active'
		      AND mq.status = 'active'
		  )
		ORDER BY g.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
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

func (s *Store) InsertTestGame(ctx context.Context, name string) (*Game, error) {
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
