package runtimeenv

import (
	"os"
	"strings"
)

// IsProductionEnv reports whether the backend runs in a production-like environment.
func IsProductionEnv() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "production", "prod":
		return true
	}
	pub := strings.ToLower(strings.TrimSpace(os.Getenv("LOBBY_PUBLIC_URL")))
	return strings.HasPrefix(pub, "https://")
}
