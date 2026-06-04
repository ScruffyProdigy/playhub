package gameclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProvisionMatchSendsLobbyBlock(t *testing.T) {
	var gotBody map[string]any
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/matches" {
			http.NotFound(w, r)
			return
		}
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(srv.Close)

	client := NewClient()
	err := client.ProvisionMatch(t.Context(), ProvisionRequest{
		APIBaseURL:   srv.URL,
		ServiceToken: "lobby-svc-secret",
		LobbyID:      "https://joinquest.cc",
		Lobby: LobbyInfo{
			ReturnURL:    "https://joinquest.cc/return",
			GraphqlURL:   "https://joinquest.cc/graphql",
			ServiceToken: "lobby-svc-secret",
		},
		Assignment: Assignment{
			ExternalMatchID: "match-1",
			GameMode:        "duel",
			Seats: []AssignmentSeat{
				{SeatKey: "a", LobbyUserID: "user-a"},
			},
		},
	})
	if err != nil {
		t.Fatalf("ProvisionMatch: %v", err)
	}
	if gotBody["lobbyId"] != "https://joinquest.cc" {
		t.Fatalf("lobbyId = %v", gotBody["lobbyId"])
	}
	lobby, ok := gotBody["lobby"].(map[string]any)
	if !ok {
		t.Fatalf("lobby missing: %#v", gotBody)
	}
	if lobby["returnUrl"] != "https://joinquest.cc/return" {
		t.Fatalf("returnUrl = %v", lobby["returnUrl"])
	}
	if lobby["graphqlUrl"] != "https://joinquest.cc/graphql" {
		t.Fatalf("graphqlUrl = %v", lobby["graphqlUrl"])
	}
	if lobby["serviceToken"] != "lobby-svc-secret" {
		t.Fatalf("serviceToken = %v", lobby["serviceToken"])
	}
	assignment, ok := gotBody["assignment"].(map[string]any)
	if !ok {
		t.Fatalf("assignment missing: %#v", gotBody)
	}
	if assignment["gameMode"] != "duel" {
		t.Fatalf("gameMode = %v", assignment["gameMode"])
	}
	if gotAuth != "Bearer lobby-svc-secret" {
		t.Fatalf("Authorization = %q, want Bearer lobby-svc-secret", gotAuth)
	}
}

func TestProvisionMatchRequiresLobbyID(t *testing.T) {
	client := NewClient()
	err := client.ProvisionMatch(t.Context(), ProvisionRequest{
		APIBaseURL: "http://127.0.0.1:1",
		LobbyID:    "  ",
		Lobby: LobbyInfo{
			ReturnURL:  "https://joinquest.cc/return",
			GraphqlURL: "https://joinquest.cc/graphql",
		},
	})
	if err == nil {
		t.Fatal("expected error for empty lobby id")
	}
}
