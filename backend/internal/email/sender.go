package email

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// LogSender writes magic links to the application log (default for local dev).
type LogSender struct{}

func (LogSender) SendMagicLink(_ context.Context, msg MagicLinkEmail) error {
	if msg.Link == "" {
		log.Printf("email: magic link for %s (URL not configured)", msg.To)
		return nil
	}
	log.Printf("email: magic link for %s -> %s", msg.To, msg.Link)
	return nil
}

// Config holds SMTP connection settings.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// ConfigFromEnv loads SMTP settings. Returns ok=false when SMTP is not configured.
func ConfigFromEnv() (Config, bool, error) {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	if host == "" {
		return Config{}, false, nil
	}

	port := 587
	if rawPort := strings.TrimSpace(os.Getenv("SMTP_PORT")); rawPort != "" {
		parsed, err := strconv.Atoi(rawPort)
		if err != nil {
			return Config{}, false, fmt.Errorf("email: invalid SMTP_PORT %q: %w", rawPort, err)
		}
		port = parsed
	}

	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
	if from == "" {
		return Config{}, false, fmt.Errorf("email: SMTP_FROM is required when SMTP_HOST is set")
	}

	return Config{
		Host:     host,
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     from,
		FromName: strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
	}, true, nil
}

// SenderFromEnv returns an SMTP sender when configured, otherwise LogSender.
func SenderFromEnv() (Sender, error) {
	cfg, ok, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}
	if !ok {
		return LogSender{}, nil
	}
	return NewSMTPSender(cfg)
}
