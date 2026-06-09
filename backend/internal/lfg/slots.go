package lfg

// SlotsFromAssignments converts persisted assignments to mutable slots.
func SlotsFromAssignments(assignments []Assignment) []SeatSlot {
	out := make([]SeatSlot, len(assignments))
	for i, a := range assignments {
		out[i] = SeatSlot{
			SeatKey:     a.SeatKey,
			QueuePath:   a.QueuePath,
			AffinityKey: a.AffinityKey,
			OccupantID:  a.UserID,
		}
	}
	return out
}

// SeatSlot is one assignable cell on a forming match map.
type SeatSlot struct {
	SeatKey     string
	QueuePath   string
	AffinityKey string
	OccupantID  string
}
