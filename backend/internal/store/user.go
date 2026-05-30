package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scanUser(row interface{ Scan(dest ...any) error }) (*User, error) {
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.DisplayName, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

const userColumns = `id, email, username, display_name, created_at`

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+userColumns+`
		FROM users
		WHERE id = $1 AND is_active = true
	`, id)
	return scanUser(row)
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+userColumns+`
		FROM users
		WHERE email = $1 AND is_active = true
	`, strings.ToLower(strings.TrimSpace(email)))
	return scanUser(row)
}

func (s *Store) CreateUser(ctx context.Context, params CreateUserParams) (*User, error) {
	email := strings.ToLower(strings.TrimSpace(params.Email))
	displayName := strings.TrimSpace(params.DisplayName)
	if displayName == "" {
		displayName = strings.Split(email, "@")[0]
	}

	username, err := s.uniqueUsername(ctx, email)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO users (email, username, display_name)
		VALUES ($1, $2, $3)
		RETURNING `+userColumns+`
	`, email, username, displayName)
	return scanUser(row)
}

func (s *Store) UpdateUserLastLogin(ctx context.Context, id uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return ensureRowsAffected(result, ErrNotFound)
}

func (s *Store) uniqueUsername(ctx context.Context, email string) (string, error) {
	base := strings.Split(email, "@")[0]
	base = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, strings.ToLower(base))
	if base == "" {
		base = "user"
	}
	if len(base) > 40 {
		base = base[:40]
	}

	for i := 0; i < 5; i++ {
		candidate := base
		if i > 0 {
			suffix := strings.ReplaceAll(uuid.NewString(), "-", "")[:8]
			candidate = fmt.Sprintf("%s_%s", base, suffix)
		}

		var exists bool
		if err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, candidate).Scan(&exists); err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("store: failed to generate unique username for %s", email)
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func ensureRowsAffected(result sql.Result, notFound error) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return notFound
	}
	return nil
}
