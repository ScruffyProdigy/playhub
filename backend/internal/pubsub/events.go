package pubsub

import "encoding/json"

// QueueEventStatus is delivered to queue subscribers.
type QueueEventStatus string

const (
	QueueStatusWaiting QueueEventStatus = "WAITING"
	QueueStatusMatched QueueEventStatus = "MATCHED"
	QueueStatusLeft    QueueEventStatus = "LEFT"
)

// QueueEvent is published on per-user Redis channels when queue state changes.
type QueueEvent struct {
	GameID      string           `json:"gameId"`
	Status      QueueEventStatus `json:"status"`
	SessionID   string           `json:"sessionId,omitempty"`
	JoinURL     string           `json:"joinUrl,omitempty"`
	QueuedCount int              `json:"queuedCount,omitempty"`
}

// UserQueueChannel returns the Redis channel for a user's queue updates.
func UserQueueChannel(userID string) string {
	return "lobby:user:" + userID + ":queue"
}

func MarshalQueueEvent(event QueueEvent) ([]byte, error) {
	return json.Marshal(event)
}

func UnmarshalQueueEvent(payload []byte) (QueueEvent, error) {
	var event QueueEvent
	err := json.Unmarshal(payload, &event)
	return event, err
}
