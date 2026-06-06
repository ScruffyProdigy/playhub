package graph

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func TestSelectStarterAvatarAndPlayerLookup(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()

	t.Setenv("LOBBY_PUBLIC_URL", "https://joinquest.test")
	t.Setenv("LOBBY_GAME_SERVICE_TOKEN", "test-game-service-token")

	user, err := env.Store.CreateUser(ctx, store.CreateUserParams{
		Email: "avatar-flow-" + uuid.NewString() + "@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)
	_, cookie := createTestUserSessionForUser(t, env, user.ID)

	selectQuery := `mutation Select($key: ID!) {
		selectStarterAvatar(key: $key) {
			id
			avatarKey
			avatarUrl
			avatarSource
		}
	}`
	selectBody := postGraphQL(t, env.Handler, selectQuery, map[string]any{"key": "campfire"}, cookie)
	var selectResp struct {
		Data struct {
			SelectStarterAvatar struct {
				AvatarKey *string `json:"avatarKey"`
				AvatarURL *string `json:"avatarUrl"`
			} `json:"selectStarterAvatar"`
		} `json:"data"`
	}
	if err := json.Unmarshal(selectBody, &selectResp); err != nil {
		t.Fatalf("decode select: %v body=%s", err, selectBody)
	}
	selected := selectResp.Data.SelectStarterAvatar
	if selected.AvatarKey == nil || *selected.AvatarKey != "campfire" {
		t.Fatalf("avatarKey: %+v", selected.AvatarKey)
	}
	if selected.AvatarURL == nil || *selected.AvatarURL != "https://joinquest.test/avatars/campfire.png" {
		t.Fatalf("avatarUrl: %+v", selected.AvatarURL)
	}

	playerQuery := `query Player($id: ID!) { player(id: $id) { id avatarUrl avatarSource } }`
	playerBody := postGraphQLWithBearer(t, env.Handler, "test-game-service-token", playerQuery, map[string]any{"id": user.ID.String()})
	var playerResp struct {
		Data struct {
			Player *struct {
				AvatarURL *string `json:"avatarUrl"`
			} `json:"player"`
		} `json:"data"`
	}
	if err := json.Unmarshal(playerBody, &playerResp); err != nil {
		t.Fatalf("decode player: %v", err)
	}
	if playerResp.Data.Player == nil || playerResp.Data.Player.AvatarURL == nil {
		t.Fatalf("expected player avatar, got %+v", playerResp.Data.Player)
	}

	catalogQuery := `query { starterAvatars { key name slot imageUrl } }`
	catalogBody := postGraphQL(t, env.Handler, catalogQuery, nil, cookie)
	var catalogResp struct {
		Data struct {
			StarterAvatars []struct {
				Key string `json:"key"`
			} `json:"starterAvatars"`
		} `json:"data"`
	}
	if err := json.Unmarshal(catalogBody, &catalogResp); err != nil {
		t.Fatalf("decode catalog: %v", err)
	}
	if len(catalogResp.Data.StarterAvatars) != 5 {
		t.Fatalf("expected 5 starter avatars, got %d", len(catalogResp.Data.StarterAvatars))
	}

	profileQuery := `mutation Profile($displayName: String!, $avatarKey: ID!) {
		updatePlayerProfile(displayName: $displayName, avatarKey: $avatarKey) {
			displayName
			avatarKey
			avatarUrl
		}
	}`
	profileBody := postGraphQL(t, env.Handler, profileQuery, map[string]any{
		"displayName": "River",
		"avatarKey":   "storm",
	}, cookie)
	var profileResp struct {
		Data struct {
			UpdatePlayerProfile struct {
				DisplayName *string `json:"displayName"`
				AvatarKey   *string `json:"avatarKey"`
				AvatarURL   *string `json:"avatarUrl"`
			} `json:"updatePlayerProfile"`
		} `json:"data"`
	}
	if err := json.Unmarshal(profileBody, &profileResp); err != nil {
		t.Fatalf("decode profile: %v body=%s", err, profileBody)
	}
	profile := profileResp.Data.UpdatePlayerProfile
	if profile.DisplayName == nil || *profile.DisplayName != "River" {
		t.Fatalf("displayName: %+v", profile.DisplayName)
	}
	if profile.AvatarKey == nil || *profile.AvatarKey != "storm" {
		t.Fatalf("avatarKey: %+v", profile.AvatarKey)
	}
}

func TestMyRoomIncludesMemberAvatars(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()

	t.Setenv("LOBBY_PUBLIC_URL", "https://joinquest.test")

	user, err := env.Store.CreateUser(ctx, store.CreateUserParams{
		Email: "avatar-room-" + uuid.NewString() + "@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	cleaner.TrackUser(user.ID)
	if _, err := env.Store.SetUserStarterAvatar(ctx, user.ID, "compass", "https://joinquest.test"); err != nil {
		t.Fatalf("SetUserStarterAvatar: %v", err)
	}
	_, cookie := createTestUserSessionForUser(t, env, user.ID)

	postGraphQL(t, env.Handler, `mutation { createRoom { inviteCode } }`, nil, cookie)

	myRoomBody := postGraphQL(t, env.Handler, `query {
		myRoom {
			members { id avatarUrl avatarKey }
			host { id avatarUrl }
		}
	}`, nil, cookie)
	var roomResp struct {
		Data struct {
			MyRoom *struct {
				Members []struct {
					AvatarURL *string `json:"avatarUrl"`
				} `json:"members"`
			} `json:"myRoom"`
		} `json:"data"`
	}
	if err := json.Unmarshal(myRoomBody, &roomResp); err != nil {
		t.Fatalf("decode myRoom: %v body=%s", err, myRoomBody)
	}
	if roomResp.Data.MyRoom == nil || len(roomResp.Data.MyRoom.Members) == 0 {
		t.Fatal("expected room members")
	}
	if roomResp.Data.MyRoom.Members[0].AvatarURL == nil {
		t.Fatal("expected member avatarUrl")
	}
}
