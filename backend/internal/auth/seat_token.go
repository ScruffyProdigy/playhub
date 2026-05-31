package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const defaultSeatTokenTTL = 2 * time.Hour

// SignSeatToken issues a short-lived JWT for launching into a third-party game.
// Claims: iss, aud, sub, jti, matchId, seatKey, name (optional), nbf, iat, exp.
func (s *Signer) SignSeatToken(userID uuid.UUID, audience, externalMatchID, seatKey, displayName string, ttl time.Duration) (string, error) {
	audience = normalizeAudience(audience)
	if audience == "" {
		return "", fmt.Errorf("auth: audience is required")
	}
	if externalMatchID == "" || seatKey == "" {
		return "", fmt.Errorf("auth: match id and seat key are required")
	}
	if ttl <= 0 {
		ttl = defaultSeatTokenTTL
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":     LobbyIssuer(),
		"aud":     audience,
		"sub":     userID.String(),
		"jti":     uuid.NewString(),
		"matchId": externalMatchID,
		"seatKey": seatKey,
		"nbf":     now.Unix(),
		"iat":     now.Unix(),
		"exp":     now.Add(ttl).Unix(),
	}
	if displayName != "" {
		claims["name"] = displayName
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = s.kid
	return token.SignedString(s.privateKey)
}

func normalizeAudience(audience string) string {
	return strings.TrimRight(strings.TrimSpace(audience), "/")
}
