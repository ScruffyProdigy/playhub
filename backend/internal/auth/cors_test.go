package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddlewareAllowsDevOrigin(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	CORSMiddleware(next).ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("expected allowed origin header, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("expected credentials header, got %q", got)
	}
}

func TestWebSocketOriginAllowedLoopback(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	if !loopbackOriginsEquivalent("http://localhost:5173", "http://127.0.0.1:5173") {
		t.Fatal("expected localhost and 127.0.0.1 origins to be equivalent")
	}

	for _, origin := range []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://[::1]:5173",
		"",
	} {
		req := httptest.NewRequest(http.MethodGet, "/graphql", nil)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		if !WebSocketOriginAllowed(req) {
			t.Fatalf("expected websocket origin %q to be allowed", origin)
		}
	}
}

func TestCORSMiddlewarePreflight(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not run for preflight")
	})

	req := httptest.NewRequest(http.MethodOptions, "/graphql", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	CORSMiddleware(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
