package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSignAndVerifyUserToken(t *testing.T) {
	signer, err := LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("LoadSignerFromEnv failed: %v", err)
	}

	userID := uuid.New()
	token, err := signer.SignUserToken(userID, time.Hour)
	if err != nil {
		t.Fatalf("SignUserToken failed: %v", err)
	}

	parsedID, err := signer.VerifyUserToken(token)
	if err != nil {
		t.Fatalf("VerifyUserToken failed: %v", err)
	}
	if parsedID != userID {
		t.Fatalf("expected %s, got %s", userID, parsedID)
	}
}
