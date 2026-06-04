package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const smtpDialTimeout = 15 * time.Second

// SMTPSender delivers email over SMTP.
type SMTPSender struct {
	config Config
}

// NewSMTPSender creates an SMTP-backed sender.
func NewSMTPSender(config Config) (*SMTPSender, error) {
	if strings.TrimSpace(config.Host) == "" {
		return nil, fmt.Errorf("email: SMTP host is required")
	}
	if config.Port == 0 {
		config.Port = 587
	}
	if strings.TrimSpace(config.From) == "" {
		return nil, fmt.Errorf("email: SMTP from address is required")
	}
	return &SMTPSender{config: config}, nil
}

func (s *SMTPSender) SendMagicLink(ctx context.Context, msg MagicLinkEmail) error {
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
	body := buildMessage(from, msg.To, magicLinkSubject(), magicLinkBody(msg))
	return s.send(ctx, msg.To, body)
}

func (s *SMTPSender) send(ctx context.Context, to string, message []byte) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	auth := smtpAuth(s.config)

	if s.config.Port == 465 {
		return sendMailTLS(ctx, addr, s.config.Host, auth, s.config.From, []string{to}, message)
	}
	return sendMailSTARTTLS(ctx, addr, s.config.Host, auth, s.config.From, []string{to}, message)
}

func smtpAuth(config Config) smtp.Auth {
	if config.Username == "" && config.Password == "" {
		return nil
	}
	return smtp.PlainAuth("", config.Username, config.Password, config.Host)
}

func formatAddress(email, name string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

func buildMessage(from, to, subject, body string) []byte {
	var msg strings.Builder
	msg.WriteString("From: " + from + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)
	return []byte(msg.String())
}

func dialSMTP(ctx context.Context, addr, host string) (*smtp.Client, error) {
	dialer := &net.Dialer{Timeout: smtpDialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("email: connect to SMTP server: %w", err)
	}
	return smtp.NewClient(conn, host)
}

func sendMailSTARTTLS(ctx context.Context, addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	client, err := dialSMTP(ctx, addr, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Hello("localhost"); err != nil {
		return fmt.Errorf("email: SMTP hello: %w", err)
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: host}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("email: STARTTLS: %w", err)
		}
	}

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); !ok {
			return fmt.Errorf("email: SMTP server does not support AUTH")
		}
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("email: SMTP auth: %w", err)
		}
	}

	if err := sendMailMessage(client, from, to, msg); err != nil {
		return err
	}
	return client.Quit()
}

func sendMailTLS(ctx context.Context, addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	dialer := &net.Dialer{Timeout: smtpDialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("email: connect to SMTP server: %w", err)
	}

	tlsConn := tls.Client(conn, &tls.Config{ServerName: host})
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		_ = conn.Close()
		return fmt.Errorf("email: TLS handshake: %w", err)
	}

	client, err := smtp.NewClient(tlsConn, host)
	if err != nil {
		_ = tlsConn.Close()
		return fmt.Errorf("email: create SMTP client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("email: SMTP auth: %w", err)
		}
	}

	if err := sendMailMessage(client, from, to, msg); err != nil {
		return err
	}
	return client.Quit()
}

func sendMailMessage(client *smtp.Client, from string, to []string, msg []byte) error {
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("email: MAIL FROM: %w", err)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("email: RCPT TO %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: DATA: %w", err)
	}

	if _, err := writer.Write(msg); err != nil {
		_ = writer.Close()
		return fmt.Errorf("email: write message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("email: close message: %w", err)
	}
	return nil
}
