package tarot

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Slot is a chapter of the traveler's journey.
type Slot struct {
	Key  string
	Name string
}

// JourneySlots are assigned to drawn cards in order.
var JourneySlots = []Slot{
	{Key: "compass", Name: "Compass"},
	{Key: "coin", Name: "Coin"},
	{Key: "storm", Name: "Storm"},
	{Key: "campfire", Name: "Campfire"},
	{Key: "beacon", Name: "Beacon"},
}

// MajorArcana names in Rider–Waite order (The Fool = 0 … The World = 21).
var MajorArcana = []string{
	"The Fool",
	"The Magician",
	"The High Priestess",
	"The Empress",
	"The Emperor",
	"The Hierophant",
	"The Lovers",
	"The Chariot",
	"Strength",
	"The Hermit",
	"Wheel of Fortune",
	"Justice",
	"The Hanged Man",
	"Death",
	"Temperance",
	"The Devil",
	"The Tower",
	"The Star",
	"The Moon",
	"The Sun",
	"Judgement",
	"The World",
}

const cardCount = 22
const drawCount = 5

// CardName returns the Major Arcana name for index 0–21.
func CardName(index int) (string, error) {
	if index < 0 || index >= len(MajorArcana) {
		return "", fmt.Errorf("tarot: invalid card index %d", index)
	}
	return MajorArcana[index], nil
}

// Draw performs a random-without-replacement draw of five distinct cards from 0–21.
func Draw() ([]int, error) {
	pool := make([]int, cardCount)
	for i := range pool {
		pool[i] = i
	}
	draw := make([]int, drawCount)
	for i := 0; i < drawCount; i++ {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool)-i)))
		if err != nil {
			return nil, err
		}
		j := int(jBig.Int64()) + i
		pool[i], pool[j] = pool[j], pool[i]
		draw[i] = pool[i]
	}
	return draw, nil
}
