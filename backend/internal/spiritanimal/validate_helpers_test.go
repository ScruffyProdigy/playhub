package spiritanimal

import "testing"

func TestLooksLikeQuestion(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{"A steady lantern", false},
		{"Would you turn back?", true},
		{"What path calls to you", true},
		{"Curiosity", false},
		{"Do you trust the storm", true},
	}
	for _, tc := range cases {
		if got := looksLikeQuestion(tc.text); got != tc.want {
			t.Fatalf("looksLikeQuestion(%q) = %v, want %v", tc.text, got, tc.want)
		}
	}
}

func TestIsLikelyAnimalSpecies(t *testing.T) {
	if !isLikelyAnimalSpecies("fox") {
		t.Fatal("fox should be valid")
	}
	if isLikelyAnimalSpecies("Cloud Puffball") {
		t.Fatal("cloud puffball should be invalid")
	}
	if isLikelyAnimalSpecies("") {
		t.Fatal("empty should be invalid")
	}
}

func TestAnswerLabelsTooSimilar(t *testing.T) {
	if msg := answerLabelsTooSimilar([]string{
		"Watchful eye", "Watchful path", "Bold leap", "Quiet trust", "Ancient map",
	}); msg == "" {
		t.Fatal("expected overlap detection")
	}
	if msg := answerLabelsTooSimilar([]string{
		"Steady lantern", "Wild impulse", "Shared feast", "Hidden door", "First light",
	}); msg != "" {
		t.Fatalf("expected distinct answers, got %q", msg)
	}
}

func TestValidateSubmittedAnswers(t *testing.T) {
	questions := mockQuestionsJSON(`{"draw":[1,2,3,4,5]}`)
	if err := validateSubmittedAnswers(questions, []string{"A", "B", "C", "D", "E"}); err != nil {
		t.Fatalf("expected valid answers: %v", err)
	}
	if err := validateSubmittedAnswers(questions, []string{"A", "B", "C", "D", "Z"}); err == nil {
		t.Fatal("expected invalid answer to fail")
	}
}

func TestValidateTotemsJSONAcceptsMockTotems(t *testing.T) {
	if err := validateTotemsJSON(mockTotemsJSON()); err != nil {
		t.Fatalf("mock totems should validate: %v", err)
	}
}
