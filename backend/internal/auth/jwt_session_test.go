package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestVerifyUserTokenRejectsNonSessionTyp(t *testing.T) {
	signer, err := LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("LoadSignerFromEnv: %v", err)
	}

	userID := uuid.New()
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": LobbyIssuer(),
		"aud": sessionTokenAudience(),
		"typ": "seat",
		"sub": userID.String(),
		"jti": uuid.NewString(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = signer.kid
	raw, err := token.SignedString(signer.privateKey)
	if err != nil {
		t.Fatalf("SignedString: %v", err)
	}

	if _, err := signer.VerifyUserToken(raw); err == nil {
		t.Fatal("expected seat typ token to be rejected")
	}
}
