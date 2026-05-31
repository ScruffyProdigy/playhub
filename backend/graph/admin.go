package graph

import (
	"context"
	"fmt"

	"github.com/scruffyprodigy/playhub/internal/auth"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func (r *Resolver) requireAdmin(ctx context.Context) (*store.User, error) {
	authService, err := r.requireAuth()
	if err != nil {
		return nil, err
	}
	user, err := authService.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("authentication required")
	}
	if !auth.IsAdminEmail(user.Email) {
		return nil, auth.ErrNotAdmin
	}
	return user, nil
}
