package auth

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func openAuthTestService(t *testing.T) *Service {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping auth integration test")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("ping database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	signer, err := LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("load signer: %v", err)
	}

	service, err := NewService(store.New(db), signer)
	if err != nil {
		t.Fatalf("new auth service: %v", err)
	}

	return service
}

func TestMagicLinkLoginFlow(t *testing.T) {
	service := openAuthTestService(t)
	ctx := context.Background()

	email := "auth-" + uuid.NewString() + "@example.com"
	if err := service.RequestMagicLink(ctx, email); err != nil {
		t.Fatalf("RequestMagicLink failed: %v", err)
	}

	link, err := service.store.GetLatestMagicLinkByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetLatestMagicLinkByEmail failed: %v", err)
	}

	user, sessionToken, err := service.CompleteMagicLogin(ctx, link.Token)
	if err != nil {
		t.Fatalf("CompleteMagicLogin failed: %v", err)
	}
	if user.Email != email {
		t.Fatalf("expected email %q, got %q", email, user.Email)
	}
	if sessionToken == "" {
		t.Fatal("expected session token")
	}

	userID, err := service.signer.VerifyUserToken(sessionToken)
	if err != nil {
		t.Fatalf("verify session token: %v", err)
	}
	if userID != user.ID {
		t.Fatalf("expected user id %s in token, got %s", user.ID, userID)
	}

	authCtx := WithUserID(ctx, user.ID.String())
	found, err := service.GetAuthenticatedUser(authCtx)
	if err != nil {
		t.Fatalf("GetAuthenticatedUser failed: %v", err)
	}
	if found == nil || found.ID != user.ID {
		t.Fatalf("expected authenticated user %s, got %+v", user.ID, found)
	}

	if _, _, err := service.CompleteMagicLogin(ctx, link.Token); err == nil {
		t.Fatal("expected reused magic link to fail")
	}
}

func TestRequestMagicLinkRejectsInvalidEmail(t *testing.T) {
	service := openAuthTestService(t)
	err := service.RequestMagicLink(context.Background(), "not-an-email")
	if err != ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}
