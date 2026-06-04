package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func resendEmailsURL() string {
	base := strings.TrimSpace(os.Getenv("RESEND_API_URL"))
	if base == "" {
		base = "https://api.resend.com"
	}
	return strings.TrimRight(base, "/") + "/emails"
}

// ResendConfig holds Resend HTTP API settings.
type ResendConfig struct {
	APIKey   string
	From     string
	FromName string
}

// ResendSender delivers email via the Resend REST API (preferred over SMTP).
type ResendSender struct {
	config     ResendConfig
	httpClient *http.Client
}

// NewResendSender creates a Resend API-backed sender.
func NewResendSender(config ResendConfig) (*ResendSender, error) {
	if strings.TrimSpace(config.APIKey) == "" {
		return nil, fmt.Errorf("email: Resend API key is required")
	}
	if strings.TrimSpace(config.From) == "" {
		return nil, fmt.Errorf("email: from address is required")
	}
	return &ResendSender{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (s *ResendSender) SendMagicLink(ctx context.Context, msg MagicLinkEmail) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(msg.To) == "" {
		return fmt.Errorf("email: recipient is required")
	}
	if strings.TrimSpace(msg.Link) == "" {
		return fmt.Errorf("email: magic link URL is required")
	}

	from := formatAddress(s.config.From, s.config.FromName)
	payload := resendEmailRequest{
		From:    from,
		To:      []string{msg.To},
		Subject: magicLinkSubject(),
		Text:    magicLinkBody(msg),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("email: marshal Resend payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEmailsURL(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("email: create Resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("email: Resend API request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
	if err != nil {
		return fmt.Errorf("email: read Resend response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		detail := strings.TrimSpace(string(respBody))
		if detail == "" {
			detail = resp.Status
		}
		return fmt.Errorf("email: Resend API %s: %s", resp.Status, detail)
	}
	return nil
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}
