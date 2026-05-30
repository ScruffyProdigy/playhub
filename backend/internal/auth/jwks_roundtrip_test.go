package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestDevJWTKeyPersistsAcrossLoad(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "dev.pem")
	t.Setenv("JWT_DEV_KEY_FILE", keyPath)
	t.Setenv("JWT_PRIVATE_KEY_PEM", "")

	first, err := LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("first load: %v", err)
	}
	second, err := LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("second load: %v", err)
	}
	if first.publicX != second.publicX {
		t.Fatalf("expected same public key after reload, got %q vs %q", first.publicX, second.publicX)
	}
}

func TestSeatTokenVerifiedViaJWKSHandler(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("JWT_DEV_KEY_FILE", filepath.Join(dir, "test.pem"))

	signer, err := LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("load signer: %v", err)
	}

	userID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	token, err := signer.SignSeatToken(userID, "match-1", "b", "Bob", time.Hour)
	if err != nil {
		t.Fatalf("SignSeatToken: %v", err)
	}

	srv := httptest.NewServer(signer.JWKSHandler())
	t.Cleanup(srv.Close)

	var jwks struct {
		Keys []map[string]string `json:"keys"`
	}
	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("GET jwks: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("jwks status: %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		t.Fatalf("decode jwks: %v", err)
	}
	if len(jwks.Keys) != 1 || jwks.Keys[0]["kid"] != signer.kid {
		t.Fatalf("unexpected jwks keys: %+v", jwks.Keys)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodEdDSA {
			return nil, jwt.ErrTokenUnverifiable
		}
		if token.Header["kid"] != signer.kid {
			t.Fatalf("token kid: %v", token.Header["kid"])
		}
		return signer.publicKey, nil
	})
	if err != nil {
		t.Fatalf("verify seat token: %v", err)
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		t.Fatal("invalid claims")
	}
	if claims["matchId"] != "match-1" || claims["seatKey"] != "b" {
		t.Fatalf("claims: %+v", claims)
	}
}
