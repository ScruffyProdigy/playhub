package graph

import (
	"context"
	"fmt"

	"github.com/scruffyprodigy/playhub/internal/auth"
)

func requireGameServiceAuth(ctx context.Context) error {
	if auth.GameServiceFromContext(ctx) {
		return nil
	}
	if auth.GameServiceTokenFromEnv() == "" {
		return nil
	}
	return fmt.Errorf("game service authentication required")
}
