package auth

import "testing"

func TestVerifyLoginCode(t *testing.T) {
	code := "042815"
	hash := hashLoginCode(code)

	if !verifyLoginCode(code, hash) {
		t.Fatal("expected valid code to verify")
	}
	if verifyLoginCode("042816", hash) {
		t.Fatal("expected invalid code to fail")
	}
	if verifyLoginCode("12345", hash) {
		t.Fatal("expected short code to fail")
	}
}

func TestGenerateLoginCode(t *testing.T) {
	code, err := generateLoginCode()
	if err != nil {
		t.Fatalf("generateLoginCode() error: %v", err)
	}
	if !loginCodePattern.MatchString(code) {
		t.Fatalf("expected 6-digit code, got %q", code)
	}
}
