package returnlink

import (
	"net/url"
	"strings"
)

// AppendMatchID appends match={externalMatchId} to a Lobby return hub URL.
// Canonical server-side helper — game clients should mirror this algorithm
// (see docs/player-return-routing.md).
func AppendMatchID(returnURL, externalMatchID string) string {
	base := strings.TrimSpace(returnURL)
	matchID := strings.TrimSpace(externalMatchID)
	if base == "" {
		if matchID == "" {
			return "/return"
		}
		return "/return?match=" + url.QueryEscape(matchID)
	}
	if matchID == "" {
		return base
	}

	u, err := url.Parse(base)
	if err != nil || u.Scheme == "" || u.Host == "" {
		sep := "?"
		if strings.Contains(base, "?") {
			sep = "&"
		}
		return base + sep + "match=" + url.QueryEscape(matchID)
	}

	q := u.Query()
	q.Set("match", matchID)
	u.RawQuery = q.Encode()
	return u.String()
}
