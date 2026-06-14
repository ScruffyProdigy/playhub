package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func TestSignAndVerifyOAuthState(t *testing.T) {
	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	signer, err := loadOrCreateDevSigner("test-kid")
	if err != nil {
		t.Fatalf("load signer: %v", err)
	}

	userID := uuid.New()
	token, err := signer.SignOAuthState(OAuthState{
		Provider:     "google",
		Mode:         OAuthModeLink,
		UserID:       userID,
		ConfirmMerge: true,
	}, time.Minute)
	if err != nil {
		t.Fatalf("SignOAuthState: %v", err)
	}

	state, err := signer.VerifyOAuthState(token)
	if err != nil {
		t.Fatalf("VerifyOAuthState: %v", err)
	}
	if state.Provider != "google" {
		t.Fatalf("provider = %q", state.Provider)
	}
	if state.Mode != OAuthModeLink {
		t.Fatalf("mode = %q", state.Mode)
	}
	if state.UserID != userID {
		t.Fatalf("user id mismatch")
	}
	if !state.ConfirmMerge {
		t.Fatal("expected confirm merge")
	}
}

func TestParseDiscordProfile(t *testing.T) {
	profile, err := parseDiscordProfile([]byte(`{
		"id": 1515731455333105884,
		"username": "player",
		"global_name": "Player One",
		"email": "player@example.com",
		"verified": true
	}`))
	if err != nil {
		t.Fatalf("parseDiscordProfile: %v", err)
	}
	if profile.Subject != "1515731455333105884" {
		t.Fatalf("subject = %q", profile.Subject)
	}
	if profile.DisplayName != "Player One" {
		t.Fatalf("display name = %q", profile.DisplayName)
	}
}

func TestExchangeDiscordCode(t *testing.T) {
	t.Setenv("LOBBY_ISSUER_URL", "https://joinquest.cc")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		if r.Form.Get("grant_type") != "authorization_code" {
			t.Fatalf("grant_type = %q", r.Form.Get("grant_type"))
		}
		if r.Form.Get("code") != "abc123" {
			t.Fatalf("code = %q", r.Form.Get("code"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"token","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	oauthCfg := &oauth2.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://joinquest.cc/auth/oauth/discord/callback",
		Endpoint: oauth2.Endpoint{
			TokenURL:  server.URL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	token, err := exchangeDiscordCode(context.Background(), oauthCfg, "abc123")
	if err != nil {
		t.Fatalf("exchangeDiscordCode: %v", err)
	}
	if token.AccessToken != "token" {
		t.Fatalf("access token = %q", token.AccessToken)
	}
}

func TestEnabledOAuthProviders(t *testing.T) {
	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "")
	t.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "")
	t.Setenv("DISCORD_OAUTH_CLIENT_ID", "")
	t.Setenv("DISCORD_OAUTH_CLIENT_SECRET", "")
	if len(EnabledOAuthProviders()) != 0 {
		t.Fatal("expected no providers")
	}

	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "id")
	t.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "secret")
	providers := EnabledOAuthProviders()
	if len(providers) != 1 || providers[0] != "google" {
		t.Fatalf("providers = %#v", providers)
	}
}
