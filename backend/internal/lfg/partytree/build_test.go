package partytree

import "testing"

func TestBuildFromPinnedSeats_FlatDuelUsesDefaultPath(t *testing.T) {
	tree := BuildFromPinnedSeats(
		[]PinnedSeat{{UserID: "user-a", SeatKey: "1"}},
		map[string]string{"1": "", "2": ""},
	)
	if len(tree.AllMembers()) != 1 || tree.AllMembers()[0].UserID != "user-a" {
		t.Fatalf("members = %+v", tree.AllMembers())
	}
	if EffectiveRole(tree.Role) != "" {
		t.Fatalf("role = %q, want empty default path", tree.Role)
	}
}
