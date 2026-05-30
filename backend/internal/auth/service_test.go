package auth

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/store"
	"github.com/scruffyprodigy/playhub/internal/testdb"
)

func openAuthTestService(t *testing.T) *Service {
	t.Helper()

	databaseURL := testdb.RequireURL(t)

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
	cleaner := service.store.NewTestCleaner(t)
	ctx := context.Background()

	email := "auth-" + uuid.NewString() + "@example.com"
	cleaner.TrackEmail(email)

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
	cleaner.TrackUser(user.ID)

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

func TestLogoutClearsSessionCookie(t *testing.T) {
	service := openAuthTestService(t)

	recorder := httptest.NewRecorder()
	ctx := WithResponseWriter(context.Background(), recorder)
	service.Logout(ctx)

	cookies := recorder.Result().Cookies()
	cfg := service.CookieConfig()

	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == cfg.Name {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Fatalf("expected %q cookie to be cleared", cfg.Name)
	}
	if sessionCookie.Value != "" {
		t.Fatalf("expected empty session cookie value, got %q", sessionCookie.Value)
	}
	if sessionCookie.MaxAge >= 0 {
		t.Fatalf("expected session cookie MaxAge < 0, got %d", sessionCookie.MaxAge)
	}
}
