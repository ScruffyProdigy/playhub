package gameclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProvisionMatchParsesLaunchURLs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"match": map[string]any{"id": "m1"},
			"launchUrls": map[string]string{
				"user-a": "https://play.example.com/?match=abc&seat=1",
				"user-b": "https://play.example.com/?match=abc&seat=2",
			},
		})
	}))
	t.Cleanup(srv.Close)

	client := NewClient()
	result, err := client.ProvisionMatch(t.Context(), ProvisionRequest{
		APIBaseURL: srv.URL,
		LobbyID:    "https://joinquest.cc",
		Lobby: LobbyInfo{
			ReturnURL:  "https://joinquest.cc/return",
			GraphqlURL: "https://joinquest.cc/graphql",
		},
		Assignment: Assignment{
			ExternalMatchID: "abc",
			GameMode:        "duel",
			Seats: []AssignmentSeat{
				{SeatKey: "1", LobbyUserID: "user-a"},
				{SeatKey: "2", LobbyUserID: "user-b"},
			},
		},
	})
	if err != nil {
		t.Fatalf("ProvisionMatch: %v", err)
	}
	if len(result.LaunchURLs) != 2 {
		t.Fatalf("LaunchURLs = %v", result.LaunchURLs)
	}
	if result.LaunchURLs["user-a"] != "https://play.example.com/?match=abc&seat=1" {
		t.Fatalf("user-a url = %q", result.LaunchURLs["user-a"])
	}
}

func TestParseProvisionResponseIgnoresEmptyEntries(t *testing.T) {
	result := parseProvisionResponse([]byte(`{"launchUrls":{"u1":"https://a.test","u2":"  "}}`))
	if len(result.LaunchURLs) != 1 {
		t.Fatalf("LaunchURLs = %v", result.LaunchURLs)
	}
}
