package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

const gameServiceTokenVersion = "v1"

// GameServiceTokenPepper returns the HMAC key for per-game service tokens.
// Falls back to LOBBY_GAME_SERVICE_TOKEN for local dev when pepper is unset.
func GameServiceTokenPepper() string {
	if p := strings.TrimSpace(os.Getenv("LOBBY_GAME_TOKEN_PEPPER")); p != "" {
		return p
	}
	return strings.TrimSpace(os.Getenv("LOBBY_GAME_SERVICE_TOKEN"))
}

// FormatGameServiceToken returns a scoped credential for one catalog game.
// Format: v1.{gameUUID}.{hex_hmac}
func FormatGameServiceToken(gameID uuid.UUID) (string, error) {
	pepper := GameServiceTokenPepper()
	if pepper == "" {
		return "", errors.New("auth: LOBBY_GAME_TOKEN_PEPPER or LOBBY_GAME_SERVICE_TOKEN is required")
	}
	mac := hmac.New(sha256.New, []byte(pepper))
	_, _ = mac.Write([]byte(gameServiceTokenVersion + ":" + gameID.String()))
	return fmt.Sprintf("%s.%s.%s", gameServiceTokenVersion, gameID.String(), hex.EncodeToString(mac.Sum(nil))), nil
}

// ParseGameServiceToken validates a per-game token and returns the game ID.
func ParseGameServiceToken(token string) (uuid.UUID, error) {
	token = strings.TrimSpace(token)
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] != gameServiceTokenVersion {
		return uuid.Nil, errors.New("auth: not a per-game service token")
	}
	gameID, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth: invalid game id in token: %w", err)
	}
	pepper := GameServiceTokenPepper()
	if pepper == "" {
		return uuid.Nil, errors.New("auth: game service token pepper not configured")
	}
	expected, err := FormatGameServiceToken(gameID)
	if err != nil {
		return uuid.Nil, err
	}
	if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
		return uuid.Nil, errors.New("auth: invalid game service token")
	}
	return gameID, nil
}

// MatchesGameServiceToken accepts per-game tokens or the legacy global token (dev).
func MatchesGameServiceToken(token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}
	if _, err := ParseGameServiceToken(token); err == nil {
		return true
	}
	expected := GameServiceTokenFromEnv()
	if expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}
