package store

import (
	"fmt"
	"strings"
)

type seatAssignment struct {
	Entry   QueueEntry
	SeatKey string
}

func seatQueuePathValue(seat GameModeSeat) string {
	if seat.QueuePath == nil {
		return ""
	}
	return strings.TrimSpace(*seat.QueuePath)
}

func entryQueuePathValue(entry QueueEntry) string {
	if entry.QueuePath == nil {
		return ""
	}
	return strings.TrimSpace(*entry.QueuePath)
}

func modeUsesComposition(seats []GameModeSeat) bool {
	for _, seat := range seats {
		if seatQueuePathValue(seat) != "" {
			return true
		}
	}
	return false
}

func validQueuePaths(seats []GameModeSeat) map[string]struct{} {
	paths := make(map[string]struct{})
	for _, seat := range seats {
		if path := seatQueuePathValue(seat); path != "" {
			paths[path] = struct{}{}
		}
	}
	return paths
}

func validateJoinQueuePath(seats []GameModeSeat, queuePath string) error {
	queuePath = strings.TrimSpace(queuePath)
	if !modeUsesComposition(seats) {
		if queuePath != "" {
			return fmt.Errorf("store: queue path is not used for this mode")
		}
		return nil
	}
	if queuePath == "" {
		return fmt.Errorf("store: queue path is required for this mode")
	}
	if _, ok := validQueuePaths(seats)[queuePath]; !ok {
		return fmt.Errorf("store: invalid queue path %q", queuePath)
	}
	return nil
}

// tryFormMatch pairs waiting players to mode seats by queue_path (fifo within each path).
func tryFormMatch(seats []GameModeSeat, waiting []QueueEntry) ([]seatAssignment, bool) {
	if len(seats) == 0 {
		return nil, false
	}

	byPath := make(map[string][]QueueEntry)
	for _, entry := range waiting {
		path := entryQueuePathValue(entry)
		byPath[path] = append(byPath[path], entry)
	}

	used := make(map[string]struct{}, len(waiting))
	cursor := make(map[string]int, len(byPath))
	for path := range byPath {
		cursor[path] = 0
	}

	out := make([]seatAssignment, 0, len(seats))
	for _, seat := range seats {
		path := seatQueuePathValue(seat)
		bucket := byPath[path]
		idx := cursor[path]
		var picked *QueueEntry
		for idx < len(bucket) {
			candidate := bucket[idx]
			idx++
			key := candidate.UserID.String()
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			picked = &candidate
			break
		}
		cursor[path] = idx
		if picked == nil {
			return nil, false
		}
		out = append(out, seatAssignment{Entry: *picked, SeatKey: seat.SeatKey})
	}
	return out, true
}

func nullQueuePathColumn(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func sameQueuePath(existing *string, requested string) bool {
	return entryQueuePathValue(QueueEntry{QueuePath: existing}) == strings.TrimSpace(requested)
}
