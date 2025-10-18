package graph

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/scruffyprodigy/playhub/graph/generated"
	"github.com/scruffyprodigy/playhub/graph/model"
)

func TestMeResolver(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		Me struct {
			ID          string
			Email       *string
			DisplayName *string
			CreatedAt   string
		}
	}

	err := c.Post(`query { 
		me { 
			id 
			email 
			displayName 
			createdAt 
		} 
	}`, &resp)
	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	// Verify the response structure
	if resp.Me.ID == "" {
		t.Error("Expected me.id to be non-empty")
	}
	if resp.Me.Email == nil || *resp.Me.Email == "" {
		t.Error("Expected me.email to be non-empty")
	}
	if resp.Me.DisplayName == nil || *resp.Me.DisplayName == "" {
		t.Error("Expected me.displayName to be non-empty")
	}
	if resp.Me.CreatedAt == "" {
		t.Error("Expected me.createdAt to be non-empty")
	}
}

func TestGamesResolver(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		Games []struct {
			ID        string
			Name      string
			CreatedAt string
		}
	}

	err := c.Post(`query { 
		games { 
			id 
			name 
			createdAt 
		} 
	}`, &resp)
	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	// Verify we get some games back
	if len(resp.Games) == 0 {
		t.Error("Expected to get at least one game")
	}

	// Verify the structure of the first game
	game := resp.Games[0]
	if game.ID == "" {
		t.Error("Expected game.id to be non-empty")
	}
	if game.Name == "" {
		t.Error("Expected game.name to be non-empty")
	}
	if game.CreatedAt == "" {
		t.Error("Expected game.createdAt to be non-empty")
	}
}

func TestGamesWithPagination(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		Games []struct {
			ID   string
			Name string
		}
	}

	// Test with limit
	err := c.Post(`query { 
		games(limit: 1) { 
			id 
			name 
		} 
	}`, &resp)
	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	// Should get at most 1 game
	if len(resp.Games) > 1 {
		t.Errorf("Expected at most 1 game with limit=1, got %d", len(resp.Games))
	}
}

func TestCreateGameMutation(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		CreateGame struct {
			ID        string
			Name      string
			CreatedAt string
		}
	}

	err := c.Post(`mutation { 
		createGame(input: { name: "Test Game" }) { 
			id 
			name 
			createdAt 
		} 
	}`, &resp)
	if err != nil {
		t.Fatalf("GraphQL mutation failed: %v", err)
	}

	// Verify the response
	if resp.CreateGame.ID == "" {
		t.Error("Expected createGame.id to be non-empty")
	}
	if resp.CreateGame.Name != "Test Game" {
		t.Errorf("Expected createGame.name to be 'Test Game', got %q", resp.CreateGame.Name)
	}
	if resp.CreateGame.CreatedAt == "" {
		t.Error("Expected createGame.createdAt to be non-empty")
	}
}

func TestJoinGameMutation(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		JoinGame struct {
			Queued    bool
			SessionID *string
			JoinURL   *string
		}
	}

	err := c.Post(`mutation { 
		joinGame(gameId: "test-game-id") { 
			queued 
			sessionId 
			joinUrl 
		} 
	}`, &resp)
	if err != nil {
		t.Fatalf("GraphQL mutation failed: %v", err)
	}

	// Verify the response
	if !resp.JoinGame.Queued {
		t.Error("Expected joinGame.queued to be true")
	}
	if resp.JoinGame.SessionID == nil || *resp.JoinGame.SessionID == "" {
		t.Error("Expected joinGame.sessionId to be non-empty")
	}
	if resp.JoinGame.JoinURL == nil || *resp.JoinGame.JoinURL == "" {
		t.Error("Expected joinGame.joinUrl to be non-empty")
	}
}

func TestGoodsResolver(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		Goods []struct {
			ID          string
			Code        string
			Name        string
			Description *string
		}
	}

	err := c.Post(`query { 
		goods { 
			id 
			code 
			name 
			description 
		} 
	}`, &resp)
	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	// Verify we get some goods back
	if len(resp.Goods) == 0 {
		t.Error("Expected to get at least one good")
	}

	// Verify the structure of the first good
	good := resp.Goods[0]
	if good.ID == "" {
		t.Error("Expected good.id to be non-empty")
	}
	if good.Code == "" {
		t.Error("Expected good.code to be non-empty")
	}
	if good.Name == "" {
		t.Error("Expected good.name to be non-empty")
	}
}

func TestMyInventoryResolver(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		MyInventory []struct {
			Good struct {
				ID   string
				Code string
				Name string
			}
			Quantity  int
			GrantedAt string
		}
	}

	err := c.Post(`query { 
		myInventory { 
			good { 
				id 
				code 
				name 
			} 
			quantity 
			grantedAt 
		} 
	}`, &resp)
	if err != nil {
		t.Fatalf("GraphQL query failed: %v", err)
	}

	// Verify we get some inventory back
	if len(resp.MyInventory) == 0 {
		t.Error("Expected to get at least one inventory item")
	}

	// Verify the structure of the first inventory item
	item := resp.MyInventory[0]
	if item.Good.ID == "" {
		t.Error("Expected inventory.good.id to be non-empty")
	}
	if item.Good.Code == "" {
		t.Error("Expected inventory.good.code to be non-empty")
	}
	if item.Good.Name == "" {
		t.Error("Expected inventory.good.name to be non-empty")
	}
	if item.Quantity <= 0 {
		t.Error("Expected inventory.quantity to be positive")
	}
	if item.GrantedAt == "" {
		t.Error("Expected inventory.grantedAt to be non-empty")
	}
}

// Test error handling
func TestGameNotFound(t *testing.T) {
	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	c := client.New(srv)

	var resp struct {
		Game *struct {
			ID   string
			Name string
		}
	}

	err := c.Post(`query { 
		game(id: "non-existent-id") { 
			id 
			name 
		} 
	}`, &resp)

	// The resolver returns an error for non-existent games, which is expected
	if err == nil {
		t.Error("Expected GraphQL query to return an error for non-existent game")
	}

	// Verify the error message contains "game not found"
	if err != nil && err.Error() != `[{"message":"game not found","path":["game"]}]` {
		t.Errorf("Expected 'game not found' error, got: %v", err)
	}
}

// Test direct resolver calls
func TestDirectResolverCalls(t *testing.T) {
	resolver := &Resolver{}
	queryResolver := resolver.Query()
	mutationResolver := resolver.Mutation()

	// Test direct healthz call
	healthz, err := queryResolver.Healthz(context.Background())
	if err != nil {
		t.Fatalf("Direct healthz call failed: %v", err)
	}
	if healthz != "ok" {
		t.Errorf("Expected healthz to return 'ok', got %q", healthz)
	}

	// Test direct version call
	version, err := queryResolver.Version(context.Background())
	if err != nil {
		t.Fatalf("Direct version call failed: %v", err)
	}
	if version != "1.0.0" {
		t.Errorf("Expected version to return '1.0.0', got %q", version)
	}

	// Test direct createGame call
	game, err := mutationResolver.CreateGame(context.Background(), model.CreateGameInput{Name: "Direct Test Game"})
	if err != nil {
		t.Fatalf("Direct createGame call failed: %v", err)
	}
	if game.Name != "Direct Test Game" {
		t.Errorf("Expected game name to be 'Direct Test Game', got %q", game.Name)
	}
}
