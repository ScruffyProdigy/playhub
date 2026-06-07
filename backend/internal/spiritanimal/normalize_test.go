package spiritanimal

import (
	"encoding/json"
	"testing"
)

func TestNormalizeQuestionsJSONNumberedSlots(t *testing.T) {
	raw := json.RawMessage(`{
		"cards": [
			{"slot": 1, "slot_name": "Compass", "card": "The Fool", "question": "Q1?",
			 "answers": [{"id":"A","label":"a"},{"id":"B","label":"b"},{"id":"C","label":"c"},{"id":"D","label":"d"},{"id":"E","label":"e"}]},
			{"slot": 2, "slot_name": "Coin", "card": "The Magician", "question": "Q2?",
			 "answers": {"A":"a","B":"b","C":"c","D":"d","E":"e"}},
			{"slot": "storm", "slot_name": "Storm", "card": "Death", "question": "Q3?",
			 "answers": ["a","b","c","d","e"]},
			{"slot": 4, "slot_name": "Campfire", "card": "Temperance", "question": "Q4?",
			 "answers": [{"label":"a"},{"label":"b"},{"label":"c"},{"label":"d"},{"label":"e"}]},
			{"slot": 5, "slot_name": "Beacon", "card": "The World", "question": "Q5?",
			 "answers": [{"id":"A","label":"a"},{"id":"B","label":"b"},{"id":"C","label":"c"},{"id":"D","label":"d"},{"id":"E","label":"e"}]}
		]
	}`)
	out, err := normalizeQuestionsJSON(raw)
	if err != nil {
		t.Fatalf("normalizeQuestionsJSON: %v", err)
	}
	var payload struct {
		Cards []struct {
			Slot    string `json:"slot"`
			Answers []struct {
				ID string `json:"id"`
			} `json:"answers"`
		} `json:"cards"`
	}
	if err := json.Unmarshal(out, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload.Cards[0].Slot != "compass" {
		t.Fatalf("slot 0: %q", payload.Cards[0].Slot)
	}
	if payload.Cards[2].Slot != "storm" {
		t.Fatalf("slot 2: %q", payload.Cards[2].Slot)
	}
	for i, card := range payload.Cards {
		if len(card.Answers) != 5 {
			t.Fatalf("card %d: %d answers", i, len(card.Answers))
		}
	}
}
