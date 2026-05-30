package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const defaultSeatTokenTTL = 2 * time.Hour

// SignSeatToken issues a short-lived JWT for launching into a third-party game.
// Claims match demo-game-rps: sub, matchId, seatKey, name (optional).
func (s *Signer) SignSeatToken(userID uuid.UUID, externalMatchID, seatKey, displayName string, ttl time.Duration) (string, error) {
	if externalMatchID == "" || seatKey == "" {
		return "", fmt.Errorf("auth: match id and seat key are required")
	}
	if ttl <= 0 {
		ttl = defaultSeatTokenTTL
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":     userID.String(),
		"matchId": externalMatchID,
		"seatKey": seatKey,
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
