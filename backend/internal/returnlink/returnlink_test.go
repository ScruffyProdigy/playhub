package returnlink

import "testing"

func TestAppendMatchID(t *testing.T) {
	tests := []struct {
		name      string
		returnURL string
		matchID   string
		want      string
	}{
		{
			name:      "full hub url",
			returnURL: "https://joinquest.cc/return",
			matchID:   "550e8400-e29b-41d4-a716-446655440000",
			want:      "https://joinquest.cc/return?match=550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:      "no match id",
			returnURL: "https://joinquest.cc/return",
			matchID:   "",
			want:      "https://joinquest.cc/return",
		},
		{
			name:      "empty base",
			returnURL: "",
			matchID:   "",
			want:      "/return",
		},
		{
			name:      "existing query params",
			returnURL: "https://joinquest.cc/return?foo=bar",
			matchID:   "abc",
			want:      "https://joinquest.cc/return?foo=bar&match=abc",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := AppendMatchID(tc.returnURL, tc.matchID); got != tc.want {
				t.Fatalf("AppendMatchID() = %q, want %q", got, tc.want)
			}
		})
	}
}
