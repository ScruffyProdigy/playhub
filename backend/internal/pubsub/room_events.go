package pubsub

import (
	"context"
	"encoding/json"
)

const (
	RoomEventUpdated = "updated"
	RoomEventMessage = "message"
)

// RoomEvent is published on a room channel for membership or chat updates.
type RoomEvent struct {
	Type      string `json:"type"`
	RoomID    string `json:"roomId"`
	MessageID string `json:"messageId,omitempty"`
}

// RoomChannel returns the Redis channel for room updates.
func RoomChannel(roomID string) string {
	return "lobby:room:" + roomID
}

func MarshalRoomEvent(event RoomEvent) ([]byte, error) {
	return json.Marshal(event)
}

func UnmarshalRoomEvent(payload []byte) (RoomEvent, error) {
	var event RoomEvent
	err := json.Unmarshal(payload, &event)
	return event, err
}

// PublishRoomEvent marshals and publishes a room event.
func PublishRoomEvent(ctx context.Context, broker Broker, roomID string, event RoomEvent) error {
	if broker == nil {
		return nil
	}
	event.RoomID = roomID
	payload, err := MarshalRoomEvent(event)
	if err != nil {
		return err
	}
	return broker.Publish(ctx, RoomChannel(roomID), payload)
}
