package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/gameclient"
)

var gameSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type RegisterMyGameParams struct {
	OwnerUserID      uuid.UUID
	Slug             string
	Name             string
	ShortDescription string
	APIBaseURL       string
	ContactEmail     string
	WebsiteURL       *string
	CommunityURL     *string
}

type RegisterMyGameResult struct {
	Game          *Game
	WebhookSecret string
	Connected     bool
	ConnectError  string
}

type GameIntegrationCheck struct {
	CheckID    string
	Status     string
	Message    *string
	DetailJSON json.RawMessage
	RanAt      time.Time
}

func ValidateGameSlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("store: slug is required")
	}
	if len(slug) > 64 {
		return fmt.Errorf("store: slug must be at most 64 characters")
	}
	if !gameSlugPattern.MatchString(slug) {
		return fmt.Errorf("store: slug must be lowercase letters, numbers, and hyphens")
	}
	return nil
}

// RegisterMyGame saves a developer-owned game and optionally connects the API.
func (s *Store) RegisterMyGame(ctx context.Context, params RegisterMyGameParams, manifest *gameclient.Manifest, connectErr error) (*RegisterMyGameResult, error) {
	slug := strings.TrimSpace(params.Slug)
	if err := ValidateGameSlug(slug); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(params.Name)
	if name == "" {
		return nil, fmt.Errorf("store: name is required")
	}
	shortDescription := strings.TrimSpace(params.ShortDescription)
	if shortDescription == "" {
		return nil, fmt.Errorf("store: short description is required")
	}
	apiBaseURL := strings.TrimRight(strings.TrimSpace(params.APIBaseURL), "/")
	if apiBaseURL == "" {
		return nil, fmt.Errorf("store: api base URL is required")
	}
	contactEmail := strings.TrimSpace(params.ContactEmail)
	if contactEmail == "" {
		return nil, fmt.Errorf("store: contact email is required")
	}

	secret, err := generateWebhookSecret()
	if err != nil {
		return nil, err
	}

	visibility := GameVisibilityDraft
	status := ModeStatusInactive
	if manifest != nil && connectErr == nil {
		visibility = GameVisibilityPrivateTesting
		status = ModeStatusActive
	}

	iconURL := DefaultGameIconURL(slug)
	heroURL := DefaultGameHeroURL(slug)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO games (
			name, short_description, icon_url, hero_url, slug, api_base_url,
			category, status, visibility, owner_user_id, contact_email, website_url, community_url, webhook_secret
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'catalog', $7, $8, $9, $10, $11, $12, $13)
		RETURNING `+gameColumns+`
	`, name, shortDescription, iconURL, heroURL, slug, apiBaseURL, status, visibility, params.OwnerUserID, contactEmail, params.WebsiteURL, params.CommunityURL, secret)
	game, err := scanGame(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("store: game slug %q already exists", slug)
		}
		return nil, err
	}

	result := &RegisterMyGameResult{
		Game:          game,
		WebhookSecret: secret,
		Connected:     manifest != nil && connectErr == nil,
	}
	if connectErr != nil {
		result.ConnectError = connectErr.Error()
	}

	if manifest != nil && connectErr == nil {
		if _, _, err := applyManifestTx(ctx, tx, game, manifest); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	game, err = s.GetGameByID(ctx, game.ID)
	if err != nil {
		return nil, err
	}
	result.Game = game
	return result, nil
}

// ListMyGames returns games owned by the given user, newest first.
func (s *Store) ListMyGames(ctx context.Context, ownerUserID uuid.UUID) ([]Game, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+gameColumns+`
		FROM games
		WHERE owner_user_id = $1
		ORDER BY created_at DESC
	`, ownerUserID)
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

// GetOwnedGame loads a game owned by ownerUserID.
func (s *Store) GetOwnedGame(ctx context.Context, gameID, ownerUserID uuid.UUID) (*Game, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+gameColumns+`
		FROM games
		WHERE id = $1 AND owner_user_id = $2
	`, gameID, ownerUserID)
	return scanGame(row)
}

// AssertCanCreateTableForGame enforces visibility rules for table creation.
func (s *Store) AssertCanCreateTableForGame(ctx context.Context, gameID, userID uuid.UUID) error {
	game, err := s.GetGameByID(ctx, gameID)
	if err != nil {
		return err
	}
	if !game.AllowsRoomTables() {
		return fmt.Errorf("store: game is not ready for testing yet — connect your API first")
	}
	if game.IsPublicCatalog() {
		return nil
	}
	if game.OwnerUserID != nil && *game.OwnerUserID == userID {
		return nil
	}
	return fmt.Errorf("store: only the game owner can create tables for this game")
}

// UpsertIntegrationCheck stores the latest result for a checklist row.
func (s *Store) UpsertIntegrationCheck(ctx context.Context, gameID uuid.UUID, checkID, status string, message *string, detail json.RawMessage) (*GameIntegrationCheck, error) {
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO game_integration_checks (game_id, check_id, status, message, detail_json, ran_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (game_id, check_id) DO UPDATE SET
			status = EXCLUDED.status,
			message = EXCLUDED.message,
			detail_json = EXCLUDED.detail_json,
			ran_at = NOW()
		RETURNING check_id, status, message, detail_json, ran_at
	`, gameID, checkID, status, message, nullableJSON(detail))
	return scanIntegrationCheck(row)
}

func nullableJSON(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	return raw
}

func scanIntegrationCheck(row interface{ Scan(dest ...any) error }) (*GameIntegrationCheck, error) {
	var check GameIntegrationCheck
	var message sql.NullString
	var detail []byte
	if err := row.Scan(&check.CheckID, &check.Status, &message, &detail, &check.RanAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if message.Valid {
		check.Message = &message.String
	}
	if len(detail) > 0 {
		check.DetailJSON = json.RawMessage(detail)
	}
	return &check, nil
}

// ListIntegrationChecks returns stored checklist results for a game.
func (s *Store) ListIntegrationChecks(ctx context.Context, gameID uuid.UUID) ([]GameIntegrationCheck, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT check_id, status, message, detail_json, ran_at
		FROM game_integration_checks
		WHERE game_id = $1
		ORDER BY check_id ASC
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []GameIntegrationCheck
	for rows.Next() {
		check, err := scanIntegrationCheck(rows)
		if err != nil {
			return nil, err
		}
		checks = append(checks, *check)
	}
	return checks, rows.Err()
}
