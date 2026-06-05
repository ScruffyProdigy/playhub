package graph

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
)

func TestRoomGraphQLCreateJoinMessage(t *testing.T) {
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")

	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()

	_, hostCookie := createTestUserSession(t, ctx, env, cleaner)
	_, guestCookie := createTestUserSession(t, ctx, env, cleaner)

	var createResp struct {
		CreateRoom struct {
			ID         string
			InviteCode string
			JoinURL    string
		} `json:"createRoom"`
	}
	if err := env.Client.Post(`mutation {
		createRoom {
			id
			inviteCode
			joinUrl
		}
	}`, &createResp, client.AddCookie(hostCookie)); err != nil {
		t.Fatalf("createRoom: %v", err)
	}
	if createResp.CreateRoom.InviteCode == "" {
		t.Fatal("expected invite code")
	}
	if createResp.CreateRoom.JoinURL == "" {
		t.Fatal("expected join url")
	}

	var joinResp struct {
		JoinRoom struct {
			ID         string
			InviteCode string
		} `json:"joinRoom"`
	}
	joinQuery := `mutation Join($code: String!) {
		joinRoom(inviteCode: $code) { id inviteCode }
	}`
	if err := env.Client.Post(joinQuery, &joinResp, client.AddCookie(guestCookie), client.Var("code", createResp.CreateRoom.InviteCode)); err != nil {
		t.Fatalf("joinRoom: %v", err)
	}
	if joinResp.JoinRoom.ID != createResp.CreateRoom.ID {
		t.Fatalf("joinRoom returned wrong room")
	}

	var sendResp struct {
		SendRoomMessage struct {
			ID   string
			Body string
		} `json:"sendRoomMessage"`
	}
	sendQuery := `mutation Send($roomId: ID!, $body: String!) {
		sendRoomMessage(roomId: $roomId, body: $body) { id body }
	}`
	if err := env.Client.Post(sendQuery, &sendResp, client.AddCookie(hostCookie),
		client.Var("roomId", createResp.CreateRoom.ID),
		client.Var("body", "hello room"),
	); err != nil {
		t.Fatalf("sendRoomMessage: %v", err)
	}
	if sendResp.SendRoomMessage.Body != "hello room" {
		t.Fatalf("unexpected message body %q", sendResp.SendRoomMessage.Body)
	}
}
