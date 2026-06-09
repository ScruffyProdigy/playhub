package partytree

import (
	"sort"
	"strings"
)

func normalizeRoot(root Node) Node {
	if IsRootRole(root.Role) && len(root.Members) == 0 && len(root.Children) > 0 {
		return Node{Children: root.Children}
	}
	return root
}

func splitSeatKey(seatKey string) []string {
	parts := strings.Split(strings.TrimSpace(seatKey), "-")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func isInstanceToken(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func sortStrings(list []string) {
	sort.Strings(list)
}

// SortChildrenStable sorts children by role for deterministic matching.
func SortChildrenStable(n *Node) {
	if n == nil {
		return
	}
	sort.SliceStable(n.Children, func(i, j int) bool {
		return EffectiveRole(n.Children[i].Role) < EffectiveRole(n.Children[j].Role)
	})
	for i := range n.Children {
		SortChildrenStable(&n.Children[i])
	}
}
