package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const oauthCodeKeyPrefix = "oauth:code:"

// OAuthCodeStore deduplicates concurrent OAuth callbacks (common on mobile Safari).
type OAuthCodeStore interface {
	TryClaim(ctx context.Context, provider, code string) (claimed bool, err error)
	SetCompleted(ctx context.Context, provider, code, userID string) error
	GetCompleted(ctx context.Context, provider, code string) (userID string, ok bool, err error)
}

type memoryOAuthCodeStore struct {
	mu    sync.Mutex
	items map[string]string
}

type redisOAuthCodeStore struct {
	client *redis.Client
}

func NewOAuthCodeStoreFromEnv() (OAuthCodeStore, error) {
	url := strings.TrimSpace(os.Getenv("REDIS_URL"))
	if url == "" {
		return newMemoryOAuthCodeStore(), nil
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("auth: parse REDIS_URL for oauth code store: %w", err)
	}
	opts.DialTimeout = 2 * time.Second
	opts.ReadTimeout = 2 * time.Second
	opts.WriteTimeout = 2 * time.Second
	opts.PoolTimeout = 3 * time.Second
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("auth: ping redis for oauth code store: %w", err)
	}
	return &redisOAuthCodeStore{client: client}, nil
}

func newMemoryOAuthCodeStore() *memoryOAuthCodeStore {
	return &memoryOAuthCodeStore{items: make(map[string]string)}
}

func oauthCodeKey(provider, code string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(provider)) + ":" + strings.TrimSpace(code)))
	return oauthCodeKeyPrefix + hex.EncodeToString(sum[:])
}

func (s *memoryOAuthCodeStore) TryClaim(_ context.Context, provider, code string) (bool, error) {
	key := oauthCodeKey(provider, code)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[key]; ok {
		return false, nil
	}
	s.items[key] = "processing"
	time.AfterFunc(2*time.Minute, func() {
		s.mu.Lock()
		delete(s.items, key)
		s.mu.Unlock()
	})
	return true, nil
}

func (s *memoryOAuthCodeStore) SetCompleted(_ context.Context, provider, code, userID string) error {
	key := oauthCodeKey(provider, code)
	s.mu.Lock()
	s.items[key] = "done:" + strings.TrimSpace(userID)
	s.mu.Unlock()
	return nil
}

func (s *memoryOAuthCodeStore) GetCompleted(_ context.Context, provider, code string) (string, bool, error) {
	key := oauthCodeKey(provider, code)
	s.mu.Lock()
	val, ok := s.items[key]
	s.mu.Unlock()
	if !ok || !strings.HasPrefix(val, "done:") {
		return "", false, nil
	}
	return strings.TrimPrefix(val, "done:"), true, nil
}

func (s *redisOAuthCodeStore) TryClaim(ctx context.Context, provider, code string) (bool, error) {
	key := oauthCodeKey(provider, code)
	ok, err := s.client.SetNX(ctx, key, "processing", 2*time.Minute).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *redisOAuthCodeStore) SetCompleted(ctx context.Context, provider, code, userID string) error {
	key := oauthCodeKey(provider, code)
	return s.client.Set(ctx, key, "done:"+strings.TrimSpace(userID), 2*time.Minute).Err()
}

func (s *redisOAuthCodeStore) GetCompleted(ctx context.Context, provider, code string) (string, bool, error) {
	key := oauthCodeKey(provider, code)
	val, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if !strings.HasPrefix(val, "done:") {
		return "", false, nil
	}
	return strings.TrimPrefix(val, "done:"), true, nil
}
