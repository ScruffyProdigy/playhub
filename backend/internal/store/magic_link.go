package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func scanMagicLink(row interface{ Scan(dest ...any) error }) (*MagicLink, error) {
	var link MagicLink
	var userID sql.NullString
	var usedAt sql.NullTime
	if err := row.Scan(&link.ID, &userID, &link.Email, &link.Token, &link.ExpiresAt, &usedAt, &link.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if userID.Valid {
		id, err := uuid.Parse(userID.String)
		if err != nil {
			return nil, err
		}
		link.UserID = &id
	}
	if usedAt.Valid {
		t := usedAt.Time
		link.UsedAt = &t
	}
	return &link, nil
}

const magicLinkColumns = `id, user_id, email, token, expires_at, used_at, created_at`

func (s *Store) CreateMagicLink(ctx context.Context, params CreateMagicLinkParams) (*MagicLink, error) {
	email := strings.ToLower(strings.TrimSpace(params.Email))
	token := strings.TrimSpace(params.Token)
	if email == "" || token == "" {
		return nil, fmt.Errorf("store: magic link email and token are required")
	}
	if params.ExpiresAt.IsZero() {
		return nil, fmt.Errorf("store: magic link expiry is required")
	}

	var userID any
	if params.UserID != nil {
		userID = *params.UserID
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO magic_links (user_id, email, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING `+magicLinkColumns+`
	`, userID, email, token, params.ExpiresAt)
	return scanMagicLink(row)
}

func (s *Store) GetMagicLinkByToken(ctx context.Context, token string) (*MagicLink, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+magicLinkColumns+`
		FROM magic_links
		WHERE token = $1
	`, strings.TrimSpace(token))
	return scanMagicLink(row)
}

func (s *Store) GetLatestMagicLinkByEmail(ctx context.Context, email string) (*MagicLink, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+magicLinkColumns+`
		FROM magic_links
		WHERE email = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, strings.ToLower(strings.TrimSpace(email)))
	return scanMagicLink(row)
}

func (s *Store) MarkMagicLinkUsed(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE magic_links
		SET used_at = NOW()
		WHERE id = $1 AND used_at IS NULL
	`, id)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result, ErrNotFound)
}

func (s *Store) IsMagicLinkValid(link *MagicLink, now time.Time) bool {
	if link == nil || link.UsedAt != nil {
		return false
	}
	return now.Before(link.ExpiresAt)
}
