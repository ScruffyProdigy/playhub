package pubsub

import (
	"context"
	"encoding/json"
)

const (
	TableSeatStatusStarted = "STARTED"
	TableSeatStatusLeft    = "LEFT"
)

// TableSeatEvent is published on per-user channels when a private table seat changes.
type TableSeatEvent struct {
	Status  string `json:"status"`
	TableID string `json:"tableId,omitempty"`
	JoinURL string `json:"joinUrl,omitempty"`
}

// UserTableSeatChannel returns the channel for a user's table seat updates.
func UserTableSeatChannel(userID string) string {
	return "lobby:user:" + userID + ":table-seat"
}

func MarshalTableSeatEvent(event TableSeatEvent) ([]byte, error) {
	return json.Marshal(event)
}

func UnmarshalTableSeatEvent(payload []byte) (TableSeatEvent, error) {
	var event TableSeatEvent
	err := json.Unmarshal(payload, &event)
	return event, err
}

// PublishTableSeatEvent marshals and publishes a table seat event to the user channel.
func PublishTableSeatEvent(ctx context.Context, broker Broker, userID string, event TableSeatEvent) error {
	if broker == nil {
		return nil
	}
	payload, err := MarshalTableSeatEvent(event)
	if err != nil {
		return err
	}
	return broker.Publish(ctx, UserTableSeatChannel(userID), payload)
}
