package partytree

import "strings"

// NormalizeForPlacement rewrites party trees so flat templates (e.g. {"count":2})
// place onto any open numbered seat, not only the player's local table seat key.
//
// Legacy table-backfill trees encoded a solo player as {role:"1", members:[user]}
// which blocks the second duelist from taking seat "2". Those become SoloNode layouts.
func NormalizeForPlacement(tree Node) Node {
	members := tree.AllMembers()
	if len(members) != 1 {
		return tree
	}

	member := members[0]
	queuePath := strings.TrimSpace(member.QueuePath)

	// Already catalog-style solo layout (empty default path).
	if len(tree.Members) == 1 && EffectiveRole(tree.Role) == queuePath {
		return tree
	}

	// Legacy flat layout: single child whose role is the local seat key ("1", "2", ...).
	if len(tree.Children) == 1 && len(tree.Members) == 0 {
		child := tree.Children[0]
		seatRole := EffectiveRole(child.Role)
		if len(child.Members) == 1 && isFlatSeatRole(seatRole) &&
			(queuePath == "" || queuePath == seatRole) {
			out := SoloNode(member.UserID, "")
			out.TableID = tree.TableID
			return out
		}
	}

	return tree
}

func isFlatSeatRole(role string) bool {
	role = strings.TrimSpace(role)
	if role == "" {
		return false
	}
	for _, r := range role {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
