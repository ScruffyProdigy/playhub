package spiritanimal

import (
	"fmt"
	"strings"
	"unicode"
)

var answerQuestionPrefixes = []string{
	"do ", "does ", "did ", "would ", "could ", "should ", "can ", "is ", "are ",
	"was ", "were ", "will ", "what ", "when ", "where ", "why ", "how ", "which ",
	"who ", "if ",
}

func looksLikeQuestion(text string) bool {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return false
	}
	if strings.HasSuffix(trimmed, "?") {
		return true
	}
	lower := strings.ToLower(trimmed)
	for _, prefix := range answerQuestionPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

var answerStopWords = map[string]struct{}{
	"a": {}, "an": {}, "the": {}, "my": {}, "your": {}, "to": {}, "of": {}, "and": {}, "or": {},
}

func answersTooSimilar(answers []struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}) string {
	labels := make([]string, 0, len(answers))
	for _, answer := range answers {
		labels = append(labels, strings.TrimSpace(answer.Label))
	}
	return answerLabelsTooSimilar(labels)
}

func answerLabelsTooSimilar(labels []string) string {
	return labelsTooSimilarStrict(labels, "answer")
}

func labelsTooSimilarStrict(labels []string, labelKind string) string {
	normalized := make([]string, len(labels))
	tokenSets := make([]map[string]struct{}, len(labels))
	firstWords := map[string]int{}

	for i, label := range labels {
		normalized[i] = normalizeAnswerLabel(label)
		tokenSets[i] = answerWordSet(label)
		if words := significantAnswerWords(label); len(words) > 0 {
			firstWords[words[0]]++
		}
	}

	seen := map[string]struct{}{}
	for i, norm := range normalized {
		if norm == "" {
			continue
		}
		if _, ok := seen[norm]; ok {
			return fmt.Sprintf("duplicate %s %q", labelKind, labels[i])
		}
		seen[norm] = struct{}{}
	}

	for i := 0; i < len(labels); i++ {
		for j := i + 1; j < len(labels); j++ {
			a, b := normalized[i], normalized[j]
			if a == b {
				return fmt.Sprintf("%q and %q are duplicate %ss", labels[i], labels[j], labelKind)
			}
			if strings.Contains(a, b) || strings.Contains(b, a) {
				return fmt.Sprintf("%q and %q overlap too much", labels[i], labels[j])
			}
			if answerWordOverlap(tokenSets[i], tokenSets[j]) >= 0.55 {
				return fmt.Sprintf("%q and %q are too similar", labels[i], labels[j])
			}
		}
	}

	for word, count := range firstWords {
		if count >= 3 {
			return fmt.Sprintf("too many %ss start with %q", labelKind, word)
		}
	}

	wordUse := map[string]int{}
	for _, label := range labels {
		seenInLabel := map[string]struct{}{}
		for _, word := range significantAnswerWords(label) {
			if _, ok := seenInLabel[word]; ok {
				continue
			}
			seenInLabel[word] = struct{}{}
			wordUse[word]++
		}
	}
	for word, count := range wordUse {
		if count >= 2 && len(word) >= 4 {
			return fmt.Sprintf("%ss repeat %q", labelKind, word)
		}
	}

	return ""
}

func normalizeAnswerLabel(label string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(label))), " ")
}

func answerWordSet(label string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, word := range significantAnswerWords(label) {
		set[word] = struct{}{}
	}
	return set
}

func significantAnswerWords(label string) []string {
	words := strings.Fields(strings.ToLower(label))
	out := make([]string, 0, len(words))
	for _, word := range words {
		word = strings.Trim(word, ".,!?\"'")
		if word == "" {
			continue
		}
		if _, stop := answerStopWords[word]; stop {
			continue
		}
		out = append(out, word)
	}
	return out
}

func answerWordOverlap(a, b map[string]struct{}) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	shared := 0
	for word := range a {
		if _, ok := b[word]; ok {
			shared++
		}
	}
	union := len(a) + len(b) - shared
	if union == 0 {
		return 0
	}
	return float64(shared) / float64(union)
}

func isLikelyAnimalSpecies(animal string) bool {
	animal = strings.TrimSpace(strings.ToLower(animal))
	if animal == "" {
		return false
	}
	disallowed := []string{
		"cloud", "puffball", "spirit", "object", "lamp", "coin", "star", "crystal",
		"robot", "machine", "vehicle", "plant", "tree", "flower", "rock", "stone",
		"blob", "shape", "human", "person", "wizard", "knight",
	}
	for _, word := range disallowed {
		if strings.Contains(animal, word) {
			return false
		}
	}
	if len(animal) < 2 {
		return false
	}
	for _, r := range animal {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
