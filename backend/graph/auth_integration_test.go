package graph

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/graph/generated"
	"github.com/scruffyprodigy/playhub/internal/auth"
	"github.com/scruffyprodigy/playhub/internal/store"
	_ "github.com/lib/pq"
)

func newAuthGraphQLTestClient(t *testing.T) (*client.Client, http.Handler, *store.Store) {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping GraphQL auth integration test")
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

	st := store.New(db)
	signer, err := auth.LoadSignerFromEnv()
	if err != nil {
		t.Fatalf("load signer: %v", err)
	}

	authService, err := auth.NewService(st, signer)
	if err != nil {
		t.Fatalf("new auth service: %v", err)
	}

	resolver := NewResolver(st, authService)
	gql := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	handlerWithAuth := auth.Middleware(signer, gql)
	c := client.New(handlerWithAuth)

	return c, handlerWithAuth, st
}

func postGraphQLWithCookies(t *testing.T, handler http.Handler, query string, variables map[string]any) ([]*http.Cookie, []byte) {
	t.Helper()

	payload := map[string]any{"query": query}
	if variables != nil {
		payload["variables"] = variables
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal GraphQL request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code >= http.StatusBadRequest {
		t.Fatalf("GraphQL request failed with HTTP %d: %s", rec.Code, rec.Body.String())
	}

	return rec.Result().Cookies(), rec.Body.Bytes()
}

func sessionCookieOptions(cookies []*http.Cookie) []client.Option {
	cfg := auth.CookieConfigFromEnv()
	opts := make([]client.Option, 0, len(cookies))
	for _, cookie := range cookies {
		if cookie.Name == cfg.Name && cookie.Value != "" {
			opts = append(opts, client.AddCookie(cookie))
		}
	}
	return opts
}

func TestAuthGraphQLFlow(t *testing.T) {
	c, handlerWithAuth, st := newAuthGraphQLTestClient(t)
	ctx := context.Background()

	email := "graphql-auth-" + uuid.NewString() + "@example.com"

	var loginResp struct {
		LoginMagic bool `json:"loginMagic"`
	}
	if err := c.Post(`mutation LoginMagic($email: String!) {
		loginMagic(email: $email)
	}`, &loginResp, client.Var("email", email)); err != nil {
		t.Fatalf("loginMagic failed: %v", err)
	}
	if !loginResp.LoginMagic {
		t.Fatal("expected loginMagic to return true")
	}

	link, err := st.GetLatestMagicLinkByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetLatestMagicLinkByEmail failed: %v", err)
	}

	completeQuery := `mutation CompleteMagic($token: ID!) {
		completeMagic(token: $token) {
			email
			displayName
		}
	}`
	sessionCookies, completeBody := postGraphQLWithCookies(t, handlerWithAuth, completeQuery, map[string]any{"token": link.Token})
	if len(sessionCookieOptions(sessionCookies)) == 0 {
		t.Fatal("expected completeMagic to set a session cookie")
	}

	var completeResp struct {
		Data struct {
			CompleteMagic struct {
				Email       *string `json:"email"`
				DisplayName *string `json:"displayName"`
			} `json:"completeMagic"`
		} `json:"data"`
	}
	if err := json.Unmarshal(completeBody, &completeResp); err != nil {
		t.Fatalf("decode completeMagic response: %v", err)
	}
	if completeResp.Data.CompleteMagic.Email == nil || *completeResp.Data.CompleteMagic.Email != email {
		t.Fatalf("expected completed user email %q, got %+v", email, completeResp.Data.CompleteMagic.Email)
	}
	if completeResp.Data.CompleteMagic.DisplayName == nil || *completeResp.Data.CompleteMagic.DisplayName != store.DefaultDisplayName(email) {
		t.Fatalf("expected provisional display name, got %+v", completeResp.Data.CompleteMagic.DisplayName)
	}

	var meResp struct {
		Me *struct {
			Email *string `json:"email"`
		} `json:"me"`
	}
	if err := c.Post(`query {
		me {
			email
		}
	}`, &meResp, sessionCookieOptions(sessionCookies)...); err != nil {
		t.Fatalf("me query failed: %v", err)
	}
	if meResp.Me == nil || meResp.Me.Email == nil || *meResp.Me.Email != email {
		t.Fatalf("expected authenticated me for %q, got %+v", email, meResp.Me)
	}

	var logoutResp struct {
		Logout bool `json:"logout"`
	}
	if err := c.Post(`mutation {
		logout
	}`, &logoutResp, sessionCookieOptions(sessionCookies)...); err != nil {
		t.Fatalf("logout failed: %v", err)
	}
	if !logoutResp.Logout {
		t.Fatal("expected logout to return true")
	}

	if err := c.Post(`query {
		me {
			email
		}
	}`, &meResp); err != nil {
		t.Fatalf("me query after logout failed: %v", err)
	}
	if meResp.Me != nil && meResp.Me.Email != nil {
		t.Fatalf("expected me to be unauthenticated after logout, got %+v", meResp.Me)
	}
}
