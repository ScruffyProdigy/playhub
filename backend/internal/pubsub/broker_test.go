package pubsub

import (
	"context"
	"testing"
	"time"
)

func TestMemoryBrokerPublishSubscribe(t *testing.T) {
	broker := NewMemory()
	defer broker.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, unsubscribe, err := broker.Subscribe(ctx, UserQueueChannel("user-1"))
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	defer unsubscribe()

	event := QueueEvent{
		GameID:      "game-1",
		Status:      QueueStatusMatched,
		SessionID:   "session-1",
		JoinURL:     "http://localhost:5174/join/session-1",
		QueuedCount: 0,
	}
	if err := PublishQueueEvent(ctx, broker, "user-1", event); err != nil {
		t.Fatalf("PublishQueueEvent failed: %v", err)
	}

	select {
	case payload := <-ch:
		got, err := UnmarshalQueueEvent(payload)
		if err != nil {
			t.Fatalf("UnmarshalQueueEvent failed: %v", err)
		}
		if got.Status != QueueStatusMatched || got.SessionID != "session-1" {
			t.Fatalf("unexpected event: %+v", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for queue event")
	}
}
