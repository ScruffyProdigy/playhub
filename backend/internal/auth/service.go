package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/store"
)

var (
	ErrInvalidEmail    = errors.New("auth: invalid email address")
	ErrInvalidMagicLink = errors.New("auth: invalid or expired magic link")
)

// Service handles authentication workflows.
type Service struct {
	store            *store.Store
	signer           *Signer
	cookie           CookieConfig
	magicLinkTTL     time.Duration
	sessionTTL       time.Duration
	magicLinkBaseURL string
}

// NewService creates an authentication service backed by the store and JWT signer.
func NewService(st *store.Store, signer *Signer) (*Service, error) {
	if st == nil {
		return nil, errors.New("auth: store is required")
	}
	if signer == nil {
		return nil, errors.New("auth: signer is required")
	}

	return &Service{
		store:            st,
		signer:           signer,
		cookie:           CookieConfigFromEnv(),
		magicLinkTTL:     durationFromEnv("MAGIC_LINK_TTL", 15*time.Minute),
		sessionTTL:       durationFromEnv("SESSION_TTL", 7*24*time.Hour),
		magicLinkBaseURL: strings.TrimSpace(os.Getenv("MAGIC_LINK_BASE_URL")),
	}, nil
}

// Signer returns the JWT signer used by the service.
func (s *Service) Signer() *Signer {
	return s.signer
}

// CookieConfig returns the session cookie configuration.
func (s *Service) CookieConfig() CookieConfig {
	return s.cookie
}

// RequestMagicLink creates a magic link for the given email address.
func (s *Service) RequestMagicLink(ctx context.Context, email string) error {
	normalized, err := normalizeEmail(email)
	if err != nil {
		return err
	}

	token := uuid.NewString()
	_, err = s.store.CreateMagicLink(ctx, store.CreateMagicLinkParams{
		Email:     normalized,
		Token:     token,
		ExpiresAt: time.Now().Add(s.magicLinkTTL),
	})
	if err != nil {
		return err
	}

	s.logMagicLink(normalized, token)
	return nil
}

// CompleteMagicLogin validates a magic link token and returns the user plus session JWT.
func (s *Service) CompleteMagicLogin(ctx context.Context, token string) (*store.User, string, error) {
	link, err := s.store.GetMagicLinkByToken(ctx, token)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, "", ErrInvalidMagicLink
		}
		return nil, "", err
	}

	if !s.store.IsMagicLinkValid(link, time.Now()) {
		return nil, "", ErrInvalidMagicLink
	}

	user, err := s.findOrCreateUser(ctx, link.Email)
	if err != nil {
		return nil, "", err
	}

	if err := s.store.MarkMagicLinkUsed(ctx, link.ID); err != nil {
		return nil, "", err
	}

	if err := s.store.UpdateUserLastLogin(ctx, user.ID); err != nil {
		return nil, "", err
	}

	sessionToken, err := s.signer.SignUserToken(user.ID, s.sessionTTL)
	if err != nil {
		return nil, "", err
	}

	return user, sessionToken, nil
}

// GetAuthenticatedUser loads the current user from context, if any.
func (s *Service) GetAuthenticatedUser(ctx context.Context) (*store.User, error) {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		return nil, nil
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("auth: invalid user id in context: %w", err)
	}

	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) findOrCreateUser(ctx context.Context, email string) (*store.User, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}

	displayName := strings.Split(email, "@")[0]
	return s.store.CreateUser(ctx, store.CreateUserParams{
		Email:       email,
		DisplayName: displayName,
	})
}

func (s *Service) logMagicLink(email, token string) {
	link := s.magicLinkURL(token)
	if link == "" {
		log.Printf("auth: magic link for %s token=%s", email, token)
		return
	}
	log.Printf("auth: magic link for %s -> %s", email, link)
}

func (s *Service) magicLinkURL(token string) string {
	if s.magicLinkBaseURL == "" {
		return ""
	}
	if strings.Contains(s.magicLinkBaseURL, "{token}") {
		return strings.ReplaceAll(s.magicLinkBaseURL, "{token}", token)
	}
	if strings.Contains(s.magicLinkBaseURL, "%s") {
		return fmt.Sprintf(s.magicLinkBaseURL, token)
	}
	if strings.Contains(s.magicLinkBaseURL, "?") {
		return s.magicLinkBaseURL + token
	}
	return strings.TrimRight(s.magicLinkBaseURL, "/") + "?token=" + token
}

func normalizeEmail(email string) (string, error) {
	address, err := mail.ParseAddress(strings.TrimSpace(email))
	if err != nil {
		return "", ErrInvalidEmail
	}
	return strings.ToLower(address.Address), nil
}

func durationFromEnv(name string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
