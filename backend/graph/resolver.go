package graph

import (
	"fmt"

	"github.com/scruffyprodigy/playhub/internal/auth"
	"github.com/scruffyprodigy/playhub/internal/store"
)

type Resolver struct {
	Store *store.Store
	Auth  *auth.Service
}

// NewResolver creates a resolver backed by the store and auth service.
func NewResolver(st *store.Store, authService *auth.Service) *Resolver {
	return &Resolver{Store: st, Auth: authService}
}

func (r *Resolver) requireStore() (*store.Store, error) {
	if r == nil || r.Store == nil {
		return nil, fmt.Errorf("database store is not configured")
	}
	return r.Store, nil
}

func (r *Resolver) requireAuth() (*auth.Service, error) {
	if r == nil || r.Auth == nil {
		return nil, fmt.Errorf("auth service is not configured")
	}
	return r.Auth, nil
}
