package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
)

const (
	ModeStatusActive   = "active"
	ModeStatusInactive = "inactive"
	ModeQueueActive    = "active"
	ModeQueueInactive  = "inactive"
	defaultQueueName   = "Default"
)

// RegisterGame inserts a catalog game and applies the fetched manifest.
func (s *Store) RegisterGame(ctx context.Context, params RegisterGameParams, manifest *gameclient.Manifest) (*ApplyManifestResult, error) {
	if manifest == nil {
		return nil, fmt.Errorf("store: manifest is required")
	}

	slug := strings.TrimSpace(params.Slug)
	if slug == "" {
		return nil, fmt.Errorf("store: slug is required")
	}
	playURL := strings.TrimRight(strings.TrimSpace(params.PlayURL), "/")
	apiBaseURL := strings.TrimRight(strings.TrimSpace(params.APIBaseURL), "/")
	if playURL == "" || apiBaseURL == "" {
		return nil, fmt.Errorf("store: playUrl and apiBaseUrl are required")
	}

	name := strings.TrimSpace(params.Name)
	if name == "" {
		name = strings.TrimSpace(manifest.Status.Game)
	}
	if name == "" {
		name = slug
	}

	secret, err := generateWebhookSecret()
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO games (name, description, slug, play_url, api_base_url, category, status, webhook_secret)
		VALUES ($1, $2, $3, $4, $5, 'catalog', 'active', $6)
		RETURNING `+gameColumns+`
	`, name, params.Description, slug, playURL, apiBaseURL, secret)
	game, err := scanGame(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("store: game slug %q already exists", slug)
		}
		return nil, err
	}

	kicked, changed, err := applyManifestTx(ctx, tx, game, manifest)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	game, err = s.GetGameByID(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	return &ApplyManifestResult{
		Game:          game,
		Changed:       changed,
		Kicked:        kicked,
		WebhookSecret: secret,
	}, nil
}

// ApplyGameManifest refreshes cached modes/seats/queues from a newly fetched manifest.
func (s *Store) ApplyGameManifest(ctx context.Context, gameID uuid.UUID, manifest *gameclient.Manifest) (*ApplyManifestResult, error) {
	if manifest == nil {
		return nil, fmt.Errorf("store: manifest is required")
	}

	game, err := s.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.APIBaseURL == nil || strings.TrimSpace(*game.APIBaseURL) == "" {
		return nil, fmt.Errorf("store: game has no api_base_url configured")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	kicked, changed, err := applyManifestTx(ctx, tx, game, manifest)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	game, err = s.GetGameByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return &ApplyManifestResult{
		Game:    game,
		Changed: changed,
		Kicked:  kicked,
	}, nil
}

func applyManifestTx(ctx context.Context, tx *sql.Tx, game *Game, manifest *gameclient.Manifest) ([]KickedWaiter, bool, error) {
	if game.ManifestHash != nil && *game.ManifestHash == manifest.SHA256Hash {
		return nil, false, nil
	}

	now := time.Now()
	version := strings.TrimSpace(manifest.Status.Version)
	_, err := tx.ExecContext(ctx, `
		UPDATE games
		SET manifest_hash = $2,
		    manifest_etag = $3,
		    manifest_synced_at = $4,
		    game_version = NULLIF($5, ''),
		    updated_at = NOW()
		WHERE id = $1
	`, game.ID, manifest.SHA256Hash, nullIfEmpty(manifest.ETag), now, version)
	if err != nil {
		return nil, false, err
	}

	existingModes, err := listGameModesTx(ctx, tx, game.ID)
	if err != nil {
		return nil, false, err
	}

	existingByKey := make(map[string]*GameMode, len(existingModes))
	for i := range existingModes {
		existingByKey[existingModes[i].ModeKey] = &existingModes[i]
	}

	seenKeys := make(map[string]struct{}, len(manifest.Modes))
	var kicked []KickedWaiter

	for _, modeDef := range manifest.Modes {
		modeKey := strings.TrimSpace(modeDef.Key)
		seenKeys[modeKey] = struct{}{}

		displayName := strings.TrimSpace(modeDef.DisplayName)
		if displayName == "" {
			displayName = modeKey
		}

		minPlayers := modeDef.MinPlayers
		if minPlayers <= 0 {
			minPlayers = len(modeDef.Seats)
		}
		maxPlayers := modeDef.MaxPlayers
		if maxPlayers <= 0 {
			maxPlayers = len(modeDef.Seats)
		}

		mode, err := upsertGameModeTx(ctx, tx, game.ID, modeKey, displayName, minPlayers, maxPlayers)
		if err != nil {
			return nil, false, err
		}

		newSeatKeys := seatKeysFromManifest(modeDef.Seats)
		oldSeatKeys, err := listModeSeatKeysTx(ctx, tx, mode.ID)
		if err != nil {
			return nil, false, err
		}

		if !equalStringSets(oldSeatKeys, newSeatKeys) && len(oldSeatKeys) > 0 {
			k, err := deactivateModeQueuesTx(ctx, tx, mode.ID, "This queue closed because the game mode seats changed.")
			if err != nil {
				return nil, false, err
			}
			kicked = append(kicked, k...)
		}

		if err := replaceModeSeatsTx(ctx, tx, mode.ID, modeDef.Seats); err != nil {
			return nil, false, err
		}

		playersToStart := len(modeDef.Seats)
		if _, err := ensureDefaultModeQueueTx(ctx, tx, mode.ID, playersToStart); err != nil {
			return nil, false, err
		}
	}

	for _, mode := range existingModes {
		if _, ok := seenKeys[mode.ModeKey]; ok {
			continue
		}
		if err := setGameModeStatusTx(ctx, tx, mode.ID, ModeStatusInactive); err != nil {
			return nil, false, err
		}
		k, err := deactivateModeQueuesTx(ctx, tx, mode.ID, "This queue closed because the game mode was removed.")
		if err != nil {
			return nil, false, err
		}
		kicked = append(kicked, k...)
	}

	return kicked, true, nil
}

func (s *Store) ListGameModesByGameID(ctx context.Context, gameID uuid.UUID) ([]GameMode, error) {
	return listGameModes(ctx, s.db, gameID)
}

func (s *Store) ListGameModeSeats(ctx context.Context, modeID uuid.UUID) ([]GameModeSeat, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, mode_id, seat_key, team, role, sort_order
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
		var team, role sql.NullString
		if err := rows.Scan(&seat.ID, &seat.ModeID, &seat.SeatKey, &team, &role, &seat.SortOrder); err != nil {
			return nil, err
		}
		if team.Valid {
			seat.Team = &team.String
		}
		if role.Valid {
			seat.Role = &role.String
		}
		seats = append(seats, seat)
	}
	return seats, rows.Err()
}

func (s *Store) ListModeQueuesByModeID(ctx context.Context, modeID uuid.UUID) ([]ModeQueue, error) {
	return listModeQueues(ctx, s.db, modeID)
}

func (s *Store) CountWaitingInModeQueue(ctx context.Context, modeQueueID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM game_queues
		WHERE mode_queue_id = $1 AND status = 'waiting'
	`, modeQueueID).Scan(&count)
	return count, err
}

