package auth

import (
	"errors"
	"os"
	"strings"
)

var ErrNotAdmin = errors.New("admin access required")

// AdminEmailsFromEnv returns allowlisted admin emails from LOBBY_ADMIN_EMAILS (comma-separated).
func AdminEmailsFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("LOBBY_ADMIN_EMAILS"))
	if raw == "" {
		return nil
	}
	var admins []string
	for part := range strings.SplitSeq(raw, ",") {
		normalized, err := normalizeEmail(part)
		if err != nil {
			continue
		}
		admins = append(admins, normalized)
	}
	return admins
}

// IsAdminEmail reports whether email is in LOBBY_ADMIN_EMAILS.
func IsAdminEmail(email string) bool {
	normalized, err := normalizeEmail(email)
	if err != nil {
		return false
	}
	for _, admin := range AdminEmailsFromEnv() {
		if normalized == admin {
			return true
		}
	}
	return false
}
