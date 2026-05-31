package auth

import "testing"

func TestIsAdminEmail(t *testing.T) {
	t.Setenv("LOBBY_ADMIN_EMAILS", "ryan.c.kohler@gmail.com, Admin@Example.com ")

	if !IsAdminEmail("ryan.c.kohler@gmail.com") {
		t.Fatal("expected allowlisted email to be admin")
	}
	if !IsAdminEmail("Admin@Example.com") {
		t.Fatal("expected admin check to be case-insensitive")
	}
	if IsAdminEmail("other@example.com") {
		t.Fatal("expected non-allowlisted email to be rejected")
	}
	if IsAdminEmail("not-an-email") {
		t.Fatal("expected invalid email to be rejected")
	}
}
