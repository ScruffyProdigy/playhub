package store

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// TestCleaner removes rows created during integration tests.
type TestCleaner struct {
	st     *Store
	games  []uuid.UUID
	users  []uuid.UUID
	emails []string
}

// NewTestCleaner registers automatic cleanup when the test finishes.
func (s *Store) NewTestCleaner(t *testing.T) *TestCleaner {
	t.Helper()

	cleaner := &TestCleaner{st: s}
	t.Cleanup(func() {
		cleaner.run(context.Background())
	})
	return cleaner
}

func (c *TestCleaner) TrackGame(id uuid.UUID) {
	c.games = append(c.games, id)
}

func (c *TestCleaner) TrackUser(id uuid.UUID) {
	c.users = append(c.users, id)
}

func (c *TestCleaner) TrackEmail(email string) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return
	}
	c.emails = append(c.emails, email)
}

func (c *TestCleaner) run(ctx context.Context) {
	for _, gameID := range c.games {
		_ = c.st.deleteTestGame(ctx, gameID)
	}
	for _, email := range c.emails {
		_ = c.st.deleteTestMagicLinksByEmail(ctx, email)
	}
	for _, userID := range c.users {
		_ = c.st.deleteTestUser(ctx, userID)
	}
}

func (s *Store) deleteTestGame(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM games
		WHERE id = $1 AND (category IS NULL OR category <> 'demo')
	`, id)
	return err
}

func (s *Store) deleteTestUser(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (s *Store) deleteTestMagicLinksByEmail(ctx context.Context, email string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM magic_links WHERE email = $1`, strings.ToLower(strings.TrimSpace(email)))
	return err
}
