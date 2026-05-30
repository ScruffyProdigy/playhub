package auth

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

// CORSMiddleware allows browser clients to call the API with credentials.
func CORSMiddleware(next http.Handler) http.Handler {
	allowedOrigins := corsOriginsFromEnv()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" && originAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func corsOriginsFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		return []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		}
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

// AllowedOrigins returns configured browser origins for CORS and websocket checks.
func AllowedOrigins() []string {
	return corsOriginsFromEnv()
}

// WebSocketOriginAllowed reports whether the request Origin may open a GraphQL websocket.
func WebSocketOriginAllowed(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" || origin == "null" {
		return true
	}
	if originAllowed(origin, AllowedOrigins()) {
		return true
	}
	// Allow loopback browser origins (e.g. [::1]:5173) not listed in CORS_ALLOWED_ORIGINS.
	return isLoopbackOrigin(origin)
}

func originAllowed(origin string, allowed []string) bool {
	for _, candidate := range allowed {
		if candidate == "*" || candidate == origin {
			return true
		}
		if loopbackOriginsEquivalent(candidate, origin) {
			return true
		}
	}
	return false
}

// loopbackOriginsEquivalent treats localhost and 127.0.0.1 as the same origin when ports match.
func loopbackOriginsEquivalent(a, b string) bool {
	ua, err := url.Parse(a)
	if err != nil {
		return false
	}
	ub, err := url.Parse(b)
	if err != nil {
		return false
	}
	if ua.Scheme != ub.Scheme {
		return false
	}
	if normalizePort(ua) != normalizePort(ub) {
		return false
	}
	return isLoopbackHost(ua.Hostname()) && isLoopbackHost(ub.Hostname())
}

func normalizePort(u *url.URL) string {
	if port := u.Port(); port != "" {
		return port
	}
	if u.Scheme == "https" {
		return "443"
	}
	return "80"
}

func isLoopbackHost(host string) bool {
	switch strings.ToLower(host) {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		return false
	}
}

func isLoopbackOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return isLoopbackHost(u.Hostname())
}
