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

const MaxLoginCodeAttempts = 5

func scanMagicLink(row interface{ Scan(dest ...any) error }) (*MagicLink, error) {
	var link MagicLink
	var userID sql.NullString
	var tokenHash sql.NullString
	var codeHash sql.NullString
	var usedAt sql.NullTime
	if err := row.Scan(
		&link.ID, &userID, &link.Email, &tokenHash, &link.FailedAttempts,
		&codeHash, &link.ExpiresAt, &usedAt, &link.CreatedAt,
	); err != nil {
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
	if tokenHash.Valid {
		link.TokenHash = tokenHash.String
	}
	if codeHash.Valid {
		link.CodeHash = codeHash.String
	}
	if usedAt.Valid {
		t := usedAt.Time
		link.UsedAt = &t
	}
	return &link, nil
}

const magicLinkColumns = `id, user_id, email, token_hash, failed_attempts, code_hash, expires_at, used_at, created_at`

func (s *Store) CreateMagicLink(ctx context.Context, params CreateMagicLinkParams) (*MagicLink, error) {
	email := strings.ToLower(strings.TrimSpace(params.Email))
	tokenHash := strings.TrimSpace(params.TokenHash)
	if email == "" || tokenHash == "" {
		return nil, fmt.Errorf("store: magic link email and token hash are required")
	}
	if params.ExpiresAt.IsZero() {
		return nil, fmt.Errorf("store: magic link expiry is required")
	}

	var userID any
	if params.UserID != nil {
		userID = *params.UserID
	}

	var codeHash any
	if strings.TrimSpace(params.CodeHash) != "" {
		codeHash = params.CodeHash
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO magic_links (user_id, email, token_hash, code_hash, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+magicLinkColumns+`
	`, userID, email, tokenHash, codeHash, params.ExpiresAt)
	return scanMagicLink(row)
}

// DeleteMagicLink removes an unused magic link row (e.g. when outbound email fails after insert).
func (s *Store) DeleteMagicLink(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM magic_links WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result, ErrNotFound)
}

func (s *Store) CountRecentMagicLinksByEmail(ctx context.Context, email string, since time.Time) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)::int
		FROM magic_links
		WHERE email = $1 AND created_at >= $2
	`, strings.ToLower(strings.TrimSpace(email)), since).Scan(&count)
	return count, err
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

func (s *Store) GetLatestValidMagicLinkByEmail(ctx context.Context, email string, now time.Time) (*MagicLink, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+magicLinkColumns+`
		FROM magic_links
		WHERE email = $1
		  AND used_at IS NULL
		  AND expires_at > $2
		  AND failed_attempts < $3
		ORDER BY created_at DESC
		LIMIT 1
	`, strings.ToLower(strings.TrimSpace(email)), now, MaxLoginCodeAttempts)
	return scanMagicLink(row)
}

// ConsumeMagicLinkByTokenHash atomically marks a link used (URL token path).
func (s *Store) ConsumeMagicLinkByTokenHash(ctx context.Context, tokenHash string, now time.Time) (*MagicLink, error) {
	row := s.db.QueryRowContext(ctx, `
		UPDATE magic_links
		SET used_at = NOW()
		WHERE token_hash = $1
		  AND used_at IS NULL
		  AND expires_at > $2
		  AND failed_attempts < $3
		RETURNING `+magicLinkColumns+`
	`, strings.TrimSpace(tokenHash), now, MaxLoginCodeAttempts)
	link, err := scanMagicLink(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return link, nil
}

// ConsumeMagicLinkByID atomically marks a link used (login-code path after verification).
func (s *Store) ConsumeMagicLinkByID(ctx context.Context, id uuid.UUID, now time.Time) (*MagicLink, error) {
	row := s.db.QueryRowContext(ctx, `
		UPDATE magic_links
		SET used_at = NOW()
		WHERE id = $1
		  AND used_at IS NULL
		  AND expires_at > $2
		  AND failed_attempts < $3
		RETURNING `+magicLinkColumns+`
	`, id, now, MaxLoginCodeAttempts)
	link, err := scanMagicLink(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return link, nil
}

func (s *Store) IncrementMagicLinkFailedAttempts(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE magic_links
		SET failed_attempts = failed_attempts + 1
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
	if link.FailedAttempts >= MaxLoginCodeAttempts {
		return false
	}
	return now.Before(link.ExpiresAt)
}