func (s *Store) GetGameBySlug(ctx context.Context, slug string) (*Game, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+gameColumns+`
		FROM games
		WHERE slug = $1
	`, strings.TrimSpace(slug))
	return scanGame(row)
}

func listGameModes(ctx context.Context, q sqlQueryRowContext, gameID uuid.UUID) ([]GameMode, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT id, game_id, mode_key, display_name, min_players, max_players, status, created_at, updated_at
		FROM game_modes
		WHERE game_id = $1
		ORDER BY mode_key ASC
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGameModes(rows)
}

type sqlQueryRowContext interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func listGameModesTx(ctx context.Context, tx *sql.Tx, gameID uuid.UUID) ([]GameMode, error) {
	return listGameModes(ctx, tx, gameID)
}

func scanGameModes(rows *sql.Rows) ([]GameMode, error) {
	var modes []GameMode
	for rows.Next() {
		var mode GameMode
		if err := rows.Scan(
			&mode.ID, &mode.GameID, &mode.ModeKey, &mode.DisplayName,
			&mode.MinPlayers, &mode.MaxPlayers, &mode.Status,
			&mode.CreatedAt, &mode.UpdatedAt,
		); err != nil {
			return nil, err
		}
		modes = append(modes, mode)
	}
	return modes, rows.Err()
}

func upsertGameModeTx(ctx context.Context, tx *sql.Tx, gameID uuid.UUID, modeKey, displayName string, minPlayers, maxPlayers int) (*GameMode, error) {
	row := tx.QueryRowContext(ctx, `
		INSERT INTO game_modes (game_id, mode_key, display_name, min_players, max_players, status)
		VALUES ($1, $2, $3, $4, $5, 'active')
		ON CONFLICT (game_id, mode_key) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			min_players = EXCLUDED.min_players,
			max_players = EXCLUDED.max_players,
			status = 'active',
			updated_at = NOW()
		RETURNING id, game_id, mode_key, display_name, min_players, max_players, status, created_at, updated_at
	`, gameID, modeKey, displayName, minPlayers, maxPlayers)
	return scanGameModeRow(row)
}

