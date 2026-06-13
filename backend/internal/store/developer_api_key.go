package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const developerAPIKeyPrefix = "lq_dev_"

// DeveloperAPIKey is a persisted developer API key metadata row.
type DeveloperAPIKey struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Name       string
	KeyPrefix  string
	CreatedAt  time.Time
	LastUsedAt *time.Time
}

func developerAPIKeyPepper() string {
	if p := strings.TrimSpace(os.Getenv("DEVELOPER_API_KEY_PEPPER")); p != "" {
		return p
	}
	return strings.TrimSpace(os.Getenv("MAGIC_LINK_PEPPER"))
}

func hashDeveloperAPIKey(token string) string {
	pepper := developerAPIKeyPepper()
	sum := sha256.Sum256([]byte(pepper + strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

// GenerateDeveloperAPIKey creates a new opaque developer API key and its stored hash.
func GenerateDeveloperAPIKey() (raw string, displayPrefix string, keyHash string, err error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", "", "", err
	}
	raw = developerAPIKeyPrefix + hex.EncodeToString(secret)
	displayPrefix = raw[:len(developerAPIKeyPrefix)+8]
	keyHash = hashDeveloperAPIKey(raw)
	return raw, displayPrefix, keyHash, nil
}

// CreateDeveloperAPIKey stores a hashed developer API key and returns metadata plus the raw secret once.
func (s *Store) CreateDeveloperAPIKey(ctx context.Context, userID uuid.UUID, name string) (rawKey string, key DeveloperAPIKey, err error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Integration agent"
	}

	rawKey, prefix, keyHash, err := GenerateDeveloperAPIKey()
	if err != nil {
		return "", DeveloperAPIKey{}, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO developer_api_keys (user_id, name, key_prefix, key_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, key_prefix, created_at, last_used_at
	`, userID, name, prefix, keyHash)

	key, err = scanDeveloperAPIKey(row)
	if err != nil {
		return "", DeveloperAPIKey{}, err
	}
	return rawKey, key, nil
}

// ListDeveloperAPIKeys returns active keys for a user.
func (s *Store) ListDeveloperAPIKeys(ctx context.Context, userID uuid.UUID) ([]DeveloperAPIKey, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, name, key_prefix, created_at, last_used_at
		FROM developer_api_keys
		WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []DeveloperAPIKey
	for rows.Next() {
		key, err := scanDeveloperAPIKey(rows)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

// RevokeDeveloperAPIKey marks a key revoked when owned by userID.
func (s *Store) RevokeDeveloperAPIKey(ctx context.Context, keyID, userID uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE developer_api_keys
		SET revoked_at = now()
		WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL
	`, keyID, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// VerifyDeveloperAPIKey implements auth.DeveloperAPIKeyVerifier.
func (s *Store) VerifyDeveloperAPIKey(ctx context.Context, token string) (uuid.UUID, error) {
	token = strings.TrimSpace(token)
	if !strings.HasPrefix(token, developerAPIKeyPrefix) {
		return uuid.Nil, fmt.Errorf("auth: not a developer api key")
	}
	if len(token) < len(developerAPIKeyPrefix)+16 {
		return uuid.Nil, fmt.Errorf("auth: invalid developer api key")
	}

	keyHash := hashDeveloperAPIKey(token)
	var userID uuid.UUID
	err := s.db.QueryRowContext(ctx, `
		UPDATE developer_api_keys
		SET last_used_at = now()
		WHERE key_hash = $1 AND revoked_at IS NULL
		RETURNING user_id
	`, keyHash).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("auth: invalid developer api key")
	}
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func scanDeveloperAPIKey(row interface {
	Scan(dest ...any) error
}) (DeveloperAPIKey, error) {
	var key DeveloperAPIKey
	var lastUsed sql.NullTime
	if err := row.Scan(&key.ID, &key.UserID, &key.Name, &key.KeyPrefix, &key.CreatedAt, &lastUsed); err != nil {
		return DeveloperAPIKey{}, err
	}
	if lastUsed.Valid {
		key.LastUsedAt = &lastUsed.Time
	}
	return key, nil
}
