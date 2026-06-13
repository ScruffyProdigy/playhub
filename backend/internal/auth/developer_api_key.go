package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const developerAPIKeyPrefix = "lq_dev_"

// DeveloperAPIKeyVerifier validates opaque developer API keys (lq_dev_*).
type DeveloperAPIKeyVerifier interface {
	VerifyDeveloperAPIKey(ctx context.Context, token string) (uuid.UUID, error)
}

// MatchesDeveloperAPIKey reports whether token looks like a developer API key.
func MatchesDeveloperAPIKey(token string) bool {
	return strings.HasPrefix(strings.TrimSpace(token), developerAPIKeyPrefix)
}

// ValidateDeveloperAPIKeyFormat checks token shape before DB lookup.
func ValidateDeveloperAPIKeyFormat(token string) error {
	token = strings.TrimSpace(token)
	if !MatchesDeveloperAPIKey(token) {
		return errors.New("auth: not a developer api key")
	}
	if len(token) < len(developerAPIKeyPrefix)+16 {
		return errors.New("auth: invalid developer api key")
	}
	return nil
}
