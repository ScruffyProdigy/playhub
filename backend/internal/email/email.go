package email

import (
	"context"
	"fmt"
	"time"
)

// MagicLinkEmail is the payload for a sign-in email.
type MagicLinkEmail struct {
	To   string
	Link string
	Code string
	TTL  time.Duration
}

// ProductName is used in sign-in email subject and body.
const ProductName = "JoinQuest"

// Sender delivers transactional email.
type Sender interface {
	SendMagicLink(ctx context.Context, msg MagicLinkEmail) error
}

func formatTTL(ttl time.Duration) string {
	if ttl < time.Minute {
		return "less than a minute"
	}
	minutes := int(ttl.Round(time.Minute).Minutes())
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}

func magicLinkSubject() string {
	return fmt.Sprintf("Sign in to %s", ProductName)
}

func magicLinkBody(msg MagicLinkEmail) string {
	if msg.Code != "" {
		return fmt.Sprintf(`Sign in to %s

Your sign-in code: %s

Or click this link to continue:
%s

This code and link expire in %s. If you did not request this email, you can ignore it.
`, ProductName, msg.Code, msg.Link, formatTTL(msg.TTL))
	}

	return fmt.Sprintf(`Sign in to %s

Click this link to continue:
%s

This link expires in %s. If you did not request this email, you can ignore it.
`, ProductName, msg.Link, formatTTL(msg.TTL))
}
