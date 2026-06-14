package store

import (
	"context"
	"regexp"
	"testing"
)

func TestRandomGuestDisplayName(t *testing.T) {
	name, err := RandomGuestDisplayName()
	if err != nil {
		t.Fatalf("RandomGuestDisplayName: %v", err)
	}
	if !regexp.MustCompile(`^guest#[0-9]{6}$`).MatchString(name) {
		t.Fatalf("expected guest#NNNNNN format, got %q", name)
	}
}

func TestCreateGuestUser(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	user, err := st.CreateGuestUser(ctx)
	if err != nil {
		t.Fatalf("CreateGuestUser: %v", err)
	}
	cleaner.TrackUser(user.ID)

	if !user.IsGuest {
		t.Fatal("expected guest user")
	}
	if user.Email != "" {
		t.Fatalf("expected empty email, got %q", user.Email)
	}
	if user.DisplayName == "" {
		t.Fatal("expected display name")
	}
	if !regexp.MustCompile(`^guest#[0-9]{6}$`).MatchString(user.DisplayName) {
		t.Fatalf("expected guest display name, got %q", user.DisplayName)
	}
}

func TestAddVerifiedUserEmailClearsGuest(t *testing.T) {
	st := openTestStore(t)
	cleaner := st.NewTestCleaner(t)
	ctx := context.Background()

	guest, err := st.CreateGuestUser(ctx)
	if err != nil {
		t.Fatalf("CreateGuestUser: %v", err)
	}
	cleaner.TrackUser(guest.ID)

	if _, err := st.AddVerifiedUserEmail(ctx, guest.ID, "linked-"+guest.ID.String()+"@example.com", true); err != nil {
		t.Fatalf("AddVerifiedUserEmail: %v", err)
	}

	user, err := st.GetUserByID(ctx, guest.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if user.IsGuest {
		t.Fatal("expected is_guest false after linking email")
	}
	if user.Email == "" {
		t.Fatal("expected primary email on user row")
	}
}
