package email

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestResendSenderSendMagicLink(t *testing.T) {
	var got resendEmailRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method: got %s", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer re_test_key" {
			t.Fatalf("authorization: got %q", auth)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"id":"email_123"}`)
	}))
	defer server.Close()

	t.Setenv("RESEND_API_URL", server.URL)

	sender, err := NewResendSender(ResendConfig{
		APIKey:   "re_test_key",
		From:     "noreply@example.com",
		FromName: "JoinQuest",
	})
	if err != nil {
		t.Fatalf("NewResendSender() error: %v", err)
	}

	err = sender.SendMagicLink(context.Background(), MagicLinkEmail{
		To:   "player@example.com",
		Link: "https://example.com/auth?token=abc",
		Code: "123456",
		TTL:  15 * time.Minute,
	})
	if err != nil {
		t.Fatalf("SendMagicLink() error: %v", err)
	}
	if got.From != "JoinQuest <noreply@example.com>" {
		t.Fatalf("from: got %q", got.From)
	}
	if len(got.To) != 1 || got.To[0] != "player@example.com" {
		t.Fatalf("to: got %v", got.To)
	}
	if !strings.Contains(got.Text, "123456") {
		t.Fatalf("body missing code: %q", got.Text)
	}
}

func TestResendSenderRejectsMissingRecipient(t *testing.T) {
	sender, err := NewResendSender(ResendConfig{APIKey: "re_x", From: "noreply@example.com"})
	if err != nil {
		t.Fatalf("NewResendSender() error: %v", err)
	}
	err = sender.SendMagicLink(context.Background(), MagicLinkEmail{Link: "https://example.com/auth"})
	if err == nil || !strings.Contains(err.Error(), "recipient") {
		t.Fatalf("expected recipient error, got %v", err)
	}
}

func TestResendSenderReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = io.WriteString(w, `{"message":"Domain not verified"}`)
	}))
	defer server.Close()

	t.Setenv("RESEND_API_URL", server.URL)

	sender, err := NewResendSender(ResendConfig{APIKey: "re_x", From: "noreply@example.com"})
	if err != nil {
		t.Fatalf("NewResendSender() error: %v", err)
	}

	err = sender.SendMagicLink(context.Background(), MagicLinkEmail{
		To:   "user@example.com",
		Link: "https://example.com/auth?token=abc",
		Code: "123456",
	})
	if err == nil || !strings.Contains(err.Error(), "Domain not verified") {
		t.Fatalf("expected API error, got %v", err)
	}
}
