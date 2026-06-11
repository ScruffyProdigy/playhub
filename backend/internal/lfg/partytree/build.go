package partytree

import (
	"strings"
)

// BuildFromPinnedSeats constructs a party tree from table seat assignments.
// seatRoles maps seatKey -> queuePath (leaf role bucket) from game_mode_seats.
func BuildFromPinnedSeats(seats []PinnedSeat, seatRoles map[string]string) Node {
	if len(seats) == 1 {
		seatKey := strings.TrimSpace(seats[0].SeatKey)
		userID := strings.TrimSpace(seats[0].UserID)
		if userID != "" && seatKey != "" {
			queuePath := strings.TrimSpace(seatRoles[seatKey])
			if branchPrefix(seatKey, queuePath) == "" {
				return SoloNode(userID, queuePath)
			}
		}
	}

	type branch struct {
		prefix string
		roles  map[string][]string // role -> user ids
	}
	branches := map[string]*branch{}
	ungrouped := map[string][]string{}

	for _, seat := range seats {
		userID := strings.TrimSpace(seat.UserID)
		seatKey := strings.TrimSpace(seat.SeatKey)
		if userID == "" || seatKey == "" {
			continue
		}
		role := strings.TrimSpace(seatRoles[seatKey])
		if role == "" {
			role = leafRoleFromSeatKey(seatKey)
		}
		prefix := branchPrefix(seatKey, role)
		if prefix == "" {
			ungrouped[role] = append(ungrouped[role], userID)
			continue
		}
		b, ok := branches[prefix]
		if !ok {
			b = &branch{prefix: prefix, roles: map[string][]string{}}
			branches[prefix] = b
		}
		b.roles[role] = append(b.roles[role], userID)
	}

	root := Node{}
	for role, userIDs := range ungrouped {
		if len(userIDs) == 1 {
			root.Children = append(root.Children, Node{Role: role, Members: userIDs})
		} else {
			root.Children = append(root.Children, Node{Role: role, Members: userIDs})
		}
	}

	prefixes := make([]string, 0, len(branches))
	for p := range branches {
		prefixes = append(prefixes, p)
	}
	sortStrings(prefixes)

	for _, prefix := range prefixes {
		b := branches[prefix]
		teamNode := Node{Role: branchRoleName(prefix)}
		roleNames := make([]string, 0, len(b.roles))
		for role := range b.roles {
			roleNames = append(roleNames, role)
		}
		sortStrings(roleNames)
		for _, role := range roleNames {
			teamNode.Children = append(teamNode.Children, Node{
				Role:    role,
				Members: b.roles[role],
			})
		}
		root.Children = append(root.Children, teamNode)
	}
	SortChildrenStable(&root)
	return root
}

func branchPrefix(seatKey, queuePath string) string {
	if queuePath == "" {
		return ""
	}
	suffix := "-" + queuePath
	if strings.HasSuffix(seatKey, suffix) {
		prefix := strings.TrimSuffix(seatKey, suffix)
		if prefix != "" {
			return prefix
		}
		return ""
	}
	// Pooled role seats: Team-1-Guesser-2 with queuePath Guesser
	marker := "-" + queuePath + "-"
	if idx := strings.Index(seatKey, marker); idx > 0 {
		return seatKey[:idx]
	}
	return ""
}

func branchRoleName(prefix string) string {
	if idx := strings.LastIndex(prefix, "-"); idx > 0 {
		left, right := prefix[:idx], prefix[idx+1:]
		if isInstanceToken(right) {
			return left
		}
	}
	parts := splitSeatKey(prefix)
	if len(parts) >= 2 && isInstanceToken(parts[len(parts)-1]) {
		return parts[0]
	}
	return prefix
}

func leafRoleFromSeatKey(seatKey string) string {
	parts := splitSeatKey(seatKey)
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	if isInstanceToken(parts[len(parts)-1]) {
		if len(parts) >= 2 {
			return parts[len(parts)-2]
		}
	}
	return parts[0]
}
