package auth

import (
	"fmt"
	"net/url"
	"strings"
)

// MagicLinkTokenFromURL extracts the token query parameter from a sign-in email link.
func MagicLinkTokenFromURL(link string) (string, error) {
	link = strings.TrimSpace(link)
	if link == "" {
		return "", fmt.Errorf("auth: magic link URL is empty")
	}
	u, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("auth: parse magic link URL: %w", err)
	}
	token := strings.TrimSpace(u.Query().Get("token"))
	if token == "" {
		return "", fmt.Errorf("auth: magic link URL missing token query param")
	}
	return token, nil
}
