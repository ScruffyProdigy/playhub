package graph

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/auth"
)

func requireNonGuestAccount(ctx context.Context, authService *auth.Service) (*auth.Service, uuid.UUID, error) {
	user, err := authService.RequireNonGuestUser(ctx)
	if err != nil {
		return nil, uuid.Nil, err
	}
	return authService, user.ID, nil
}

func requireAuthenticatedUserID(ctx context.Context) (uuid.UUID, error) {
	userID, err := requireAuthUserID(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("authentication required")
	}
	return userID, nil
}

func mapAuthProvider(provider string) (string, error) {
	switch provider {
	case "google":
		return "GOOGLE", nil
	case "discord":
		return "DISCORD", nil
	case "apple":
		return "APPLE", nil
	case "facebook":
		return "FACEBOOK", nil
	default:
		return "", fmt.Errorf("unknown auth provider %q", provider)
	}
}
