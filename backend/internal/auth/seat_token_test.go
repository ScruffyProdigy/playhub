package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestSignSeatTokenClaims(t *testing.T) {
	signer, err := LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("load signer: %v", err)
	}

	userID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	token, err := signer.SignSeatToken(userID, "session-abc", "a", "Alice", time.Hour)
	if err != nil {
		t.Fatalf("SignSeatToken: %v", err)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return signer.publicKey, nil
	})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected map claims")
	}
	if claims["sub"] != userID.String() {
		t.Fatalf("sub: %v", claims["sub"])
	}
	if claims["matchId"] != "session-abc" {
		t.Fatalf("matchId: %v", claims["matchId"])
	}
	if claims["seatKey"] != "a" {
		t.Fatalf("seatKey: %v", claims["seatKey"])
	}
	if claims["name"] != "Alice" {
		t.Fatalf("name: %v", claims["name"])
	}
}
