package email

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewSMTPSenderValidation(t *testing.T) {
	t.Run("requires host", func(t *testing.T) {
		_, err := NewSMTPSender(Config{From: "noreply@example.com"})
		if err == nil || !strings.Contains(err.Error(), "host") {
			t.Fatalf("expected host error, got %v", err)
		}
	})

	t.Run("requires from address", func(t *testing.T) {
		_, err := NewSMTPSender(Config{Host: "smtp.example.com"})
		if err == nil || !strings.Contains(err.Error(), "from") {
			t.Fatalf("expected from error, got %v", err)
		}
	})

	t.Run("defaults port to 587", func(t *testing.T) {
		sender, err := NewSMTPSender(Config{
			Host: "smtp.example.com",
			From: "noreply@example.com",
		})
		if err != nil {
			t.Fatalf("NewSMTPSender() error: %v", err)
		}
		if sender.config.Port != 587 {
			t.Fatalf("expected default port 587, got %d", sender.config.Port)
		}
	})
}

func TestSendMagicLinkValidation(t *testing.T) {
	sender, err := NewSMTPSender(Config{
		Host: "smtp.example.com",
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("NewSMTPSender() error: %v", err)
	}

	t.Run("requires recipient", func(t *testing.T) {
		err := sender.SendMagicLink(context.Background(), MagicLinkEmail{
			Link: "http://localhost/auth/complete?token=abc",
			TTL:  15 * time.Minute,
		})
		if err == nil || !strings.Contains(err.Error(), "recipient") {
			t.Fatalf("expected recipient error, got %v", err)
		}
	})

	t.Run("requires link", func(t *testing.T) {
		err := sender.SendMagicLink(context.Background(), MagicLinkEmail{
			To:  "player@example.com",
			TTL: 15 * time.Minute,
		})
		if err == nil || !strings.Contains(err.Error(), "link") {
			t.Fatalf("expected link error, got %v", err)
		}
	})

	t.Run("respects cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := sender.SendMagicLink(ctx, MagicLinkEmail{
			To:   "player@example.com",
			Link: "http://localhost/auth/complete?token=abc",
			TTL:  15 * time.Minute,
		})
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestSendMagicLinkReturnsConnectError(t *testing.T) {
	sender, err := NewSMTPSender(Config{
		Host: "127.0.0.1",
		Port: 1,
		From: "noreply@example.com",
	})
	if err != nil {
		t.Fatalf("NewSMTPSender() error: %v", err)
	}

	err = sender.SendMagicLink(context.Background(), MagicLinkEmail{
		To:   "player@example.com",
		Link: "http://localhost/auth/complete?token=abc",
		TTL:  15 * time.Minute,
	})
	if err == nil || !strings.Contains(err.Error(), "SMTP") {
		t.Fatalf("expected SMTP connect error, got %v", err)
	}
}

func TestBuildMessage(t *testing.T) {
	message := string(buildMessage(
		"JoinQuest <noreply@example.com>",
		"player@example.com",
		"Sign in to JoinQuest",
		"Click here",
	))

	for _, want := range []string{
		"From: JoinQuest <noreply@example.com>",
		"To: player@example.com",
		"Subject: Sign in to JoinQuest",
		"Content-Type: text/plain; charset=UTF-8",
		"Click here",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected message to contain %q, got %q", want, message)
		}
	}
}

func TestFormatAddress(t *testing.T) {
	if got := formatAddress("noreply@example.com", ""); got != "noreply@example.com" {
		t.Fatalf("formatAddress without name = %q", got)
	}

	want := "JoinQuest <noreply@example.com>"
	if got := formatAddress("noreply@example.com", "JoinQuest"); got != want {
		t.Fatalf("formatAddress with name = %q, want %q", got, want)
	}
}

func TestSenderFromEnvUsesSMTPWhenConfigured(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "2525")
	t.Setenv("SMTP_FROM", "noreply@example.com")
	t.Setenv("SMTP_FROM_NAME", "JoinQuest")

	sender, err := SenderFromEnv()
	if err != nil {
		t.Fatalf("SenderFromEnv() error: %v", err)
	}
	if _, ok := sender.(*SMTPSender); !ok {
		t.Fatalf("expected SMTPSender, got %T", sender)
	}
}

func TestConfigFromEnvInvalidPort(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "not-a-port")
	t.Setenv("SMTP_FROM", "noreply@example.com")

	_, ok, err := ConfigFromEnv()
	if err == nil {
		t.Fatal("expected invalid port error")
	}
	if ok {
		t.Fatal("expected ok=false")
	}
}
