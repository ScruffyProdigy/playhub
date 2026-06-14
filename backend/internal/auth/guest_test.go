package auth

import (
	"context"
	"testing"
)

func TestCreateGuestSession(t *testing.T) {
	service := openAuthTestService(t)
	ctx := context.Background()

	user, token, err := service.CreateGuestSession(ctx)
	if err != nil {
		t.Fatalf("CreateGuestSession: %v", err)
	}
	if token == "" {
		t.Fatal("expected session token")
	}
	if !user.IsGuest {
		t.Fatal("expected guest user")
	}
	if user.Email != "" {
		t.Fatalf("expected no email, got %q", user.Email)
	}
}

func TestRequireNonGuestUserRejectsGuest(t *testing.T) {
	service := openAuthTestService(t)
	ctx := context.Background()

	guest, _, err := service.CreateGuestSession(ctx)
	if err != nil {
		t.Fatalf("CreateGuestSession: %v", err)
	}

	ctx = WithUserID(ctx, guest.ID.String())
	_, err = service.RequireNonGuestUser(ctx)
	if err == nil {
		t.Fatal("expected guest rejection")
	}
}
