package gameurl

import (
	"fmt"
	"net/url"
	"strings"
)

// AttachSeatToken merges a seat JWT into a game-minted launch URL via query params.
func AttachSeatToken(rawURL, token string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	token = strings.TrimSpace(token)
	if rawURL == "" {
		return "", fmt.Errorf("gameurl: launch URL is required")
	}
	if token == "" {
		return "", fmt.Errorf("gameurl: seat token is required")
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("gameurl: invalid launch URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("gameurl: launch URL scheme must be http or https")
	}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// SameOriginHost reports whether two URLs share the same host (case-insensitive).
func SameOriginHost(a, b string) bool {
	ua, errA := url.Parse(strings.TrimSpace(a))
	ub, errB := url.Parse(strings.TrimSpace(b))
	if errA != nil || errB != nil {
		return false
	}
	return strings.EqualFold(ua.Hostname(), ub.Hostname())
}
