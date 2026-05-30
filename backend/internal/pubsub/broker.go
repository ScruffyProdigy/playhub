package pubsub

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// Broker publishes and subscribes to channels for cross-instance fan-out.
type Broker interface {
	Publish(ctx context.Context, channel string, payload []byte) error
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
	Close() error
}

// NewFromEnv returns a Redis broker when REDIS_URL is set, otherwise an in-memory broker.
func NewFromEnv() (Broker, error) {
	url := strings.TrimSpace(os.Getenv("REDIS_URL"))
	if url == "" {
		return NewMemory(), nil
	}
	return NewRedis(url)
}

// PublishQueueEvent marshals and publishes a queue event to the given user channel.
func PublishQueueEvent(ctx context.Context, broker Broker, userID string, event QueueEvent) error {
	if broker == nil {
		return fmt.Errorf("pubsub: broker is required")
	}
	payload, err := MarshalQueueEvent(event)
	if err != nil {
		return err
	}
	return broker.Publish(ctx, UserQueueChannel(userID), payload)
}
