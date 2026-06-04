package auth

import "github.com/scruffyprodigy/playhub/internal/runtimeenv"

// IsProductionEnv reports whether the backend runs in a production-like environment.
func IsProductionEnv() bool {
	return runtimeenv.IsProductionEnv()
}