func scanGameModeRow(row *sql.Row) (*GameMode, error) {
	var mode GameMode
	if err := row.Scan(
		&mode.ID, &mode.GameID, &mode.ModeKey, &mode.DisplayName,
		&mode.MinPlayers, &mode.MaxPlayers, &mode.Status,
		&mode.CreatedAt, &mode.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &mode, nil
}

func replaceModeSeatsTx(ctx context.Context, tx *sql.Tx, modeID uuid.UUID, seats []gameclient.SeatManifest) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM game_mode_seats WHERE mode_id = $1`, modeID); err != nil {
		return err
	}
	for i, seat := range seats {
		team := optionalString(seat.Team)
		role := optionalString(seat.Role)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO game_mode_seats (mode_id, seat_key, team, role, sort_order)
			VALUES ($1, $2, $3, $4, $5)
		`, modeID, strings.TrimSpace(seat.Key), team, role, i); err != nil {
			return err
		}
	}
	return nil
}

func listModeSeatKeysTx(ctx context.Context, tx *sql.Tx, modeID uuid.UUID) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT seat_key
		FROM game_mode_seats
		WHERE mode_id = $1
		ORDER BY sort_order ASC, seat_key ASC
	`, modeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func ensureDefaultModeQueueTx(ctx context.Context, tx *sql.Tx, modeID uuid.UUID, playersToStart int) (*ModeQueue, error) {
	row := tx.QueryRowContext(ctx, `
		INSERT INTO mode_queues (mode_id, name, players_to_start, status, is_default)
		VALUES ($1, $2, $3, 'active', true)
		ON CONFLICT (mode_id, name) DO UPDATE SET
			players_to_start = EXCLUDED.players_to_start,
			status = 'active',
			updated_at = NOW()
		RETURNING id, mode_id, name, players_to_start, status, is_default, created_at, updated_at
	`, modeID, defaultQueueName, playersToStart)
	return scanModeQueueRow(row)
}

func scanModeQueueRow(row *sql.Row) (*ModeQueue, error) {
	var q ModeQueue
	if err := row.Scan(
		&q.ID, &q.ModeID, &q.Name, &q.PlayersToStart, &q.Status, &q.IsDefault, &q.CreatedAt, &q.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &q, nil
}

func listModeQueues(ctx context.Context, q sqlQueryRowContext, modeID uuid.UUID) ([]ModeQueue, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT id, mode_id, name, players_to_start, status, is_default, created_at, updated_at
		FROM mode_queues
		WHERE mode_id = $1
		ORDER BY is_default DESC, name ASC
	`, modeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queues []ModeQueue
	for rows.Next() {
		var queue ModeQueue
		if err := rows.Scan(
			&queue.ID, &queue.ModeID, &queue.Name, &queue.PlayersToStart, &queue.Status,
			&queue.IsDefault, &queue.CreatedAt, &queue.UpdatedAt,
		); err != nil {
			return nil, err
		}
		queues = append(queues, queue)
	}
	return queues, rows.Err()
}

func setGameModeStatusTx(ctx context.Context, tx *sql.Tx, modeID uuid.UUID, status string) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE game_modes SET status = $2, updated_at = NOW() WHERE id = $1
	`, modeID, status)
	return err
}

func deactivateModeQueuesTx(ctx context.Context, tx *sql.Tx, modeID uuid.UUID, message string) ([]KickedWaiter, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id FROM mode_queues WHERE mode_id = $1 AND status = 'active'
	`, modeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kicked []KickedWaiter
	for rows.Next() {
		var queueID uuid.UUID
		if err := rows.Scan(&queueID); err != nil {
			return nil, err
		}
		k, err := kickModeQueueWaitersTx(ctx, tx, queueID, message)
		if err != nil {
			return nil, err
		}
		kicked = append(kicked, k...)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE mode_queues SET status = 'inactive', updated_at = NOW()
		WHERE mode_id = $1 AND status = 'active'
	`, modeID)
	return kicked, err
}

func kickModeQueueWaitersTx(ctx context.Context, tx *sql.Tx, modeQueueID uuid.UUID, message string) ([]KickedWaiter, error) {
	rows, err := tx.QueryContext(ctx, `
		UPDATE game_queues
		SET status = 'cancelled'
		WHERE mode_queue_id = $1 AND status = 'waiting'
		RETURNING user_id, game_id
	`, modeQueueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kicked []KickedWaiter
	for rows.Next() {
		var k KickedWaiter
		k.ModeQueueID = modeQueueID
		k.Message = message
		if err := rows.Scan(&k.UserID, &k.GameID); err != nil {
			return nil, err
		}
		kicked = append(kicked, k)
	}
	return kicked, rows.Err()
}

func generateWebhookSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("store: generate webhook secret: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func seatKeysFromManifest(seats []gameclient.SeatManifest) []string {
	keys := make([]string, len(seats))
	for i, seat := range seats {
		keys[i] = strings.TrimSpace(seat.Key)
	}
	sort.Strings(keys)
	return keys
}

func equalStringSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aa := append([]string(nil), a...)
	bb := append([]string(nil), b...)
	sort.Strings(aa)
	sort.Strings(bb)
	for i := range aa {
		if aa[i] != bb[i] {
			return false
		}
	}
	return true
}

func optionalString(value *string) any {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
