package pubsub

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Redis implements Broker using Redis pub/sub.
type Redis struct {
	client *redis.Client
}

// NewRedis connects to Redis at the given URL (redis://host:port/db).
func NewRedis(url string) (*Redis, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("pubsub: parse REDIS_URL: %w", err)
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("pubsub: ping redis: %w", err)
	}
	return &Redis{client: client}, nil
}

func (r *Redis) Publish(ctx context.Context, channel string, payload []byte) error {
	return r.client.Publish(ctx, channel, payload).Err()
}

func (r *Redis) Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error) {
	pubsub := r.client.Subscribe(ctx, channel)
	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()
		return nil, nil, fmt.Errorf("pubsub: subscribe %s: %w", channel, err)
	}

	out := make(chan []byte, 16)
	done := make(chan struct{})

	go func() {
		defer close(out)
		ch := pubsub.Channel()
		for {
			select {
			case <-done:
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				out <- []byte(msg.Payload)
			}
		}
	}()

	unsubscribe := func() {
		close(done)
		_ = pubsub.Close()
	}

	return out, unsubscribe, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
