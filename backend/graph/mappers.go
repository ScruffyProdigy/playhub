package graph

import (
	"github.com/scruffyprodigy/playhub/graph/model"
	"github.com/scruffyprodigy/playhub/internal/store"
)

// ToGraphQLUser maps a store user to the GraphQL model.
func ToGraphQLUser(user *store.User) *model.User {
	if user == nil {
		return nil
	}
	email := user.Email
	displayName := user.DisplayName
	return &model.User{
		ID:          user.ID.String(),
		Email:       &email,
		DisplayName: &displayName,
		CreatedAt:   user.CreatedAt,
	}
}

// ToGraphQLGame maps a store game to the GraphQL model.
func ToGraphQLGame(game *store.Game) *model.Game {
	if game == nil {
		return nil
	}
	return &model.Game{
		ID:        game.ID.String(),
		Name:      game.Name,
		CreatedAt: game.CreatedAt,
	}
}

// ToGraphQLGames maps a slice of store games to GraphQL models.
func ToGraphQLGames(games []store.Game) []*model.Game {
	result := make([]*model.Game, len(games))
	for i := range games {
		result[i] = ToGraphQLGame(&games[i])
	}
	return result
}

// ToGraphQLSession maps a store session to the GraphQL model.
func ToGraphQLSession(session *store.Session, game *store.Game) *model.Session {
	if session == nil {
		return nil
	}
	return &model.Session{
		ID:        session.ID.String(),
		Game:      ToGraphQLGame(game),
		Status:    ToGraphQLSessionStatus(session.Status),
		CreatedAt: session.StartedAt,
	}
}

// ToGraphQLSessionStatus maps database session status values to GraphQL enums.
func ToGraphQLSessionStatus(dbStatus string) model.SessionStatus {
	switch dbStatus {
	case "active":
		return model.SessionStatusActive
	case "completed", "cancelled":
		return model.SessionStatusEnded
	default:
		return model.SessionStatusPending
	}
}

// ToGraphQLDigitalGood maps a store digital good to the GraphQL model.
func ToGraphQLDigitalGood(good *store.DigitalGood) *model.DigitalGood {
	if good == nil {
		return nil
	}
	code := good.ID.String()
	return &model.DigitalGood{
		ID:          good.ID.String(),
		Code:        code,
		Name:        good.Name,
		Description: good.Description,
	}
}

// ToGraphQLDigitalGoods maps a slice of store goods to GraphQL models.
func ToGraphQLDigitalGoods(goods []store.DigitalGood) []*model.DigitalGood {
	result := make([]*model.DigitalGood, len(goods))
	for i := range goods {
		result[i] = ToGraphQLDigitalGood(&goods[i])
	}
	return result
}

// ToGraphQLEntitlement maps a store inventory item to the GraphQL model.
func ToGraphQLEntitlement(item *store.InventoryItem) *model.Entitlement {
	if item == nil {
		return nil
	}
	return &model.Entitlement{
		Good:      ToGraphQLDigitalGood(&item.Good),
		Quantity:  item.Quantity,
		GrantedAt: item.AcquiredAt,
	}
}

// ToGraphQLEntitlements maps a slice of inventory items to GraphQL models.
func ToGraphQLEntitlements(items []store.InventoryItem) []*model.Entitlement {
	result := make([]*model.Entitlement, len(items))
	for i := range items {
		result[i] = ToGraphQLEntitlement(&items[i])
	}
	return result
}

// ToGraphQLUsers maps a slice of store users to GraphQL models.
func ToGraphQLUsers(users []store.User) []*model.User {
	result := make([]*model.User, len(users))
	for i := range users {
		result[i] = ToGraphQLUser(&users[i])
	}
	return result
}
